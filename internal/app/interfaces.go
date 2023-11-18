package app

import "gitlab.com/thefrol/notty/internal/entity"

type CustomerService interface {
	Create(entity.Customer) (entity.Customer, error)
	Get(string) (entity.Customer, error)
	Update(entity.Customer) (entity.Customer, error)
	Delete(string) error
}
