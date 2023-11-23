package service

import (
	"github.com/google/uuid"

	"gitlab.com/thefrol/notty/internal/app"
	"gitlab.com/thefrol/notty/internal/app/server"
	"gitlab.com/thefrol/notty/internal/entity"
)

type CustomerRepository interface {
	Create(entity.Customer) error
	Get(string) (entity.Customer, error)
	Update(entity.Customer) error
	Delete(string) error
}

type Customers struct {
	repo CustomerRepository
}

func NewCustomers(repo CustomerRepository) Customers {
	return Customers{
		repo: repo,
	}
}

// Create implements app.CustomerService.
func (c Customers) Create(cs entity.Customer) (entity.Customer, error) {
	//todo это логика app как будто

	if cs.Id == "" {
		cs.Id = uuid.NewString()
	} else {
		_, err := c.Get(cs.Id)
		if err == nil {
			return entity.Customer{}, app.ErrorCustomerExists
		}
	}
	err := c.repo.Create(cs)
	if err != nil {
		return entity.Customer{}, err
	}
	return cs, nil
}

// Delete implements app.CustomerService.
func (c Customers) Delete(id string) error {
	_, err := c.repo.Get(id)
	if err != nil {
		// todo NotFound
		// должно быть что-то типа RepoNotFound
		return err
	}
	return c.repo.Delete(id)
}

// Get implements app.CustomerService.
func (c Customers) Get(id string) (entity.Customer, error) {
	return c.repo.Get(id)

}

// Update implements app.CustomerService.
func (c Customers) Update(cs entity.Customer) (entity.Customer, error) {
	_, err := c.repo.Get(cs.Id)
	if err != nil {
		// todo NotFound
		// должно быть что-то типа RepoNotFound
		return entity.Customer{}, err
	}
	// todo проверки на значения кастомеров если надо, может какие-то поля менять нельяз или типа того

	err = c.repo.Update(cs)
	if err != nil {
		return entity.Customer{}, err // todo Not Modified
	}

	res, err := c.Get(cs.Id)
	if err != nil {
		return entity.Customer{}, err // todo ??
	}

	return res, nil

}

var _ server.Customerere = (*Customers)(nil)
