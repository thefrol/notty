// содержит юз-кейс удаление клиента

package app

import "context"

//todo customerRequest

func (app *App) RemoveCustomer(ctx context.Context, id string) error {

	_, err := app.customers.Get(id)
	if err != nil {
		return err
	}

	return app.customers.Delete(id)
}
