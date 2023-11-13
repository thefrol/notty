// коллекция разных сырых запросов
package rawquery

import (
	"database/sql"
)

type Messages struct {
	db *sql.DB
}

func New(db *sql.DB) Messages {
	return Messages{
		db: db,
	}
}
