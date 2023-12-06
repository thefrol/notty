// Тут находится код создания нового клиента

package app

import (
	"context"

	"gitlab.com/thefrol/notty/internal/entity"
)

//todo customerRequest

func (app *App) GetCustomer(ctx context.Context, id string) (entity.Customer, error) {

	return app.customers.Get(ctx, id)
}
