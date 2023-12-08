package config

import (
	"fmt"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	// адрес, на котором сервер будет слушать
	Addr string `env:"ADDR" envDefault:":8080"`

	// строка подключения к бд
	DSN string `env:"NOTTY_DSN"`

	// ключ для подписывания jwt
	Key string `env:"KEY"`
}

// ForAPI забирает переменные окружения в структуру Config,
// если конфиг не удался завершает работу программы.
func ForAPI() (Config, error) {
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
