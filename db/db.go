package db

import (
	"fmt"
	"github.com/TikhonP/medsenger-freestylelibre-bot/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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
`

// db is a global database.
//
// Yes, im dumb and i use global varibles for db.
// It's my second project in go, i think you can forgive me.
var db *sqlx.DB

func dataSourceName(cfg *config.Database) string {
	return fmt.Sprintf("user=%s dbname=%s sslmode=disable password=%s host=%s", cfg.User, cfg.Dbname, cfg.Password, cfg.Host)
}

// MustConnect creates a new in-memory SQLite database and initializes it with the schema.
func MustConnect(cfg *config.Database) {
	db = sqlx.MustConnect("postgres", dataSourceName(cfg))
	db.MustExec(schema)
}
