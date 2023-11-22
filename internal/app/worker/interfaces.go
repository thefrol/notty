package worker

import "gitlab.com/thefrol/notty/internal/entity"

// Sender представляет класс, который будет смски отправлять
type Sender interface {
	Send(entity.Message) error
	Work(<-chan entity.Message) (<-chan entity.Message, error)
}
