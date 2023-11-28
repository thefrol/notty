// содержит юз-кейс удаление клиента

package app

//todo customerRequest

func (app *App) RemoveCustomer(id string) error {

	_, err := app.customers.Get(id)
	if err != nil {
		return err
	}

	return app.customers.Delete(id)
}
