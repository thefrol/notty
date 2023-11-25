package sms

import "gitlab.com/thefrol/notty/internal/entity"

// SMSProxy представляет класс, который будет смски отправлять
type SMSProxy interface {
	Send(entity.Message) error
	Work(<-chan entity.Message) (<-chan entity.Message, error)
}
