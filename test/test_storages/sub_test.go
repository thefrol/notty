package storages_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/thefrol/notty/internal/entity"
	"gitlab.com/thefrol/notty/internal/storage/postgres"
	"gitlab.com/thefrol/notty/internal/storage/sqlrepo"
)

func Test_SubscriptionWorkflow(t *testing.T) {
	if !NoSkip && TestDSN == "" {
		t.Skip("Пропускаю тесты на рассылки")
	}

	conn := postgres.MustConnect(TestDSN)

	subscriptions := sqlrepo.NewSubscriptions(conn)

	date := time.Date(2011, 2, 2, 13, 21, 3, 12, time.Local)

	c := entity.Subscription{
		Id:             "test-one",
		Start:          date,
		End:            date.Add(time.Hour * 10),
		Text:           "text",
		Desc:           "desc",
		OperatorFilter: "1",
		PhoneFilter:    "123",
		TagFilter:      "1252",
	}
	_, err := subscriptions.Create(context.TODO(), c)
	require.NoError(t, err)

	g, err := subscriptions.Get(context.TODO(), c.Id)

	require.NoError(t, err)

	g.Start = c.Start // todo какаято беда с часовыми поясами
	g.End = c.End
	assert.Equal(t, c, g)

	err = subscriptions.Delete(context.TODO(), c.Id)
	require.NoError(t, err)

	_, err = subscriptions.Get(context.TODO(), c.Id)
	require.Error(t, err)

}
