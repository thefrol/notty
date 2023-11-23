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
func Test_UpdateSubscription(t *testing.T) {
	c := minimock.NewController(t)
	mc := mock.NewSubscriptionRepositoryMock(c)

	id := "test-id"

	mc.GetMock.
		Expect(id).
		Return(
			entity.Subscription{Id: id, Desc: "простая рассылка"},
			nil,
		)

	//always ok
	mc.UpdateMock.
		Expect(entity.Subscription{Id: id, Desc: "непростая рассылка"}).
		Return(nil)

	// actual business
	svc := service.NewSubscriptions(mc)
	_, err := svc.Update(entity.Subscription{
		Id:   "test-id",
		Desc: "непростая рассылка",
	})
	assert.NoError(t, err)

	const GetCallsCount = 2
	assert.Equal(t, GetCallsCount, len(mc.GetMock.Calls()))
}

func Test_UpdateSubscriptionNotFound(t *testing.T) {
	c := minimock.NewController(t)
	mc := mock.NewSubscriptionRepositoryMock(c)

	id := "test-id"

	mc.GetMock.
		Expect(id).
		Return(
			entity.Subscription{},
			app.ErrorSubscriptionNotFound,
		)

	// actual business
	svc := service.NewSubscriptions(mc)
	_, err := svc.Update(entity.Subscription{
		Id:   "test-id",
		Desc: "непростая рассылка",
	})
	assert.ErrorIs(t, err, app.ErrorSubscriptionNotFound)

	const GetCallsCount = 1
	assert.Equal(t, GetCallsCount, len(mc.GetMock.Calls()))
}

func Test_DeleteSubscription(t *testing.T) {
	c := minimock.NewController(t)
	mc := mock.NewSubscriptionRepositoryMock(c)

	id := "test-id"

	mc.GetMock.Expect(id).Return(entity.Subscription{Id: id}, nil)

	mc.DeleteMock.Expect(id).Return(nil)

	// actual business
	svc := service.NewSubscriptions(mc)
	svc.Delete("test-id")

	assert.Equal(t, calledOnce, len(mc.DeleteMock.Calls()))

}
