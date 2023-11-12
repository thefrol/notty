package main

import (
	"fmt"
	"net/http"
	"os"

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
	_ = postgres.MustConnect(dsn)

	// запускаем АПИ
	notty := api.New()

	r := chi.NewRouter()
	r.Mount("/", notty.OpenAPI())
	r.Get("/docs", api.Docs())

	http.ListenAndServe(":8080", r)
}
