package api

import (
	"net/http"

	"github.com/go-chi/chi"
)

// ListenAndServe запускает сервер Нотти, который будет выполняет
// функцию веб-фасада, описанного в openapi спеке
func (a *Server) ListenAndServe(addr string) {
	r := chi.NewRouter()

	r.Mount("/", a.OpenAPI())
	r.Get("/docs", a.Swagger())

	http.ListenAndServe(":8080", r) // todo return error
}
