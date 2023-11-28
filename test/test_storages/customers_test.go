package storages_test

import (
	"testing"

	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/thefrol/notty/internal/entity"
	"gitlab.com/thefrol/notty/internal/storage/postgres"
	"gitlab.com/thefrol/notty/internal/storage/sqlrepo"
)

// тут не надо проверять валидные и невалидные поля
// Это уже сделали юнит тесты

func Test_CustomerWorkflow(t *testing.T) {
	if !NoSkip && TestDSN == "" {
		t.Skip("Пропускаю тесты на клиентов")
	}

	conn := postgres.MustConnect(TestDSN)

	// очистить бд
	goose.Reset(conn, ".")
	goose.Up(conn, ".")

	// тестим
	clients := sqlrepo.NewCustomers(conn)

	c := entity.Customer{
		Id:       "test-one",
		Name:     "Dina",
		Operator: "916",
		Phone:    "+79161234533",
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
