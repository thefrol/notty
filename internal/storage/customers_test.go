package service_test

import (
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"gitlab.com/thefrol/notty/internal/app"
	"gitlab.com/thefrol/notty/internal/entity"
	service "gitlab.com/thefrol/notty/internal/storage"
	"gitlab.com/thefrol/notty/internal/storage/mock"
)

// тут мы проверяем, что нет ошибки, и что вызвано два раза Get и один раз Апдейт
// Какие данные он там менял, мы не проверяем
func Test_UpdateCustomer(t *testing.T) {
	c := minimock.NewController(t)
	mc := mock.NewCustomerRepositoryMock(c)

	id := "test-id"

	mc.GetMock.
		Expect(id).
		Return(
			entity.Customer{Id: id, Name: "Максим Успешнов"},
			nil,
		)

	//always ok
	mc.UpdateMock.
		Expect(entity.Customer{Id: id, Name: "Максимка Успешнов"}).
		Return(nil)

	// actual business
	svc := service.NewCustomers(mc)
	_, err := svc.Update(entity.Customer{
		Id:   "test-id",
		Name: "Максимка Успешнов",
	})
	assert.NoError(t, err)

	const GetCallsCount = 2
	assert.Equal(t, GetCallsCount, len(mc.GetMock.Calls()))
}

func Test_UpdateCustomerNotFound(t *testing.T) {
	c := minimock.NewController(t)
	mc := mock.NewCustomerRepositoryMock(c)

	id := "test-id"

	mc.GetMock.
		Expect(id).
		Return(
			entity.Customer{},
			app.ErrorCustomerNotFound,
		)

	// actual business
	svc := service.NewCustomers(mc)
	_, err := svc.Update(entity.Customer{
		Id:   "test-id",
		Name: "Максимка Успешнов",
	})
	assert.ErrorIs(t, err, app.ErrorCustomerNotFound)

	const GetCallsCount = 1
	assert.Equal(t, GetCallsCount, len(mc.GetMock.Calls()))
}

func Test_DeleteCustomer(t *testing.T) {
	c := minimock.NewController(t)
	mc := mock.NewCustomerRepositoryMock(c)

	id := "test-id"

	mc.GetMock.Expect(id).Return(entity.Customer{Id: id}, nil)

	mc.DeleteMock.Expect(id).Return(nil)

	// actual business
	svc := service.NewCustomers(mc)
	svc.Delete("test-id")

	assert.Equal(t, calledOnce, len(mc.DeleteMock.Calls()))

}

// todo testCreate

const calledOnce = 1
