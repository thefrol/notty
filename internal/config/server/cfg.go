package server

import (
	"log"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	// по умолчанию вот такие параметры
	Addr string `env:"ADDR" envDefault:":8080"`
	DSN  string `env:"NOTTY_DSN"`
}

// MustConfig забирает переменные окружения в структуру Config,
// если конфиг не удался завершает работу программы.
//
// Используются переменные:
// NOTTY_DSN - строка соединения с БД
// ADDR - адрес создаваемого сервера, по умолчанию ":8080"
func MustConfig() Config {
	cfg := Config{}
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatalf("Не удалось пропарсить переменные окружения: %v", err)
	}

	if cfg.DSN == "" {
		log.Println("Строка подключения к БД - пустая. Возможно надо назначить переменную окужения NOTTY_DSN")
	}

	return cfg
}
