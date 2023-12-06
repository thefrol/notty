package app

import (
	"context"

	"gitlab.com/thefrol/notty/internal/entity"
)

// Notifyerrrr представляет слой приложения
// для отправщика. Тут функции генерации сообщения
// и отправки сообщений
type Notifyerrrr struct {
	messages Messager
	sender   Sender
}

func NewNotifyerrrr(messages Messager, sender Sender) *Notifyerrrr {
	return &Notifyerrrr{
		messages: messages,
		sender:   sender,
	}
}

// CreateMessages создает n сообщений, или меньше.
func (notty Notifyerrrr) CreateMessages(ctx context.Context, n int) ([]entity.Message, error) {
	return notty.messages.LockedSpawn(n, entity.StatusInitial)
}

// ReserveMessages резервирует n сообщений с изначальным статусом
// fromStatus
func (notty Notifyerrrr) ReserveMessages(ctx context.Context, n int, fromStatus string) ([]entity.Message, error) {
	return notty.messages.ReserveFromStatus(n, fromStatus)
}

// SendSMS отправляет смской сообщение m
func (notty Notifyerrrr) SendSMS(ctx context.Context, m entity.Message) error {
	return notty.sender.Send(m)
}

// UpdateMessage обновляеят статус у сообщения
func (notty Notifyerrrr) UpdateMessage(ctx context.Context, m entity.Message) error {
	_, err := notty.messages.Update(m)
	if err != nil {
		return err
	}
	return nil
}

// todo
//
// смотрю на это все и понимаю, что репозиторий
// мог бы работать с совершенно урезанной структурой данных
// типа телефон, сообщение, статус... На вход по крайней мере
