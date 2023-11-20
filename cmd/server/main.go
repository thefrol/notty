package main

import (
	"gitlab.com/thefrol/notty/internal/app"
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
	cr := sqlrepo.New(db)
	cs := service.NewCustomers(cr)

	// создаем приложение
	notty := app.New(db, cs)

	// создаем сервис апи
	server := api.New(notty)

	// запускаем сервак
	server.ListenAndServe(cfg.Addr)
}
