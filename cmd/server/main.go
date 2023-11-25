package main

import (
	"gitlab.com/thefrol/notty/internal/app"
	"gitlab.com/thefrol/notty/internal/app/config/server"
	"gitlab.com/thefrol/notty/internal/storage"
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
	cs := storage.NewCustomers(cr)

	sr := sqlrepo.NewSubscriptions(db)
	ss := storage.NewSubscriptions(sr)

	mr := sqlrepo.NewMessages(db)
	stats := sqlrepo.NewStatistics(db) // мы подключаем адаптер как сервис и это хорошо!

	// создаем приложение
	notty := app.New(cs, ss, stats, mr)

	// создаем сервис апи
	server := api.New(notty)

	// запускаем сервак
	server.ListenAndServe(cfg.Addr)
}
