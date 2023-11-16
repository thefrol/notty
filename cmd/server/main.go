package main

import (
	"fmt"
	"net/http"
	"os"

	"gitlab.com/thefrol/notty/internal/app"
	"gitlab.com/thefrol/notty/internal/storage/postgres"

	"github.com/go-chi/chi"
	"gitlab.com/thefrol/notty/internal/api"
)

func main() {
	// конфигурируем
	dsn, ok := os.LookupEnv("NOTTY_DSN")
	if !ok {
		fmt.Println("Нужно передать строку подключения в переменной NOTTY_DSN")
		os.Exit(3)
	}

	// соединяемся с БД
	db := postgres.MustConnect(dsn)

	// создаем приложение
	app := app.New(db)

	// создаем сервис апи
	notty := api.New(app)

	r := chi.NewRouter()

	r.Mount("/", notty.OpenAPI())
	r.Get("/docs", api.Docs())
	r.Get("/ui", api.SwaggerUI("Нотти!"))

	http.ListenAndServe(":8080", r)
}
