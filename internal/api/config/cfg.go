package config

import (
	"github.com/caarlos0/env/v10"
	"github.com/rs/zerolog/log"
)

type Config struct {
	// по умолчанию вот такие параметры
	Addr string `env:"ADDR" envDefault:":8080"`
	DSN  string `env:"NOTTY_DSN"`
}

// MustParse забирает переменные окружения в структуру Config,
// если конфиг не удался завершает работу программы.
//
// Используются переменные:
// NOTTY_DSN - строка соединения с БД
// ADDR - адрес создаваемого сервера, по умолчанию ":8080"
func MustParse() Config {
	cfg := Config{}
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal().
			Str("Message", "Не удалось пропарсить переменные окружения").
			Err(err)
	}

	if cfg.DSN == "" {
		log.Info().
			Str("Message", "Строка подключения к БД - пустая. Возможно надо назначить переменную окужения NOTTY_DSN")
	}

	return cfg
}
