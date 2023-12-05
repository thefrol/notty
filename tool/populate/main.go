// Добавляет в базу некоторый осознанных набор данных с которым можно работать,
// тестировать и разрабатывать
package main

import (
	"os"

	"github.com/rs/zerolog/log"
	"gitlab.com/thefrol/notty/internal/storage/postgres"
	"gitlab.com/thefrol/notty/internal/storage/sqlrepo"
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
		log.Fatal().
			Str("Message", "Ошибка конфигурации").
			Str("Description", "Нужно передать строку подключения в переменной NOTTY_DSN или первым параметром командной строки")
		os.Exit(3)
	}

	conn := postgres.MustConnect(dsn)

	// Добавим клиентов

	customerRepo := sqlrepo.NewCustomers(conn)

	for _, c := range custs {
		_, err := customerRepo.Create(c)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Ошибка добавления клиента")
		}
	}
	log.Info().Msg("Клиенты добавлены")

	// Добавим рассылки

	SubsRepo := sqlrepo.NewSubscriptions(conn)

	for _, s := range subs {
		_, err := SubsRepo.Create(s)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Ошибка добавления подписки")
		}
	}
	log.Info().Msg("Рассылки добавлены")
}
