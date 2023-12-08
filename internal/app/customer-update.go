// Тут находится код создания нового клиента

package app

import (
	"context"

	"gitlab.com/thefrol/notty/internal/entity"
)

//todo customerRequest

func (app *App) UpdateCustomer(ctx context.Context, c entity.Customer) (entity.Customer, error) {
	if err := c.Validate(); err != nil {
		return entity.Customer{}, err
	}

	// не уверен, что мне нужна эта проверка. Он же инсерт делать не будет //todo
	_, err := app.customers.Get(ctx, c.Id)
	if err != nil {
		// todo NotFound
		// должно быть что-то типа RepoNotFound
		return entity.Customer{}, err
	}
	// todo проверки на значения кастомеров если надо, может какие-то поля менять нельяз или типа того

	res, err := app.customers.Update(ctx, c)
	if err != nil {
		return entity.Customer{}, err // todo Not Modified
	}

	return res, nil
}
