package notifyloop

import (
	"context"

	"gitlab.com/thefrol/notty/internal/entity"
)

type Messager interface {
	LockedSpawn(ctx context.Context, n int, status string) ([]entity.Message, error)
	ReserveFromStatus(ctx context.Context, n int, status string) ([]entity.Message, error)
	Update(context.Context, entity.Message) (entity.Message, error)
}

type Sender interface {
	Send(context.Context, entity.Message) error
}
