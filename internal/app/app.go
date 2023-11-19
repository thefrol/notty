package app

import (
	"database/sql"
	"errors"

	"gitlab.com/thefrol/notty/internal/stats"
	"gitlab.com/thefrol/notty/internal/storage/messages"
	"gitlab.com/thefrol/notty/internal/storage/subscriptions"
)

type App struct {
	Customers              CustomerService
	SubscriptionRepository subscriptions.Subscriptions
	MessageRepository      messages.Messages
	Statistics             stats.Service
}

func New(db *sql.DB, customers CustomerService) App {
	return App{
		Customers:              customers,
		SubscriptionRepository: subscriptions.New(db),
		MessageRepository:      messages.New(db),
		Statistics:             *stats.New(db),
	}
}

func (a App) FullStats() (stats.Statistics, error) {
	return a.Statistics.Filter("%", "%", "%")
}

func (a App) StatsBySubscription(id string) (stats.Statistics, error) {
	sub, err := a.SubscriptionRepository.Get(id)
	if err != nil {
		//if err=subscriptions.ErrorNotFound
		return stats.Statistics{}, ErrorSubscriptionNotFound
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
		return stats.Statistics{}, ErrorCustomerNotFound
	}

	s, err := a.Statistics.Filter("%", cl.Id, "%")
	if err != nil {
		return stats.Statistics{}, nil
	}

	return s, nil
}

var (
	ErrorSubscriptionNotFound = errors.New("subscription not found")
	ErrorCustomerNotFound     = errors.New("client not found")
	ErrorCustomerExists       = errors.New("client exists")
)

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
