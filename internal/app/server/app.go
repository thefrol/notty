package server

import (
	"database/sql"

	"gitlab.com/thefrol/notty/internal/app"
	"gitlab.com/thefrol/notty/internal/stats"
	"gitlab.com/thefrol/notty/internal/storage/messages"
)

type App struct {
	Customers         CustomerService
	Subscriptions     SubscriptionService
	MessageRepository messages.Messages
	Statistics        stats.Service
}

func New(db *sql.DB, customers CustomerService, subscriptions SubscriptionService) App {
	return App{
		Customers:         customers,
		Subscriptions:     subscriptions,
		MessageRepository: messages.New(db),
		Statistics:        *stats.New(db),
	}
}

func (a App) FullStats() (stats.Statistics, error) {
	return a.Statistics.Filter("%", "%", "%")
}

func (a App) StatsBySubscription(id string) (stats.Statistics, error) {
	sub, err := a.Subscriptions.Get(id)
	if err != nil {
		//if err=subscriptions.ErrorNotFound
		return stats.Statistics{}, app.ErrorSubscriptionNotFound
	}

	s, err := a.Statistics.Filter(sub.Id, "%", "%")
	if err != nil {
		return stats.Statistics{}, nil
	}

	return s, nil
}

func (a App) StatsByClient(id string) (stats.Statistics, error) {
	cl, err := a.Customers.Get(id)
	if err != nil {
		return stats.Statistics{}, app.ErrorCustomerNotFound
	}

	s, err := a.Statistics.Filter("%", cl.Id, "%")
	if err != nil {
		return stats.Statistics{}, nil
	}

	return s, nil
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
