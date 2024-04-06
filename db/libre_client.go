package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/TikhonP/maigo"
	"github.com/TikhonP/medsenger-freestylelibre-bot/util"
	"github.com/google/uuid"
)

type LibreClient struct {
	Id           int        `db:"id"`
	Email        string     `db:"email"`
	Password     string     `db:"password"`
	Token        *string    `db:"token"`
	TokenExpires *time.Time `db:"token_expires"`
	LastSyncDate *time.Time `db:"last_sync_date"`
	PatientId    *string    `db:"patient_id"`
}

var ErrLibreClientNotFound = errors.New("Libre Client not flound")

func NewLibreClient(email string, password string, contractId int) (*LibreClient, error) {
	contract, err := GetContractById(contractId)
	if err != nil {
		return nil, err
	}
	if contract.LibreClient != nil {
		libreClient, err := GetLibreClientById(*contract.LibreClient)
		if err != nil {
			if errors.Is(err, ErrLibreClientNotFound) {
				contract.LibreClient = nil
			} else {
				return nil, err
			}
		} else {
			libreClient.Email = email
			libreClient.Password = password
			libreClient.Token = nil
			err = libreClient.Save()
			fmt.Println("Save libre client")
			return libreClient, err
		}
	}
	query := `INSERT INTO libre_clients (email, password) VALUES ($1, $2) RETURNING *`
	var lc LibreClient
	err = db.Get(&lc, query, email, password)
	fmt.Println("inset new libre client")
	if err != nil {
		return nil, err
	}
	contract.LibreClient = &lc.Id
	err = contract.Save()
	fmt.Println("save contract")
	return &lc, err
}

func (lc *LibreClient) Save() error {
	query := `
		INSERT INTO libre_clients (id, email, password, token, token_expires, last_sync_date, patient_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT (id)
		DO UPDATE SET email = EXCLUDED.email, password = EXCLUDED.password, token = EXCLUDED.token, token_expires = EXCLUDED.token_expires, last_sync_date = EXCLUDED.last_sync_date, patient_id = EXCLUDED.patient_id
	`
	_, err := db.Exec(query, lc.Id, lc.Email, lc.Password, lc.Token, lc.TokenExpires, lc.LastSyncDate, lc.PatientId)
	return err
}

func GetLibreClientById(id int) (*LibreClient, error) {
	libreClient := new(LibreClient)
	err := db.Get(libreClient, "SELECT * FROM libre_clients WHERE id = $1", id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrLibreClientNotFound
	}
	return libreClient, err
}

func GetActiveLibreClientToFetch() ([]LibreClient, error) {
	query := `
        SELECT lc.* 
        FROM contracts c
        JOIN libre_clients lc ON lc.id = c.libre_client
        WHERE c.is_active
    `
	clients := []LibreClient{}
	err := db.Select(&clients, query)
	// rows, err := db.Query(query)
	// defer rows.Close()
	// for rows.Next() {
	// 	var (
	// 		name string
	// 		age  int
	// 	)
	// 	if err := rows.Scan(&name, &age); err != nil {
	// 		panic(err)
	// 	}
	// 	fmt.Printf("%s is %d\n", name, age)
	// }
	// if err := rows.Err(); err != nil {
	// 	panic(err)
	// }

	return clients, err
}

func (lc *LibreClient) FetchData(mc *maigo.Client) error {
	llum := util.NewLibreLinkUpManager()

	fmt.Println("Fetching token...")

	// fetch token
	if !(lc.Token != nil && time.Now().Before(*lc.TokenExpires)) {
		fmt.Println("fetch")
		user, err := llum.Login(lc.Email, lc.Password)
		if err != nil {
			return err
		}
		fmt.Println("saving")
		lc.Token = &user.AuthTicket.Token
		lc.TokenExpires = &user.AuthTicket.Expires.Time
		fmt.Printf("lc.Token: %v\n", *lc.Token)
		fmt.Printf("lc: %v\n", lc)
		fmt.Printf("user: %v\n", user.AuthTicket.Expires.Time)
		err = lc.Save()
		if err != nil {
			return err
		}
	}

	fmt.Println("Fetching patient id...")

	// fetch patient id
	if lc.PatientId == nil {
		connections, err := llum.FetchConnections(*lc.Token)
		if err != nil {
			return err
		}
		if len(connections.Data) == 0 {
			return errors.New("Account connections is empty")
		}
		patientId := connections.Data[0].PatientId.String()
		lc.PatientId = &patientId
		err = lc.Save()
		if err != nil {
			return err
		}
	}

	fmt.Println("Fetching data...")

	// fetch data
	patientUUID, err := uuid.Parse(*lc.PatientId)
	if err != nil {
		return err
	}
	graph, err := llum.FetchData(patientUUID, *lc.Token)

	println("MEASUREMENTS")

	for _, item := range graph.Data.Mesurements {
		fmt.Printf("%+v\n", item)
	}

	return err
}
