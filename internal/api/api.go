package api

import (
	"net/http"

	"gitlab.com/thefrol/notty/internal/api/generated"
	"gitlab.com/thefrol/notty/internal/app"
)

// Api представляет собой набор хендлеров для фасада нашего сервера
type Api struct {
	app app.App // арр-аррр-аррррр ! пираты
}

// GetMessage implements generated.ServerInterface.
func (a *Api) GetMessage(w http.ResponseWriter, r *http.Request, id string) {
	panic("not implemented")
}

// New создает новый сервис - набор хендлеров
func New(app app.App) Api {
	return Api{app: app}
}

// OpenAPI создает хендлер по которому будут находиться
// все маршруты нашей апишки
func (a *Api) OpenAPI() http.Handler {
	return generated.Handler(a)
}

// Swagger возвращает ручку, на которой держится сваггер текущей спецификации
// если доку не удалось создать, то там будет выводиться ошибка, но всегда
// эта ручка существует
func (a *Api) Swagger() http.HandlerFunc {
	return Docs()
}

var _ generated.ServerInterface = (*Api)(nil)
