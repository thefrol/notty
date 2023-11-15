package dto

// тут лежат ДТО

import (
	"github.com/google/uuid"
	"gitlab.com/thefrol/notty/internal/entity"
)

// Prototype это почти то, что может стать сообщение
// некая заявка на сообщение
type Prototype struct {
	CustomerId     string
	SubscriptionId string
	Text           string
	Phone          string
}

func (p Prototype) ToMessage() entity.Message {
	return entity.Message{
		Id:             uuid.New().String(),
		CustomerId:     p.CustomerId,
		SubscriptionId: p.SubscriptionId,
		Text:           p.Text,
		Phone:          p.Phone,
		Status:         "created",
		Sent:           nil,
	}
}
