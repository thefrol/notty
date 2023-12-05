package config

import (
	"fmt"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	// по умолчанию вот такие параметры
	Addr string `env:"ADDR" envDefault:":8080"`
	DSN  string `env:"NOTTY_DSN"`
}

// Parse забирает переменные окружения в структуру Config,
// если конфиг не удался завершает работу программы.
//
// Используются переменные:
// NOTTY_DSN - строка соединения с БД
// ADDR - адрес создаваемого сервера, по умолчанию ":8080"
func Parse() (Config, error) {
	cfg := Config{}
	err := env.Parse(&cfg)
	if err != nil {
		return Config{}, fmt.Errorf("не удалось пропарсить переменные окружения")
	}

	if cfg.DSN == "" {
		return Config{}, fmt.Errorf("строка соединения с БД пустая. Нужно настроить переменную окружения NOTTY_DSN")
	}

	return cfg, nil
}
