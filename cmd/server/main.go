package main

import (
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

	// создаем сервиси и репозитории
	cr := sqlrepo.NewCustomers(db)
	sr := sqlrepo.NewSubscriptions(db)
	stats := sqlrepo.NewStatistics(db)

	mr := sqlrepo.NewMessages(db)
	// создаем приложение
	notty := app.New(cr, sr, stats, mr)

	// создаем сервис апи
	server := api.New(notty)

	// запускаем сервак
	server.ListenAndServe(cfg.Addr)
}
