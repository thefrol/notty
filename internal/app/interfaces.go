package app

import (
	"context"

	"gitlab.com/thefrol/notty/internal/entity"
)

type Customerere interface {
	Create(context.Context, entity.Customer) (entity.Customer, error)
	Get(context.Context, string) (entity.Customer, error)
	Update(context.Context, entity.Customer) (entity.Customer, error)
	Delete(context.Context, string) error
}

type Subscripter interface {
	Create(context.Context, entity.Subscription) (entity.Subscription, error)
	Get(context.Context, string) (entity.Subscription, error)
	Update(context.Context, entity.Subscription) (entity.Subscription, error)
	Delete(context.Context, string) error
}

type Statister interface {
	All(context.Context) (Statistics, error)
	Filter(ctx context.Context, subId, customerId, status string) (Statistics, error)
}
