package notifyloop

import (
	"context"

	"gitlab.com/thefrol/notty/internal/entity"
)

// Notifyer это уточнение над репозиторием сообщений, который умеет хорошо
// генерировать сообщения
type Notifyer struct {
	messages Messager
	sender   Sender
}

func NewNotifyer(messages Messager, sender Sender) *Notifyer {
	return &Notifyer{
		messages: messages,
		sender:   sender,
	}
}

// CreateMessages создает n сообщений, или меньше.
func (notty Notifyer) CreateMessages(ctx context.Context, n int) ([]entity.Message, error) {
	return notty.messages.LockedSpawn(ctx, n, entity.StatusInitial)
}

// ReserveMessages резервирует n сообщений с изначальным статусом
// fromStatus
func (notty Notifyer) ReserveMessages(ctx context.Context, n int, fromStatus string) ([]entity.Message, error) {
	return notty.messages.ReserveFromStatus(ctx, n, fromStatus)
}

// SendSMS отправляет смской сообщение m
func (notty Notifyer) SendSMS(ctx context.Context, m entity.Message) error {
	return notty.sender.Send(ctx, m)
}

// UpdateMessage обновляеят статус у сообщения
func (notty Notifyer) UpdateMessage(ctx context.Context, m entity.Message) error {
	_, err := notty.messages.Update(ctx, m)
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
