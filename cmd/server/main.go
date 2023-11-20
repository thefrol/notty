package main

import (
	app "gitlab.com/thefrol/notty/internal/app/server"
	"gitlab.com/thefrol/notty/internal/config/server"
	"gitlab.com/thefrol/notty/internal/service"
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
	cs := service.NewCustomers(cr)

	sr := sqlrepo.NewSubscriptions(db)
	ss := service.NewSubscriptions(sr)

	// создаем приложение
	notty := app.New(db, cs, ss)

	// создаем сервис апи
	server := api.New(notty)

	// запускаем сервак
	server.ListenAndServe(cfg.Addr)
}
