package app

type App struct {
	customers     Customerere
	subscriptions Subscripter
	statistics    Statister
	messages      Messager
}

func New(customers Customerere, subscriptions Subscripter, stats Statister, messages Messager) App {
	return App{
		customers:     customers,
		subscriptions: subscriptions,
		statistics:    stats,
		messages:      messages,
	}
}

// todo
//
// в App должна лежать вся высокоуровневая логика,
// типа, какие рассылки мы должны обновлять и как
// может быть например, стоит проверить, что
// в рассылке уже есть
//
// Например, запретить обновлять время начала для рассылки
// Вообще у рассылки должно быть что-то вроде Put
// И тут тожно нужны какие-то Nullable типы как будето
//
// Но вернемся к логике, например, у нас есть условие,
// нельзя обновлять уже завершенные рассылки,
// или нельзя менять время начала на более раннее, если
// это время в прошлом
//
// ИЛи например, нельзя матерные слова в тексте, и фильтр
// или проверка текста на фейлы, странные символы. Например,
// запретить какие-нить кракозябры или ссылки, ну кто знает
// все это может быть тут
//
// Это вроде и валидация, просто высокоуровневая
