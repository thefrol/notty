package notifyloop

import (
	"context"

	"gitlab.com/thefrol/notty/internal/entity"
)

// Notifyer это кусок такой высокоровневой,
// логики над репозиторием и смсками, чтобы убрать с
// воркера некоторые задачки по перемещению данных между слоями
//
// Тут уже инкапсулирована логика работы с состояними. Да,
// конечно нужны ли мне состояния отдельный вопрос, но сейчас от них отказываться
// слишком стремно.
//
// Просто не хочется делать гигантской структуру воркера. Там ещё всякие тайминки,
// контесты и прочее. А касаемо работы с инфраструктурой что ли
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

// CreateMessages создает n сообщений, или меньше. Сообщения до этого не существовали
func (notty Notifyer) CreateMessages(ctx context.Context, n int) ([]entity.Message, error) {
	return notty.messages.LockedSpawn(ctx, n, entity.StatusInitial)
}

// ReserveFailed резервирует n сообщений с изначальным статусом
// fromStatus
func (notty Notifyer) ReserveFailed(ctx context.Context, n int) ([]entity.Message, error) {
	return notty.messages.ReserveFromStatus(ctx, n, entity.StatusFailed)
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
