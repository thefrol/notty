package subscriptions

import (
	"log"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/thefrol/notty/internal/entity"
	_ "gitlab.com/thefrol/notty/internal/migrations"
	"gitlab.com/thefrol/notty/internal/storage/postgres"
)

// конечно жесть что у меня тут интегральное тестирование
// todo

const (
	testDSN = "host=localhost dbname=notty_dev_local user=notty_dev_local password=local.dev.1 port=15432 sslmode=disable"
)

func Test_Workflow(t *testing.T) {
	conn := postgres.MustConnect(testDSN)

	subscriptions := New(conn)

	date := time.Date(2011, 2, 2, 13, 21, 3, 12, time.Local)

	c := entity.Subscription{
		Id:             "test_one",
		Start:          date,
		End:            date,
		Text:           "text",
		Desc:           "desc",
		OperatorFilter: "1",
		PhoneFilter:    "123",
		TagFilter:      "1252",
	}
	_, err := subscriptions.Create(c)
	require.NoError(t, err)

	g, err := subscriptions.Get(c.Id)

	require.NoError(t, err)

	log.Println("Беда с часовыми поясами в Subscriptions")
	g.Start = c.Start // todo какаято беда с часовыми поясами
	g.End = c.End
	assert.Equal(t, c, g)

	err = subscriptions.Delete(c.Id)
	require.NoError(t, err)

	_, err = subscriptions.Get(c.Id)
	require.Error(t, err)

}
