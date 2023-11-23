package postgres

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	_ "gitlab.com/thefrol/notty/internal/storage/postgres/migrations"
)

func MustConnect(dsn string) *sql.DB {

	//создаем базу
	conn, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err := conn.PingContext(context.Background()); err != nil {
		log.Fatal("не могу соединиться с базой:", err)
	}

	if err := goose.Up(conn, "."); err != nil {
		log.Fatal("не удалось мигрировать:", err)
	}

	return conn
}
