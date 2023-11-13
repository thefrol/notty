package api

import (
	"net/http"

	"github.com/google/uuid"
	"gitlab.com/thefrol/notty/internal/api/decode"
	"gitlab.com/thefrol/notty/internal/api/respond"
)

// CreateSubscription implements generated.ServerInterface.
func (a *Api) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	sub, err := decode.Subscription(r)
	if err != nil {
		respond.BadRequest(w, "%v", err) // может тут оставить место только для ошибки?
		return
	}

	// todo во эта часть -логика, которая должна быть более высокого уровня, в App
	// func (a App) CreateSubscription(Subscription) Subscription

	// если айдишник не указан - создадим сами
	if sub.Id == "" {
		sub.Id = uuid.New().String()
	}

	res, err := a.app.SubscriptionRepository.Create(sub) // todo а что если такая подписка существует??
	if err != nil {
		http.Error(w, "не удалось создать подписку"+err.Error(), http.StatusInternalServerError) // todo отвечать структурой
		return
	}

	respond.Subscription(w, res)
}

// DeleteSubscription implements generated.ServerInterface.
func (a *Api) DeleteSubscription(w http.ResponseWriter, r *http.Request, id string) {
	err := a.app.SubscriptionRepository.Delete(id)
	if err != nil {
		respond.InternalServerError(w, "Не удалось удалить подписку %v", err)
	}
}

// GetSubscription implements generated.ServerInterface.
func (a *Api) GetSubscription(w http.ResponseWriter, r *http.Request, id string) {
	sub, err := a.app.SubscriptionRepository.Get(id)
	if err != nil {
		respond.InternalServerError(w, "Не удалось найти подписку с id=%s %v", id, err)
		return
	}
	respond.Subscription(w, sub)
}

// UpdateSubscription implements generated.ServerInterface.
func (a *Api) UpdateSubscription(w http.ResponseWriter, r *http.Request, id string) {
	sub, err := decode.Subscription(r)
	if err != nil {
		respond.BadRequest(w, "%v", err) // может тут оставить место только для ошибки?
		return
	}

	sub.Id = id // заменяем айдишник на тот, что стоит в запросе

	res, err := a.app.SubscriptionRepository.Update(sub) // todo а что если такая подписка существует??
	if err != nil {
		http.Error(w, "не удалось обновить подписку"+err.Error(), http.StatusInternalServerError) // todo отвечать структурой
		return
	}

	respond.Subscription(w, res)
}
