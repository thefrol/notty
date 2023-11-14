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

	"gitlab.com/thefrol/notty/internal/app"
	"gitlab.com/thefrol/notty/internal/postman"
	"gitlab.com/thefrol/notty/internal/storage/postgres"
	"gitlab.com/thefrol/notty/internal/storage/subscriptions"
	"gitlab.com/thefrol/notty/internal/stream"
)

const batchSize = 50
const timeout = 3 * time.Second
const WorkerCount = 50

const (
	retryWait  = 3
	retryCount = 3
	token      = "123123"
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

	//создаем сервисы
	MessageStreaming := stream.Messages(db)

	SubscriptionService := subscriptions.New(db)

	PostingService := postman.New("123", retryWait, retryCount, token)

	notty := app.Notifyer{
		Messages:      MessageStreaming,
		Poster:        PostingService,
		Subscriptions: SubscriptionService,
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
