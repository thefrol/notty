package postgres

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog/log"
	_ "gitlab.com/thefrol/notty/internal/storage/postgres/migrations"
)

func Connect(dsn string) (*sql.DB, error) {

	//создаем базу
	conn, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("невозможно создать содениение с базой: %w", err)
	}

	if err := conn.PingContext(context.Background()); err != nil {
		return nil, fmt.Errorf("не могу соединиться с базой: %w", err)
	}

	if err := goose.Up(conn, "."); err != nil {
		return nil, fmt.Errorf("не удалось мигрировать БД: %w", err)
	}

	return conn, nil
}

func MustConnect(dsn string) *sql.DB {
	conn, err := Connect(dsn)
	if err != nil {
		log.Fatal().Err(err).Msg("немозможно подключиться к БД")
	}
	return conn
}
