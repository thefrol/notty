package api

import (
	"context"
	"net/http"
	"time"

	"github.com/Lavalier/zchi"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth/v5"
)

// ListenAndServe запускает сервер Нотти, который будет выполняет
// функцию веб-фасада, описанного в openapi спеке
//
// При завершении контекста берет десять секунд за закрытие соединений
// и потом возвращает управление. Вернет err==nil если завершение прошло хорошо
// , и err!=nil если проблемы с запуском сервера
func (a *Server) ListenAndServe(ctx context.Context, addr, key string) error {
	r := chi.NewRouter()

	// логирование запросов
	r.Use(zchi.Logger(a.logger))

	// добавим сваггер
	r.Get("/docs", a.Swagger())

	// и основное ап, закрытое токеном
	r.Group(func(r chi.Router) {
		// если указан ключ, то закрываем доступ
		if key != "" {
			a.logger.Info().Msg("Апи закрыто авторизацией через jwt")
			tokenAuth := jwtauth.New("HS256", []byte(key), nil)
			r.Use(jwtauth.Verifier(tokenAuth))
		} else {
			a.logger.Error().Msg("Апи открыто для доступа")
		}

		r.Mount("/", a.OpenAPI())
	})

	server := &http.Server{Addr: addr, Handler: r}

	// Запустим горутину, которая будет окончания контекса и
	// после этого остановит сервер
	stop := make(chan struct{})

	go func() {
		<-ctx.Done()
		a.logger.Error().
			Msg("Начинается остановка сервера")

		// выключаем сервер, дадим ему 10 секунд,
		// что вообще по-хорошему должно быть в конфиге
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		err := server.Shutdown(ctx)
		if err != nil {
			a.logger.Error().
				Err(err).
				Msg("Ошибка при остановке сервера")
		}

		// закрывем канал остановки, там самым дадим понять, что
		// сервер успешно остановлен
		close(stop)

	}()

	// запускаем сервер и занимаем поток
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	<-stop
	a.logger.Info().Msg("Сервер остановлен успешно")
	return nil
}
