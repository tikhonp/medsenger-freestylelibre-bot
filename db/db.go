// Package db provides a simple interface to interact with the database.
package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/tikhonp/medsenger-freestylelibre-bot/util"
)

const schema = `
	CREATE TABLE IF NOT EXISTS public.contracts (
	    id INTEGER PRIMARY KEY NOT NULL,
	    is_active BOOLEAN NOT NULL,
	    agent_token VARCHAR(254),
	    patient_name VARCHAR(254),
	    patient_email VARCHAR(254),
        locale VARCHAR(5) NULL,
        libre_client INTEGER
	);

    CREATE TABLE IF NOT EXISTS public.libre_clients (
        id SERIAL PRIMARY KEY NOT NULL,
        email VARCHAR(254) NOT NULL,
        password VARCHAR(254) NOT NULL,
        token VARCHAR(1000),
        last_sync_date TIMESTAMP,
        token_expires TIMESTAMP,
        patient_id VARCHAR(254), 
        contract_id INTEGER NOT NULL
    );

	ALTER TABLE public.libre_clients ADD COLUMN IF NOT EXISTS is_valid BOOLEAN NOT NULL DEFAULT TRUE;
	ALTER TABLE public.libre_clients ADD COLUMN IF NOT EXISTS sync_success_msg_sent BOOLEAN NOT NULL DEFAULT FALSE;
`

// db is a global database.
//
// Yes, im dumb and i use global varibles for db.
// It's my second project on go, i think you can forgive me.
var db *sqlx.DB

func dataSourceName(cfg *util.Database) string {
	return fmt.Sprintf("user=%s dbname=%s sslmode=disable password=%s host=%s", cfg.User, cfg.Dbname, cfg.Password, cfg.Host)
}

// MustConnect creates a new in-memory SQLite database and initializes it with the schema.
func MustConnect(cfg *util.Database) {
	db = sqlx.MustConnect("postgres", dataSourceName(cfg))
	db.MustExec(schema)
}
