package main

import (
	"log"
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
		PhoneFilter:    "%",
		OperatorFilter: "%",
		TagFilter:      "%love%",
	},
	{
		Id:             "my-old",
		Text:           "Превед, медвед!",
		Desc:           "Подписка, которая устарела когда-то давно. В 2003 году",
		Start:          MustParse("2006-01-02", "2003-03-08"),
		End:            MustParse("2006-01-02", "2003-06-08"),
		PhoneFilter:    "%",
		OperatorFilter: "%",
		TagFilter:      "%",
	},
}

func MustParse(format, s string) time.Time {
	t, err := time.Parse(format, s)
	if err != nil {
		log.Fatal(err)
	}
	return t
}
