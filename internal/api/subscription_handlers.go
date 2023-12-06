package api

import (
	"errors"
	"net/http"

	"gitlab.com/thefrol/notty/internal/api/decode"
	"gitlab.com/thefrol/notty/internal/api/respond"
	"gitlab.com/thefrol/notty/internal/app"
	"gitlab.com/thefrol/notty/internal/entity/valid"
)

// CreateSubscription implements generated.ServerInterface.
func (a *Server) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	c, err := decode.Subscription(r)
	if err != nil {
		respond.BadRequest(w, "%v", err) // может тут оставить место только для ошибки?
		return
	}

	res, err := a.app.NewSubscription(r.Context(), c) // todo а что если такой клиент существует??
	if err != nil {
		if errors.Is(err, app.ErrorSubscriptionExists) {
			respond.Errorf(w, http.StatusConflict, "Рассылка с id %s существует ", c.Id)
		}
		respond.InternalServerError(w, "Неизвестная ошибка %s", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	respond.Subscription(w, res)
}

// DeleteSubscription implements generated.ServerInterface.
func (a *Server) DeleteSubscription(w http.ResponseWriter, r *http.Request, id string) {
	if err := valid.Id(id); err != nil {
		respond.BadRequest(w, "%v", err)
		return
	}

	err := a.app.RemoveSubscription(r.Context(), id)
	if err != nil {
		if errors.Is(err, app.ErrorSubscriptionNotFound) {
			respond.NotFound(w, "Рассылка с id %s не обнаружена", id)
			return
		}
		respond.InternalServerError(w, "Не удалось удалить рассылку %v", err)
	}
}

// GetSubscription implements generated.ServerInterface.
func (a *Server) GetSubscription(w http.ResponseWriter, r *http.Request, id string) {
	if err := valid.Id(id); err != nil {
		respond.BadRequest(w, "%v", err)
		return
	}

	sub, err := a.app.GetSubscription(r.Context(), id)
	if err != nil {
		if errors.Is(err, app.ErrorSubscriptionNotFound) {
			respond.NotFound(w, "Рассылка с id %s не обнаружена", id)
			return
		}
		respond.InternalServerError(w, "Не удалось найти рассылку с id=%s %v", id, err)
		return
	}
	respond.Subscription(w, sub)
}

// UpdateSubscription implements generated.ServerInterface.
func (a *Server) UpdateSubscription(w http.ResponseWriter, r *http.Request, id string) {
	c, err := decode.Subscription(r)
	if err != nil {
		respond.BadRequest(w, "%v", err) // может тут оставить место только для ошибки?
		return
	}

	c.Id = id // заменяем айдишник на тот, что стоит в запросе
	// bug это что такое вообще!!!

	res, err := a.app.UpdateSubscription(r.Context(), c) // todo а что если такой клиент существует??
	if err != nil {
		if errors.Is(err, app.ErrorSubscriptionNotFound) {
			respond.NotFound(w, "Рассылка с id %s не обнаружена", id)
			return
		}
		http.Error(w, "не удалось обновить рассылку"+err.Error(), http.StatusInternalServerError) // todo отвечать структурой
		return
	}

	respond.Subscription(w, res)
}

// SubscriptionStats implements generated.ServerInterface.
func (a *Server) SubscriptionStats(w http.ResponseWriter, r *http.Request, id string) {
	// редиректим в другую функцию
	a.StatsBySubscriptionId(w, r, id)
}
