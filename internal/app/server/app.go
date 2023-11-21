package server

import (
	"gitlab.com/thefrol/notty/internal/app"
	"gitlab.com/thefrol/notty/internal/dto"
)

type App struct {
	Customers     Customerere
	Subscriptions Subscripter
	Statistics    Statister
}

func New(customers Customerere, subscriptions Subscripter, stats Statister) App {
	return App{
		Customers:     customers,
		Subscriptions: subscriptions,
		Statistics:    stats,
	}
}

// для упрощения логики статистики можно было бы выделить вот это копание в сервис,
// опустить на уровень пониже

// FullStats возвращает статистику по всем сообщениям
func (a App) FullStats() (dto.Statistics, error) {
	return a.Statistics.Filter("%", "%", "%")
}

// StatsBySubscription возвращает статистику сообщений по конкретной подписке
func (a App) StatsBySubscription(id string) (dto.Statistics, error) {
	_, err := a.Subscriptions.Get(id)
	if err != nil {
		//if err=subscriptions.ErrorNotFound
		return dto.Statistics{}, app.ErrorSubscriptionNotFound
	}

	s, err := a.Statistics.Filter(id, "%", "%")
	if err != nil {
		return dto.Statistics{}, nil
	}

	return s, nil
}

// StatsByClient возвращает статистику сообщений по клиенту
func (a App) StatsByClient(id string) (dto.Statistics, error) {
	_, err := a.Customers.Get(id)
	if err != nil {
		return dto.Statistics{}, app.ErrorCustomerNotFound
	}

	s, err := a.Statistics.Filter("%", id, "%")
	if err != nil {
		return dto.Statistics{}, nil
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
