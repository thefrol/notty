package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"gitlab.com/thefrol/notty/internal/api"
)

func main() {

	notty := api.New()

	r := chi.NewRouter()
	r.Mount("/", notty.OpenAPI())
	r.Get("/docs", api.Docs())

	http.ListenAndServe(":8080", r)
}
