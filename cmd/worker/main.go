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
	"os"
	"sync"
	"time"

	"gitlab.com/thefrol/notty/internal/app/worker"
	"gitlab.com/thefrol/notty/internal/sms"
	service "gitlab.com/thefrol/notty/internal/storage"
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
	// конфигурируем
	dsn, ok := os.LookupEnv("NOTTY_DSN")
	if !ok {
		fmt.Println("Нужно передать строку подключения в переменной NOTTY_DSN")
		os.Exit(3)
	}

	// соединяемся с БД
	db := postgres.MustConnect(dsn)

	//создаем сервисы
	mr := sqlrepo.NewMessages(db)
	MessageStreaming := service.Messages(db, mr)

	//SubscriptionRepo := sqlrepo.NewSubscriptions(db)

	PostingService := sms.New(endpoint, retryWait, retryCount, token)

	notty := worker.Notifyer{
		Messages: MessageStreaming,
		Poster:   PostingService,
	}
	//и за дело

	for {
		// спим какое-то время. Тут бы заменить по сути на тикер
		time.Sleep(timeout)
		wg := sync.WaitGroup{}

		// ищем новые сообщения
		wg.Add(1)
		go func() {
			notty.FindAndSend(batchSize, WorkerCount)
			wg.Done()
		}()

		// и одновременно пытаемся отправить неотправленные, но
		// могущие быть отправленными

		wg.Add(1)
		go func() {
			notty.TryToResend(batchSize, WorkerCount)
			wg.Done()
		}()

		// ждем когда оба процесса закончатся
		// todo а вообще надо их в отдельных приложениях например
		// или в разных горутинах
		wg.Wait()
	}
}
