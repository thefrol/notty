// Тут находится код создания нового клиента

package app

import (
	"gitlab.com/thefrol/notty/internal/entity"
)

//todo customerRequest

func (app *App) GetCustomer(id string) (entity.Customer, error) {

	return app.customers.Get(id)
}
