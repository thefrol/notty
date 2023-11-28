// Тут находится код создания нового клиента

package app

import (
	"gitlab.com/thefrol/notty/internal/entity"
)

//todo SubscriptionRequest

func (app *App) UpdateSubscription(c entity.Subscription) (entity.Subscription, error) {
	if err := c.Validate(); err != nil {
		return entity.Subscription{}, err
	}

	// не уверен, что мне нужна эта проверка. Он же инсерт делать не будет //todo
	_, err := app.subscriptions.Get(c.Id)
	if err != nil {
		// todo NotFound
		// должно быть что-то типа RepoNotFound
		return entity.Subscription{}, err
	}
	// todo проверки на значения подписок если надо, может какие-то поля менять нельяз или типа того

	res, err := app.subscriptions.Update(c)
	if err != nil {
		return entity.Subscription{}, err // todo Not Modified
	}

	return res, nil
}
