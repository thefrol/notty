package api

import (
	"net/http"

	"gitlab.com/thefrol/notty/internal/api/generated"
	"gitlab.com/thefrol/notty/internal/api/respond"
	"gitlab.com/thefrol/notty/internal/app"
)

// Server представляет собой набор хендлеров для фасада нашего сервера
type Server struct {
	app app.App // арр-аррр-аррррр ! пираты
}

// GetMessage implements generated.ServerInterface.
func (a *Server) GetMessage(w http.ResponseWriter, r *http.Request, id string) {
	respond.InternalServerError(w, "Не реализовано :(")
}

// New создает новый сервис - набор хендлеров
func New(app app.App) Server {
	return Server{app: app}
}

// OpenAPI создает хендлер по которому будут находиться
// все маршруты нашей апишки
func (a *Server) OpenAPI() http.Handler {
	return generated.Handler(a)
}

var _ generated.ServerInterface = (*Server)(nil)
