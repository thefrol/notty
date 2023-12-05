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
	"context"
	"os"
	"time"

	"github.com/rs/zerolog/log"
	"gitlab.com/thefrol/notty/internal/app"
	"gitlab.com/thefrol/notty/internal/service"
	"gitlab.com/thefrol/notty/internal/sms"
	"gitlab.com/thefrol/notty/internal/storage/postgres"
	"gitlab.com/thefrol/notty/internal/storage/sqlrepo"
)

const batchSize = 50
const timeout = 3 * time.Second
const WorkerCount = 50

const (
	retryWait  = 3
	retryCount = 3
	endpoint   = "https://probe.fbrq.cloud/v1/send/"
)

var token = os.Getenv("ENDPOINT_TOKEN")

func main() {
	// создадим корневой логгер
	rootLogger := log.With().Str("service", "worker").Logger()

	// конфигурируем
	dsn, ok := os.LookupEnv("NOTTY_DSN")
	if !ok {
		log.Info().
			Str("Message", "Неправильная конфигурации. Нужно передать строку подключения в переменной NOTTY_DSN")
		os.Exit(3)
	}

	// соединяемся с БД
	db := postgres.MustConnect(dsn)

	//создаем сервисы
	mr := sqlrepo.NewMessages(db, rootLogger)
	sms := sms.NewEndpoint(endpoint, retryWait, retryCount, token)

	notty := app.NewNotifyerrrr(mr, sms)

	// создаем приложение
	worker := service.Worker{
		Notifyer:  notty,
		Timeout:   timeout,
		BatchSize: batchSize,
		Logger:    rootLogger,
	}

	// Запускаем воркера на постоянку)
	worker.FetchAndSend(context.TODO())

}
