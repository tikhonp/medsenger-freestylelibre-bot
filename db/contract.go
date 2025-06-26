package db

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Contract represents Medsenger contract.
// Create on agent /init and persist during agent lifecycle.
type Contract struct {
	ID            int     `db:"id"`
	IsActive      bool    `db:"is_active"`
	AgentToken    *string `db:"agent_token"`
	PatientName   *string `db:"patient_name"`
	PatientEmail  *string `db:"patient_email"`
	Locale        *string `db:"locale"`
	LibreClientID *int    `db:"libre_client"`
}

// Save on Contract saves structure to database.
func (c *Contract) Save() error {
	const query = `
		INSERT INTO contracts (id, is_active, agent_token, patient_name, patient_email, locale, libre_client)
		VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT (id)
		DO UPDATE SET is_active = EXCLUDED.is_active, agent_token = EXCLUDED.agent_token, patient_name = EXCLUDED.patient_name, patient_email = EXCLUDED.patient_email, locale = EXCLUDED.locale, libre_client = EXCLUDED.libre_client
	`
	_, err := db.Exec(query, c.ID, c.IsActive, c.AgentToken, c.PatientName, c.PatientEmail, c.Locale, c.LibreClientID)
	return err
}

func (c *Contract) LibreClient() (*LibreClient, error) {
	if c.LibreClientID == nil {
		return nil, ErrLibreClientNotFound
	}
	return GetLibreClientByID(*c.LibreClientID)
}

// GetActiveContractIds returns all active contracts ids.
// Use it for medsenger status endpoint.
func GetActiveContractIds() ([]int, error) {
	var contractIds = make([]int, 0)
	err := db.Select(&contractIds, `SELECT id FROM contracts WHERE is_active = true`)
	return contractIds, err
}

// MarkInactiveContractWithID sets contract with id to inactive.
// Use it for medsenger remove endpoint.
// Equivalent to DELETE FROM contracts WHERE id = ?.
func MarkInactiveContractWithID(id int) error {
	_, err := db.Exec(`UPDATE contracts SET is_active = false WHERE id = $1`, id)
	return err
}

// GetContractByID returns contract with specified id.
func GetContractByID(id int) (*Contract, error) {
	contract := new(Contract)
	err := db.Get(contract, `SELECT * FROM contracts WHERE id = $1`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, echo.NewHTTPError(http.StatusNotFound, "Contract not found")
	}
	return contract, err
}

