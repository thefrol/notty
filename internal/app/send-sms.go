package app

import "gitlab.com/thefrol/notty/internal/entity"

type SMSer interface {
	Send(entity.Message) error
	//Work(<-chan entity.Message) (<-chan entity.Message, error)
}

// Notifyer это попытка понять, что же мы делать будем в этом слое
// можем ли мы раздробить слой приложения на куски так, чтобы
// логика не была разной для разных случаев
type Notifyer struct {
	proxy SMSer
}

func NewNotifyer(proxy SMSer, messages Messager) Notifyer {
	return Notifyer{
		proxy: proxy,
	}
}

func (n Notifyer) SendNotification(m entity.Message) error {
	return n.proxy.Send(m)
}

// тут мы делаем какие-то совсем простые функции, чтобы потом
// уже логику на каналах реализовать в юз-кейсах
