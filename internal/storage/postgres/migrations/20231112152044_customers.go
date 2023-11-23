package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCustomers, downCustomers)
}

func upCustomers(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(
		`CREATE TABLE 
			customer(
				id text PRIMARY KEY,
				name text,
				phone text,
				tag text,
				operator text)`)
	return err
}

func downCustomers(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec("DROP TABLE customer")
	return err
}
