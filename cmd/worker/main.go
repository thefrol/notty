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
	"sync"

	"gitlab.com/thefrol/notty/internal/storage/customers"
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
	subs := subscriptions.New(db)
	cs := customers.New(db)

	active, err := subs.Active()
	if err != nil {
		log.Fatal(err)
	}

	wg := sync.WaitGroup{}
	for _, s := range active {
		s := s
		ch, err := cs.Filter(s.TagFilter, s.OperatorFilter, 300)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			continue
		}

		wg.Add(1)
		go func() {
			for c := range ch {
				fmt.Println(c, s)
			}
			wg.Done()
		}()
	}

	wg.Wait()
	log.Println("all done")
}
