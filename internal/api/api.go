package api

import (
	"net/http"

	"github.com/rs/zerolog"
	"gitlab.com/thefrol/notty/internal/api/generated"
	"gitlab.com/thefrol/notty/internal/app"
)

// тут мы воплощаем интерфейс для сгенерированного кода
var _ generated.ServerInterface = (*Server)(nil)

// Server представляет собой набор хендлеров для фасада нашего сервера
type Server struct {
	app    app.App // арр-аррр-аррррр ! пираты
	logger zerolog.Logger
}

// New создает новый сервис - набор хендлеров
func New(app app.App, logger zerolog.Logger) Server {
	return Server{app: app, logger: logger}
}

// Handler создает хендлер по которому будут находиться
// все маршруты нашей апишки
func (a *Server) Handler() http.Handler {
	return generated.Handler(a)
}
