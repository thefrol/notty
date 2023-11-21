package service

import (
	"github.com/google/uuid"

	"gitlab.com/thefrol/notty/internal/app"
	"gitlab.com/thefrol/notty/internal/app/server"
	"gitlab.com/thefrol/notty/internal/entity"
)

type SubscriptionRepository interface {
	Create(entity.Subscription) error
	Get(string) (entity.Subscription, error)
	Update(entity.Subscription) error
	Delete(string) error
}

type Subscriptions struct {
	repo SubscriptionRepository
}

func NewSubscriptions(repo SubscriptionRepository) Subscriptions {
	return Subscriptions{
		repo: repo,
	}
}

// Create implements app.SubscriptionService.
func (c Subscriptions) Create(cs entity.Subscription) (entity.Subscription, error) {
	//todo это логика app как будто

	if cs.Id == "" {
		cs.Id = uuid.NewString()
	} else {
		_, err := c.Get(cs.Id)
		if err == nil {
			return entity.Subscription{}, app.ErrorSubscriptionExists
		}
	}
	err := c.repo.Create(cs)
	if err != nil {
		return entity.Subscription{}, err
	}
	return cs, nil
}

// Delete implements app.SubscriptionService.
func (c Subscriptions) Delete(id string) error {
	_, err := c.repo.Get(id)
	if err != nil {
		// todo NotFound
		// должно быть что-то типа RepoNotFound
		return err
	}
	return c.repo.Delete(id)
}

// Get implements app.SubscriptionService.
func (c Subscriptions) Get(id string) (entity.Subscription, error) {
	return c.repo.Get(id)
}

// Update implements app.SubscriptionService.
func (c Subscriptions) Update(cs entity.Subscription) (entity.Subscription, error) {
	_, err := c.repo.Get(cs.Id)
	if err != nil {
		// todo NotFound
		// должно быть что-то типа RepoNotFound
		return entity.Subscription{}, err
	}
	// todo проверки на значения кастомеров если надо, может какие-то поля менять нельяз или типа того

	err = c.repo.Update(cs)
	if err != nil {
		return entity.Subscription{}, err // todo Not Modified
	}

	res, err := c.Get(cs.Id)
	if err != nil {
		return entity.Subscription{}, err // todo ??
	}

	return res, nil

}

var _ server.Subscripter = (*Subscriptions)(nil)
