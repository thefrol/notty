// содержит юз-кейс удаление клиента

package app

//todo SubscriptionRequest

func (app *App) RemoveSubscription(id string) error {

	_, err := app.subscriptions.Get(id)
	if err != nil {
		return err
	} // todo а нам нужна эта проверка? можем и обойтись думаю, можно просто проверять сколько полей было измененено

	return app.subscriptions.Delete(id)
}
