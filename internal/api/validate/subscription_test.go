// Используется для валидации данных. Запросы и ответы разделены специально
package validate

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.com/thefrol/notty/internal/entity"
)

func TestSubscriptionRequest(t *testing.T) {

	tests := []struct {
		name string
		sub  entity.Subscription
		has  error
	}{
		{
			name: "Валидная",
			sub:  ValidSub(),
			has:  nil,
		},
		{
			name: "С матом",
			sub:  ExplicitText(),
			has:  ErrorExplicitLanguage,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SubscriptionRequest(tt.sub)
			if tt.has == nil {
				assert.NoError(t, err, "не должно быть ошибки")
				return
			}
			assert.ErrorIs(t, err, tt.has)
		})
	}
}

func ValidSub() entity.Subscription {
	return entity.Subscription{
		Id:             "my-id",
		Text:           "hey! its cool",
		Desc:           "my first sub",
		Start:          time.Now(),
		End:            time.Now().Add(10 * time.Hour),
		PhoneFilter:    "+7123123", // todo valid regexp??
		OperatorFilter: "beeline",
		TagFilter:      "best",
	}
}

func ExplicitText() entity.Subscription {
	s := ValidSub()
	s.Text = "Господи ну и пиздец"
	return s
}
