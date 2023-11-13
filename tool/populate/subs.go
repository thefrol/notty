package main

import (
	"time"

	"gitlab.com/thefrol/notty/internal/entity"
)

var subs = []entity.Subscription{
	{
		Id:             "my-lovers",
		Text:           "Привет, зайчик!",
		Desc:           "Подписка, которая смотрит за тегом love. Действует 1444 часа с момента создания",
		Start:          time.Now(),
		End:            time.Now().Add(1444 * time.Hour),
		PhoneFilter:    "",
		OperatorFilter: "",
		TagFilter:      "^love$",
	},
}
