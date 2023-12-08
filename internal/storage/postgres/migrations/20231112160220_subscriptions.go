package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upSubscriptions, downSubscriptions)
}

// Desc           string    `json:"desc"`
// End            time.Time `json:"end"`
// Id             string    `json:"id"`
// OperatorFilter string    `json:"operator_filter"`
// PhoneFilter    string    `json:"phone_filter"`
// Start          time.Time `json:"start"`

// // Text Текст сообщения в рассылке
// Text string `json:"text"`

func upSubscriptions(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(
		`CREATE TABLE 
			subscription(
				id text PRIMARY KEY,
				msg_text text, 
				sub_start timestamptz,
				sub_end timestamptz,
				operator_filter text,
				phone_filter text,
				tag_filter text,
				description text)`)
	return err
}

func downSubscriptions(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec("DROP TABLE subscription")
	return err
}
