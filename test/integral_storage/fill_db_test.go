//go:build integral

package integralstorage

import (
	"database/sql"
	"time"

	"github.com/pressly/goose/v3"
	"gitlab.com/thefrol/notty/internal/entity"
	_ "gitlab.com/thefrol/notty/internal/storage/postgres/migrations"
	"gitlab.com/thefrol/notty/internal/storage/sqlrepo"
)

func Purge(db *sql.DB) error {
	err := goose.Reset(db, ".")
	if err != nil {
		return err
	}

	return goose.Up(db, ".")
}

func Fill(db *sql.DB) error {
	cs := sqlrepo.NewCustomers(db)
	ss := sqlrepo.NewSubscriptions(db)

	err := FillCustomers(cs)
	if err != nil {
		return err
	}

	err = FillSubs(ss)
	if err != nil {
		return err
	}

	return nil
}

func FillCustomers(cs sqlrepo.Customers) error {
	var err error
	err = cs.Create(entity.Customer{
		Id:       "anna",
		Name:     "Анна Новикова",
		Phone:    "+79163332211",
		Operator: "mts",
		Tag:      "moscow.best",
	})
	if err != nil {
		return err
	}

	err = cs.Create(entity.Customer{
		Id:       "ksu",
		Name:     "Ксения Иванова",
		Phone:    "+79163332211",
		Operator: "mts",
		Tag:      "rostov.best",
	})
	if err != nil {
		return err
	}

	err = cs.Create(entity.Customer{
		Id:       "ivan-testov",
		Name:     "Иван Тестов",
		Phone:    "+79163332211",
		Operator: "mts",
		Tag:      "moscow.test",
	})
	if err != nil {
		return err
	}

	return nil
}

func FillSubs(ss sqlrepo.Subscriptions) error {
	var err error
	err = ss.Create(entity.Subscription{
		Id:             "will-send",
		Text:           "Текст сообщения",
		Desc:           "рассылка, которая разошлется",
		PhoneFilter:    "%",
		OperatorFilter: "%",
		TagFilter:      "moscow.%",
		Start:          time.Now(),
		End:            time.Now().Add(time.Hour * 144),
	})
	if err != nil {
		return err
	}

	err = ss.Create(entity.Subscription{
		Id:             "no-send",
		Text:           "Что-то из далекого прошло случано попало в будущее",
		Desc:           "не должно быть отправлено",
		PhoneFilter:    "%",
		OperatorFilter: "%",
		TagFilter:      "moscow.%",
		Start:          time.Date(2006, 2, 11, 11, 22, 33, 12, time.Local),
		End:            time.Date(2007, 2, 11, 11, 22, 33, 12, time.Local),
	})
	if err != nil {
		return err
	}

	return nil
}
