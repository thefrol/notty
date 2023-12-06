// Тут находится код создания нового клиента

package app

import (
	"context"

	"gitlab.com/thefrol/notty/internal/entity"
)

//todo SubscriptionRequest

func (app *App) GetSubscription(ctx context.Context, id string) (entity.Subscription, error) {

	return app.subscriptions.Get(id)
}
