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
	"context"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gitlab.com/thefrol/notty/internal/config"
	"gitlab.com/thefrol/notty/internal/notifyloop"
	"gitlab.com/thefrol/notty/internal/notifyloop/fabrique"
	"gitlab.com/thefrol/notty/internal/storage/postgres"
	"gitlab.com/thefrol/notty/internal/storage/sqlrepo"
)

func main() {
	// это корневой контекст приложения
	rootContext := context.Background()

	// создадим корневой логгер
	rootLogger := log.With().
		Str("service", "worker").
		Str("instance_id", uuid.NewString()).
		Logger()

	// читаем переменные окружения
	cfg, err := config.ForWorker()
	if err != nil {
		rootLogger.Fatal().
			Err(err).
			Msg("Ошибка конфигурации")
	}

	// соединяемся с БД
	db, err := postgres.Connect(cfg.DSN)
	if err != nil {
		rootLogger.Fatal().
			Err(err).
			Msg("Не удалось подключить к базе данных")
	}

	//создаем репозитории/адаптеры
	mr := sqlrepo.NewMessages(db, rootLogger)
	sender := fabrique.NewEndpoint(
		cfg.SMSEndoint,
		int(cfg.RetryInterval),
		int(cfg.RetryCount),
		cfg.SMSToken)

	// создаем приложение
	notty := notifyloop.NewNotifyer(mr, sender)

	worker := notifyloop.Worker{
		Notifyer:  notty,
		Timeout:   cfg.Interval,
		BatchSize: int(cfg.BatchSize),
		Logger:    rootLogger,
	}

	// создадим контекст, который завершается при получении указанных сигналов
	ctx, stop := signal.NotifyContext(rootContext,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	defer stop()

	// Запускаем воркера на постоянку)
	worker.FetchAndSend(ctx)

	rootLogger.Info().
		Msg("окончание main()")

}
