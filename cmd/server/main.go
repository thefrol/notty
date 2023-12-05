package main

import (
	"github.com/rs/zerolog/log"
	"gitlab.com/thefrol/notty/internal/api/config"
	"gitlab.com/thefrol/notty/internal/app"
	"gitlab.com/thefrol/notty/internal/storage/postgres"
	"gitlab.com/thefrol/notty/internal/storage/sqlrepo"

	"gitlab.com/thefrol/notty/internal/api"
)

func main() {
	// Создадим корневой логгер
	rootLogger := log.With().Str("service", "server").Logger()

	// читаем переменные окружения
	cfg := config.MustParse()

	// соединяемся с БД
	db := postgres.MustConnect(cfg.DSN)

	// создаем репозитории
	cr := sqlrepo.NewCustomers(db)
	sr := sqlrepo.NewSubscriptions(db)
	stats := sqlrepo.NewStatistics(db)

	ml := log.With().Str("repository", "messages").Logger()
	mr := sqlrepo.NewMessages(db, ml)

	// создаем приложение
	notty := app.New(cr, sr, stats, mr)

	// создаем веб-апи

	server := api.New(notty, rootLogger)

	// запускаем сервак
	server.ListenAndServe(cfg.Addr)
}
