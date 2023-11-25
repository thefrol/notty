package app

import (
	"gitlab.com/thefrol/notty/internal/dto"
	"gitlab.com/thefrol/notty/internal/entity"
)

type Customerere interface {
	Create(entity.Customer) (entity.Customer, error)
	Get(string) (entity.Customer, error)
	Update(entity.Customer) (entity.Customer, error)
	Delete(string) error
}

type Subscripter interface {
	Create(entity.Subscription) (entity.Subscription, error)
	Get(string) (entity.Subscription, error)
	Update(entity.Subscription) (entity.Subscription, error)
	Delete(string) error
}

type Statister interface {
	All() (dto.Statistics, error)
	Filter(subId, customerId, status string) (dto.Statistics, error)
}

type Messager interface {
	Spawn(n int, status string) ([]entity.Message, error)
}
