package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gitlab.com/thefrol/notty/internal/app"
	"gitlab.com/thefrol/notty/internal/config"
	"gitlab.com/thefrol/notty/internal/storage/postgres"
	"gitlab.com/thefrol/notty/internal/storage/sqlrepo"

	"gitlab.com/thefrol/notty/internal/api"
)

func main() {
	// это корневой контекст приложения
	rootContext := context.Background()

	// Создадим корневой логгер
	rootLogger := log.With().
		Str("service", "server").
		Str("instance_id", uuid.NewString()).
		Logger()

	// читаем переменные окружения
	cfg, err := config.Parse()
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

	// создаем репозитории
	cr := sqlrepo.NewCustomers(db)
	sr := sqlrepo.NewSubscriptions(db)
	stats := sqlrepo.NewStatistics(db)

	//ml := log.With().Str("repository", "messages").Logger()
	//mr := sqlrepo.NewMessages(db, ml)

	// создаем приложение
	notty := app.New(cr, sr, stats)

	// создаем веб-апи

	server := api.New(notty, rootLogger)

	// создадим контекст, который завершается при получении указанных сигналов
	ctx, stop := signal.NotifyContext(rootContext,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	defer stop()

	// запустим сервер, который нежно завершится при завершении контекта
	rootLogger.Info().
		Str("addr", cfg.Addr).
		Msg("Запускается сервер")
	if err = server.ListenAndServe(ctx, cfg.Addr, cfg.Key); err != nil {
		rootLogger.Fatal().
			Err(err).
			Msg("Не удалось запустить сервер")
	}

	// если мы оказались тут, значит сервер аккуратно завершился
	rootLogger.Info().
		Msg("main() завершается")

}
