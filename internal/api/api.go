package api

import (
	"net/http"

	"gitlab.com/thefrol/notty/internal/api/generated"
)

// Api представляет собой набор хендлеров для фасада нашего сервера
type Api struct {
}

// New создает новый сервис - набор хендлеров
func New() Api {
	return Api{}
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

// CreateClient implements generated.ServerInterface.
func (*Api) CreateClient(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

// CreateSubscription implements generated.ServerInterface.
func (*Api) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

// DeleteClient implements generated.ServerInterface.
func (*Api) DeleteClient(w http.ResponseWriter, r *http.Request, id string) {
	panic("unimplemented")
}

// DeleteSubscription implements generated.ServerInterface.
func (*Api) DeleteSubscription(w http.ResponseWriter, r *http.Request, id string) {
	panic("unimplemented")
}

// GetClient implements generated.ServerInterface.
func (*Api) GetClient(w http.ResponseWriter, r *http.Request, id string) {
	panic("unimplemented")
}

// GetSubscription implements generated.ServerInterface.
func (*Api) GetSubscription(w http.ResponseWriter, r *http.Request, id string) {
	panic("unimplemented")
}

// UpdateClient implements generated.ServerInterface.
func (*Api) UpdateClient(w http.ResponseWriter, r *http.Request, id string) {
	panic("unimplemented")
}

// UpdateSubscription implements generated.ServerInterface.
func (*Api) UpdateSubscription(w http.ResponseWriter, r *http.Request, id string) {
	panic("unimplemented")
}

var _ generated.ServerInterface = (*Api)(nil)
