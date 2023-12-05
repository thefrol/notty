package main

import (
	"github.com/rs/zerolog/log"
	"gitlab.com/thefrol/notty/internal/app"
	"gitlab.com/thefrol/notty/internal/app/config/server"
	"gitlab.com/thefrol/notty/internal/storage/postgres"
	"gitlab.com/thefrol/notty/internal/storage/sqlrepo"

	"gitlab.com/thefrol/notty/internal/api"
)

func main() {
	// читаем переменные окружения
	cfg := server.MustConfig()

	// соединяемся с БД
	db := postgres.MustConnect(cfg.DSN)

	// создаем репозитории
	cr := sqlrepo.NewCustomers(db)
	sr := sqlrepo.NewSubscriptions(db)
	stats := sqlrepo.NewStatistics(db)

	ml := log.With().Str("resository", "messages").Logger()
	mr := sqlrepo.NewMessages(db, ml)

	// создаем приложение
	notty := app.New(cr, sr, stats, mr)

	// создаем веб-апи
	serverLogger := log.With().Str("service", "server").Logger()
	server := api.New(notty, serverLogger)

	// запускаем сервак
	server.ListenAndServe(cfg.Addr)
}
