// Это инструмент для очищения базы данных, принимает в первом аргументе
// DSN базы данных, или читает из переменных окружения
package main

import (
	"os"

	"github.com/rs/zerolog/log"
	"gitlab.com/thefrol/notty/internal/storage/postgres"
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

	_, err := conn.Exec("DELETE FROM customer")
	if err != nil {
		log.Fatal().Err(err)
	}

	_, err = conn.Exec("DELETE FROM subscription")
	if err != nil {
		log.Fatal().Err(err)
	}

	log.Info().Msg("База очищена")
}
