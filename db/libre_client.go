package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/TikhonP/maigo"
	"github.com/getsentry/sentry-go"
	"github.com/google/uuid"
	"github.com/tikhonp/medsenger-freestylelibre-bot/util/libre_client"
)

// LibreClient contains information about LibreLinkUp account connected to contraact.
type LibreClient struct {
	Id           int        `db:"id"`
	Email        string     `db:"email"`
	Password     string     `db:"password"`
	Token        *string    `db:"token"`
	TokenExpires *time.Time `db:"token_expires"`
	LastSyncDate *time.Time `db:"last_sync_date"`
	PatientId    *uuid.UUID `db:"patient_id"`
	ContractId   int        `db:"contract_id"`
	IsValid      bool       `db:"is_valid"`
}

var ErrLibreClientNotFound = errors.New("libre client not found")

var llum = libreclient.NewLibreLinkUpManager()

func NewLibreClient(email string, password string, contractId int) (*LibreClient, error) {
	contract, err := GetContractById(contractId)
	if err != nil {
		return nil, err
	}
	if contract.LibreClientId != nil {
		libreClient, err := contract.LibreClient()
		if err != nil {
			if errors.Is(err, ErrLibreClientNotFound) {
				contract.LibreClientId = nil
			} else {
				return nil, err
			}
		} else {
			libreClient.Email = email
			libreClient.Password = password
			libreClient.Token = nil
			libreClient.IsValid = true
			err = libreClient.Save()
			return libreClient, err
		}
	}
	const query = `INSERT INTO libre_clients (email, password, contract_id) VALUES ($1, $2, $3) RETURNING *`
	var lc LibreClient
	err = db.Get(&lc, query, email, password, contractId)
	if err != nil {
		return nil, err
	}
	contract.LibreClientId = &lc.Id
	err = contract.Save()
	return &lc, err
}

func (lc *LibreClient) Contract() (*Contract, error) {
	return GetContractById(lc.ContractId)
}

func (lc *LibreClient) Save() error {
	const query = `
		INSERT INTO libre_clients (id, email, password, token, token_expires, last_sync_date, patient_id, contract_id, is_valid)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) ON CONFLICT (id)
		DO UPDATE SET email = EXCLUDED.email, password = EXCLUDED.password, token = EXCLUDED.token, token_expires = EXCLUDED.token_expires, last_sync_date = EXCLUDED.last_sync_date, patient_id = EXCLUDED.patient_id, contract_id = EXCLUDED.contract_id, is_valid = EXCLUDED.is_valid
	`
	_, err := db.Exec(query, lc.Id, lc.Email, lc.Password, lc.Token, lc.TokenExpires, lc.LastSyncDate, lc.PatientId, lc.ContractId, lc.IsValid)
	return err
}

func GetLibreClientById(id int) (*LibreClient, error) {
	libreClient := new(LibreClient)
	err := db.Get(libreClient, `SELECT * FROM libre_clients WHERE id = $1`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrLibreClientNotFound
	}
	return libreClient, err
}

func GetActiveLibreClientToFetch() ([]LibreClient, error) {
	const query = `
        SELECT lc.* 
        FROM contracts c
        JOIN libre_clients lc ON lc.id = c.libre_client
        WHERE c.is_active
    `
	clients := []LibreClient{}
	err := db.Select(&clients, query)
	return clients, err
}

func (lc *LibreClient) sendMessageToDoctor(mc *maigo.Client, text string) {
	if lc.IsValid {
		_, err := mc.SendMessage(lc.ContractId, text, maigo.Urgent())
		if err != nil {
			sentry.CaptureException(err)
			return
		}
		lc.IsValid = false
		err = lc.Save()
		if err != nil {
			sentry.CaptureException(err)
			return
		}
	}
}

func (lc *LibreClient) fetchToken(mc *maigo.Client) error {
	user, err := llum.Login(lc.Email, lc.Password)
	if err != nil {
		lc.sendMessageToDoctor(mc, fmt.Sprintf("Ошибка синхронизации с сервисом Libre Link Up. Не удалось войти в систему, проверьте логин и пароль. Ошибка: %s", err.Error()))
		return err
	}
	lc.Token = &user.AuthTicket.Token
	lc.TokenExpires = &user.AuthTicket.Expires.Time
	return lc.Save()
}

func (lc *LibreClient) fetchPatientId(mc *maigo.Client) error {
	connections, err := llum.FetchConnections(*lc.Token)
	if err != nil {
		lc.sendMessageToDoctor(mc, fmt.Sprintf("Ошибка синхронизации с сервисом Libre Link Up. Ошибка: %s", err.Error()))
		return err
	}
	if len(connections) == 0 {
		lc.sendMessageToDoctor(mc, "Ошибка синхронизации с сервисом Libre Link Up. Не найдено подключенных пациентов для отслеживания.")
		return errors.New("account connections is empty")
	}
	lc.PatientId = &connections[0].PatientId
	return lc.Save()
}

// function for get the latest FactoryTimestamp from []util.GlucoseMeasurement
func getLatestTimestamp(data []libreclient.GlucoseMeasurement) *time.Time {
	if len(data) == 0 {
		return nil
	}
	timestamp := data[0].FactoryTimestamp.Time
	for _, item := range data {
		if item.FactoryTimestamp.After(timestamp) {
			timestamp = item.FactoryTimestamp.Time
		}
	}
	return &timestamp
}

func (lc *LibreClient) FetchData(mc *maigo.Client) error {
	log.Printf("Fetching data for contract %d", lc.ContractId)

	now := time.Now().UTC()

	// fetch token
	if lc.Token == nil || now.After(*lc.TokenExpires) {
		log.Printf("Token is nil or expired. Fetching new token")
		err := lc.fetchToken(mc)
		if err != nil {
			return err
		}
	}

	// fetch patient id
	if lc.PatientId == nil {
		log.Printf("Patient id is nil. Fetching new patient id")
		err := lc.fetchPatientId(mc)
		if err != nil {
			return err
		}
	}

	graph, err := llum.FetchData(*lc.PatientId, *lc.Token)
	if err != nil {
		return err
	}

	log.Printf("Got graph with %d items", len(graph.Mesurements))

	var records []maigo.Record
	for _, item := range graph.Mesurements {
		if lc.LastSyncDate == nil || item.FactoryTimestamp.After(*lc.LastSyncDate) {
			records = append(records, maigo.NewRecord("glukose", item.ValueAsString(), item.FactoryTimestamp.Time))
		}
	}
	if len(records) > 0 {
		log.Printf("Sending %d records to medsenger", len(records))
		_, err = mc.AddRecords(lc.ContractId, records)
		if err != nil {
			return err
		}
	}

	// update last sync date
	lastSyncDate := getLatestTimestamp(graph.Mesurements)
	if lastSyncDate != nil {
		lc.LastSyncDate = lastSyncDate
		return lc.Save()
	}

	return nil
}
