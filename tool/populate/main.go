// Добавляет в базу некоторый осознанных набор данных с которым можно работать,
// тестировать и разрабатывать
package main

import (
	"fmt"
	"log"
	"os"

	"gitlab.com/thefrol/notty/internal/storage/postgres"
	"gitlab.com/thefrol/notty/internal/storage/sqlrepo"
	"gitlab.com/thefrol/notty/internal/storage/subscriptions"
)

func main() {
	dsn := ""

	if len(os.Args) > 1 {
		dsn = os.Args[1]
	}

	env, ok := os.LookupEnv("NOTTY_DSN") // todo выделить например в пакет notty
	if ok {
		dsn = env
	}

	if dsn == "" {
		fmt.Println("Нужно передать строку подключения в переменной NOTTY_DSN или первым параметром командной строки")
		os.Exit(3)
	}

	conn := postgres.MustConnect(dsn)

	// Добавим клиентов

	customerRepo := sqlrepo.New(conn)

	for _, c := range custs {
		err := customerRepo.Create(c)
		if err != nil {
			log.Println(err)
		}
	}
	fmt.Println("Клиенты добавлены")

	// Добавим рассылки

	SubsRepo := subscriptions.New(conn)

	for _, s := range subs {
		_, err := SubsRepo.Create(s)
		if err != nil {
			log.Println(err)
		}
	}
	fmt.Println("Рассылки добавлены")
}
