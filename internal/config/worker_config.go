package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v10"
)

type WorkerConfig struct {

	// строка подключения к бд
	DSN string `env:"NOTTY_DSN"`

	// эндпоинт для отправки смс
	SMSEndoint string `env:"SMS_ENDPOINT"`

	// эндпоинт для отправки смс
	SMSToken string `env:"SMS_TOKEN"`

	// интервал между между поиском новых сообщений
	Interval time.Duration `env:"INTERVAL" envDefault:"5s"`

	// количество сообщений, которые мы за раз генерируем
	BatchSize uint `env:"BATCH_SIZE" envDefault:"50"`

	// размер пула воркеров, которые отправляют сообщения
	WorkerCount uint `env:"WORKER_COUNT" envDefault:"20"`

	// интервал между попытками отправить смс
	RetryInterval time.Duration `env:"RETRY_INTERVAL" envDefault:"3s"`

	// количество повторных попыток отправить смс
	RetryCount uint `env:"RETRY_COUNT" envDefault:"3"`
}

// ForWorker забирает переменные окружения в структуру Config,
// если конфиг не удался завершает работу программы.
func ForWorker() (WorkerConfig, error) {
	cfg := WorkerConfig{}
	err := env.Parse(&cfg)
	if err != nil {
		return WorkerConfig{}, fmt.Errorf("не удалось пропарсить переменные окружения")
	}

	if cfg.DSN == "" {
		return WorkerConfig{}, fmt.Errorf("строка соединения с БД пустая. Нужно настроить переменную окружения NOTTY_DSN")
	}

	if cfg.SMSEndoint == "" {
		return WorkerConfig{}, fmt.Errorf("нет смс эндпоинта в SMS_ENDPOINT")
	}

	if cfg.WorkerCount == 0 {
		return WorkerConfig{}, fmt.Errorf("количество рабочих должно быть больше 0 в WORKER_POOL должено быть больше нуля")
	}

	return cfg, nil
}
