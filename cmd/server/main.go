package main

import (
	"gitlab.com/thefrol/notty/internal/app/config/server"
	app "gitlab.com/thefrol/notty/internal/app/server"
	service "gitlab.com/thefrol/notty/internal/storage"
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

	stats := sqlrepo.NewStatistics(db) // мы подключаем адаптер как сервис и это хорошо!

	// создаем приложение
	notty := app.New(cs, ss, stats)

	// создаем сервис апи
	server := api.New(notty)

	// запускаем сервак
	server.ListenAndServe(cfg.Addr)
}
