package entity

import (
	"gitlab.com/thefrol/notty/internal/entity/valid"
)

func (s Subscription) Validate() error {

	return valid.If(
		valid.Period(s.Start, s.End),
		valid.Id(s.Id),
		valid.HasNoExplicitLanguage(s.Text),
	)
}
