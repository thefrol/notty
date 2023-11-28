// Тут находится код создания нового клиента

package app

import (
	"gitlab.com/thefrol/notty/internal/entity"
)

//todo SubscriptionRequest

func (app *App) GetSubscription(id string) (entity.Subscription, error) {

	return app.subscriptions.Get(id)
}
