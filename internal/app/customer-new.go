// Тут находится код создания нового клиента

package app

import (
	"github.com/google/uuid"
	"gitlab.com/thefrol/notty/internal/entity"
)

//todo customerRequest

func (app *App) NewCustomer(c entity.Customer) (entity.Customer, error) {
	// добавить uuid если не задан, если задан
	// то проверить существует ли такой челик в базе
	//
	// по сути с такой логикой работы мы можем избавиться от
	// create в репозитории оставить только Put как бы
	if c.Id == "" {
		c.Id = uuid.NewString()
	} else {
		_, err := app.customers.Get(c.Id)
		if err == nil {
			return entity.Customer{}, ErrorCustomerExists
		}
	}

	if err := c.Validate(); err != nil {
		return entity.Customer{}, err
	}

	return app.customers.Create(c)
}
