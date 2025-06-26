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
	libreclient "github.com/tikhonp/medsenger-freestylelibre-bot/util/libre_client"
)

// LibreClient contains information about LibreLinkUp account connected to contraact.
type LibreClient struct {
	ID           int        `db:"id"`
	Email        string     `db:"email"`
	Password     string     `db:"password"`
	Token        *string    `db:"token"`
	TokenExpires *time.Time `db:"token_expires"`
	LastSyncDate *time.Time `db:"last_sync_date"`
	PatientID    *uuid.UUID `db:"patient_id"`
	ContractID   int        `db:"contract_id"`
	IsValid      bool       `db:"is_valid"`
}

var (
	ErrLibreClientNotFound            = errors.New("libre client not found")
	ErrLibreAccountConnectionsIsEmpty = errors.New("account connections is empty")
)

var llum = libreclient.NewLibreLinkUpManager()

func NewLibreClient(email string, password string, contractID int) (*LibreClient, error) {
	contract, err := GetContractByID(contractID)
	if err != nil {
		return nil, err
	}
	if contract.LibreClientID != nil {
		libreClient, err := contract.LibreClient()
		if err != nil {
			if errors.Is(err, ErrLibreClientNotFound) {
				contract.LibreClientID = nil
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
	err = db.Get(&lc, query, email, password, contractID)
	if err != nil {
		return nil, err
	}
	contract.LibreClientID = &lc.ID
	err = contract.Save()
	return &lc, err
}

func (lc *LibreClient) Contract() (*Contract, error) {
	return GetContractByID(lc.ContractID)
}

func (lc *LibreClient) Save() error {
	const query = `
		INSERT INTO libre_clients (id, email, password, token, token_expires, last_sync_date, patient_id, contract_id, is_valid)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) ON CONFLICT (id)
		DO UPDATE SET email = EXCLUDED.email, password = EXCLUDED.password, token = EXCLUDED.token, token_expires = EXCLUDED.token_expires, last_sync_date = EXCLUDED.last_sync_date, patient_id = EXCLUDED.patient_id, contract_id = EXCLUDED.contract_id, is_valid = EXCLUDED.is_valid
	`
	_, err := db.Exec(query, lc.ID, lc.Email, lc.Password, lc.Token, lc.TokenExpires, lc.LastSyncDate, lc.PatientID, lc.ContractID, lc.IsValid)
	return err
}

func GetLibreClientByID(id int) (*LibreClient, error) {
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

// sendMessageToChat sends message to doctor AND patient with URGENT setting
func (lc *LibreClient) sendMessageToChat(mc *maigo.Client, text string) {
	if lc.IsValid {
		_, err := mc.SendMessage(lc.ContractID, text, maigo.Urgent())
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
		lc.sendMessageToChat(mc, fmt.Sprintf("Ошибка синхронизации с сервисом Libre Link Up. Не удалось войти в систему, проверьте логин и пароль. Ошибка: %s", err.Error()))
		return err
	}
	lc.Token = &user.AuthTicket.Token
	lc.TokenExpires = &user.AuthTicket.Expires.Time
	return lc.Save()
}

func (lc *LibreClient) fetchPatientID(mc *maigo.Client) error {
	connections, err := llum.FetchConnections(*lc.Token)
	if err != nil {
		lc.sendMessageToChat(mc, fmt.Sprintf("Ошибка синхронизации с сервисом Libre Link Up. Ошибка: %s", err.Error()))
		return err
	}
	if len(connections) == 0 {
		lc.sendMessageToChat(mc, "Ошибка синхронизации с сервисом Libre Link Up. Не найдено подключенных пациентов для отслеживания.")
		return ErrLibreAccountConnectionsIsEmpty
	}
	lc.PatientID = &connections[0].PatientID
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
	log.Printf("Fetching data for contract %d", lc.ContractID)

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
	if lc.PatientID == nil {
		log.Printf("Patient id is nil. Fetching new patient id")
		err := lc.fetchPatientID(mc)
		if err != nil {
			return err
		}
	}

	graph, err := llum.FetchData(*lc.PatientID, *lc.Token)
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
		_, err = mc.AddRecords(lc.ContractID, records)
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
