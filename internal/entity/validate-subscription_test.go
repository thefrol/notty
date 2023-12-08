// Используется для валидации данных. Запросы и ответы разделены специально
package entity_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.com/thefrol/notty/internal/entity"
	"gitlab.com/thefrol/notty/internal/entity/valid"
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
			name: "интервал наоборот",
			sub:  ReversedInterval(ValidSub()),
			has:  valid.ErrorInvalidPeriod,
		},
		{
			name: "корявый айдишник",
			sub:  FancyId(ValidSub()),
			has:  valid.ErrorIdValidation,
		},
		{
			name: "С матом",
			sub:  ExplicitText(),
			has:  valid.ErrorExplicitLanguage,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.sub.Validate()
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

// ReversedInterval разворачивает интервал в подписке, конец становится началм
// а начало концом, аминь
func ReversedInterval(s entity.Subscription) entity.Subscription {
	s.End, s.Start = s.Start, s.End
	return s
}

func FancyId(s entity.Subscription) entity.Subscription {
	s.Id = "my_LOVELY_id"
	return s
}
