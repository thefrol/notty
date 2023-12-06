// Тут находится код создания нового клиента

package app

import (
	"context"

	"github.com/google/uuid"
	"gitlab.com/thefrol/notty/internal/entity"
)

//todo SubscriptionRequest

func (app *App) NewSubscription(ctx context.Context, s entity.Subscription) (entity.Subscription, error) {
	// добавить uuid если не задан, если задан
	// то проверить существует ли такой челик в базе
	//
	// по сути с такой логикой работы мы можем избавиться от
	// create в репозитории оставить только Put как бы
	if s.Id == "" {
		s.Id = uuid.NewString()
	} else {
		_, err := app.subscriptions.Get(ctx, s.Id)
		if err == nil {
			return entity.Subscription{}, ErrorSubscriptionExists
		}
	}

	if err := s.Validate(); err != nil {
		return entity.Subscription{}, err
	}

	return app.subscriptions.Create(ctx, s)
}
