package customers

import (
	"testing"

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

	clients := New(conn)

	c := entity.Customer{
		Id:       "test_one",
		Name:     "Dina",
		Operator: "916",
		Phone:    "+7916",
		Tag:      "test_user",
	}
	_, err := clients.Create(c)
	require.NoError(t, err)

	g, err := clients.Get(c.Id)
	require.NoError(t, err)

	assert.Equal(t, c, g)

	err = clients.Delete(c.Id)
	require.NoError(t, err)

	_, err = clients.Get(c.Id)
	require.Error(t, err)

}
