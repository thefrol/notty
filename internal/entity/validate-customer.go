package entity

import (
	"gitlab.com/thefrol/notty/internal/entity/valid"
)

func (c Customer) Validate() error {

	return valid.If(
		valid.Phone(c.Phone),
		valid.Id(c.Id),
	)
}
