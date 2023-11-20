package server

import "gitlab.com/thefrol/notty/internal/entity"

type CustomerService interface {
	Create(entity.Customer) (entity.Customer, error)
	Get(string) (entity.Customer, error)
	Update(entity.Customer) (entity.Customer, error)
	Delete(string) error
}

type SubscriptionService interface {
	Create(entity.Subscription) (entity.Subscription, error)
	Get(string) (entity.Subscription, error)
	Update(entity.Subscription) (entity.Subscription, error)
	Delete(string) error
}
