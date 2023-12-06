package app

import (
	"context"
	"encoding/json"
	"net/http"

	"gitlab.com/thefrol/notty/internal/api/respond"
)

type Statistics map[string]int

func (s Statistics) ToResponseWriter(w http.ResponseWriter) { // todo не нравится!
	// не используем easyjson, потому что он не умеет в мапы
	err := json.NewEncoder(w).Encode(&s)
	if err != nil {
		respond.InternalServerError(w, "cant Marshall Stats to json: %v", err) // todo MarshallingError
	}
}

// FullStats возвращает статистику по всем сообщениям
func (a App) FullStats(ctx context.Context) (Statistics, error) {
	return a.statistics.Filter(ctx, "%", "%", "%")
}

// StatsBySubscription возвращает статистику сообщений по конкретной подписке
func (a App) StatsBySubscription(ctx context.Context, id string) (Statistics, error) {
	_, err := a.subscriptions.Get(ctx, id)
	if err != nil {
		//if err=subscriptions.ErrorNotFound
		return Statistics{}, ErrorSubscriptionNotFound
	}

	s, err := a.statistics.Filter(ctx, id, "%", "%")
	if err != nil {
		return Statistics{}, nil
	}

	return s, nil
}

// StatsByClient возвращает статистику сообщений по клиенту
func (a App) StatsByClient(ctx context.Context, id string) (Statistics, error) {
	_, err := a.customers.Get(ctx, id)
	if err != nil {
		return Statistics{}, ErrorCustomerNotFound
	}

	s, err := a.statistics.Filter(ctx, "%", id, "%")
	if err != nil {
		return Statistics{}, nil
	}

	return s, nil
}
