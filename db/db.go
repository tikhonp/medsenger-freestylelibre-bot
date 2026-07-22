// Package db provides a simple interface to interact with the database.
package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/tikhonp/medsenger-freestylelibre-bot/util"
)

// schema holds the idempotent DDL applied on every startup. Evolve the schema
// by appending idempotent statements here.
//
// Each statement is executed in its own (autocommit) transaction — never as one
// combined multi-statement string. A combined string runs as a single
// transaction that holds ACCESS EXCLUSIVE on every table it touches until the
// end, in the order the statements appear (libre_clients, then contracts). The
// worker's active-clients SELECT takes ACCESS SHARE on the same tables in the
// opposite order (contracts, then libre_clients), so an overlapping startup
// migration and fetch tick would deadlock (pq 40P01). Running one statement at
// a time releases each table lock before the next is taken, so the migration
// never holds both at once and no lock-wait cycle can form.
var schema = []string{
	`CREATE TABLE IF NOT EXISTS public.contracts (
	    id INTEGER PRIMARY KEY NOT NULL,
	    is_active BOOLEAN NOT NULL,
	    agent_token VARCHAR(254),
	    patient_name VARCHAR(254),
	    patient_email VARCHAR(254),
        locale VARCHAR(5) NULL,
        libre_client INTEGER
	)`,

	`CREATE TABLE IF NOT EXISTS public.libre_clients (
        id SERIAL PRIMARY KEY NOT NULL,
        email VARCHAR(254) NOT NULL,
        password VARCHAR(254) NOT NULL,
        token VARCHAR(1000),
        last_sync_date TIMESTAMP,
        token_expires TIMESTAMP,
        patient_id VARCHAR(254),
        contract_id INTEGER NOT NULL
    )`,

	`ALTER TABLE public.libre_clients ADD COLUMN IF NOT EXISTS is_valid BOOLEAN NOT NULL DEFAULT TRUE`,
	`ALTER TABLE public.libre_clients ADD COLUMN IF NOT EXISTS sync_success_msg_sent BOOLEAN NOT NULL DEFAULT FALSE`,

	`ALTER TABLE public.libre_clients ADD COLUMN IF NOT EXISTS account_id VARCHAR(256) DEFAULT NULL`,

	`ALTER TABLE public.contracts DROP COLUMN IF EXISTS agent_token`,
}

// db is a global database.
//
// Yes, im dumb and i use global varibles for db.
// It's my second project on go, i think you can forgive me.
var db *sqlx.DB

func dataSourceName(cfg *util.Database) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Dbname)
}

// MustConnect connects to Postgres and applies the schema, one statement per transaction.
func MustConnect(cfg *util.Database) {
	db = sqlx.MustConnect("postgres", dataSourceName(cfg))
	for _, stmt := range schema {
		db.MustExec(stmt)
	}
}
