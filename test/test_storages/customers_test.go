package storages_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/thefrol/notty/internal/entity"
	"gitlab.com/thefrol/notty/internal/storage/customers"
	"gitlab.com/thefrol/notty/internal/storage/postgres"
)

func Test_CustomerWorkflow(t *testing.T) {
	if !NoSkip && TestDSN == "" {
		t.Skip("Пропускаю тесты на клиентов")
	}
	conn := postgres.MustConnect(TestDSN)

	clients := customers.New(conn)

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
