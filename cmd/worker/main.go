// Воркер занимается отправкой сообщений
//
// Раз во сколько-то секунд он просыпается, и смотрит
// Есть ли активные рассылки, и если нет, то засыпает.
// Если активные рассылки есть, то создает воркеров,
// которые будут отправлять это все на эндпоинт
//
// Как бы он на кажой итерации создает такую сеть горутин.
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gitlab.com/thefrol/notty/internal/app"
	"gitlab.com/thefrol/notty/internal/storage/messages"
	"gitlab.com/thefrol/notty/internal/storage/postgres"
	"gitlab.com/thefrol/notty/internal/storage/subscriptions"
)

func main() {
	// конфигурируем
	dsn, ok := os.LookupEnv("NOTTY_DSN")
	if !ok {
		fmt.Println("Нужно передать строку подключения в переменной NOTTY_DSN")
		os.Exit(3)
	}

	// соединяемся с БД
	db := postgres.MustConnect(dsn)

	o := app.Obs{
		MessagesRepo:      messages.New(db),
		SubscriptionsRepo: subscriptions.New(db),
	}

	o.FetchAndWork(3*time.Second, 2)

	log.Println("all done")
}
