package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upMessages, downMessages)
}

func upMessages(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(
		`CREATE TABLE messages(
			id BIGSERIAL PRIMARY KEY,
			sub_id text,
			customer_id text,
			message_text text,
			phone_number text,
			status text) `)
	return err
}

func downMessages(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec("DROP TABLE messages")
	return err
}
