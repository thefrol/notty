package app

import (
	"gitlab.com/thefrol/notty/internal/entity"
)

// CreateMessages создает n сообщений, или меньше.
func (a App) CreateMessages(n int) ([]entity.Message, error) {
	return a.Messages.Spawn(n, entity.StatusInitial)
}
