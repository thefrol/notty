package app

import (
	"context"

	"gitlab.com/thefrol/notty/internal/dto"
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
	All(context.Context) (dto.Statistics, error)
	Filter(ctx context.Context, subId, customerId, status string) (dto.Statistics, error)
}

type Messager interface {
	LockedSpawn(ctx context.Context, n int, status string) ([]entity.Message, error)
	ReserveFromStatus(ctx context.Context, n int, status string) ([]entity.Message, error)
	Update(context.Context, entity.Message) (entity.Message, error)
}

type Sender interface {
	Send(context.Context, entity.Message) error
}
