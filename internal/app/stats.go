package app

import (
	"context"
)

type Statistics map[string]int

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
