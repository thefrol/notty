package server_test

import (
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/suite"
	"gitlab.com/thefrol/notty/internal/app"
	"gitlab.com/thefrol/notty/internal/app/server"
	"gitlab.com/thefrol/notty/internal/dto"
	"gitlab.com/thefrol/notty/internal/entity"
	"gitlab.com/thefrol/notty/internal/mock"
)

type AppTestSuite struct {
	suite.Suite
	app server.App
}

func (suite *AppTestSuite) SetupTest() {
	mc := minimock.NewController(suite.T())

	customers := mock.NewCustomerereMock(mc)
	subs := mock.NewSubscripterMock(mc)
	stats := mock.NewStatisterMock(mc)

	// статистика по всем сообщениям
	stats.AllMock.Expect().Return(dto.Statistics{
		"done": 4,
	}, nil)

	suite.app = server.New(customers, subs, stats)

	// наполняем моки

	// для общей статистики
	stats.FilterMock.When("%", "%", "%").Then(dto.Statistics{
		"done": 4,
	}, nil)

	// статистика по клиентам
	stats.FilterMock.When("%", "test-customer", "%").Then(dto.Statistics{
		"done": 2,
	}, nil)

	customers.GetMock.
		When("test-customer").
		Then(entity.Customer{}, nil)

	customers.GetMock.
		When("no-customer").
		Then(entity.Customer{}, app.ErrorCustomerNotFound)

	// статистика по клиентам
	stats.FilterMock.When("test-sub", "%", "%").Then(dto.Statistics{
		"done": 5,
	}, nil)

	subs.GetMock.
		When("test-sub").
		Then(entity.Subscription{}, nil)

	subs.GetMock.
		When("no-sub").
		Then(entity.Subscription{}, app.ErrorSubscriptionNotFound)

}

func TestAppSuite(t *testing.T) {
	suite.Run(t, new(AppTestSuite))
}

// Тестируем статистику

func (suite *AppTestSuite) TestAllStatistics() {
	st, err := suite.app.FullStats()
	suite.NoError(err)
	suite.Equal(dto.Statistics{"done": 4}, st)
}

func (suite *AppTestSuite) TestCustomerStatistics() {
	st, err := suite.app.StatsByClient("test-customer")
	suite.NoError(err)
	suite.Equal(dto.Statistics{"done": 2}, st)
}

func (suite *AppTestSuite) TestCustomerNotExistsStatistics() {
	_, err := suite.app.StatsByClient("no-customer")
	suite.ErrorIs(err, app.ErrorCustomerNotFound)
}

func (suite *AppTestSuite) TestSubStatistics() {
	st, err := suite.app.StatsBySubscription("test-sub")
	suite.NoError(err)
	suite.Equal(dto.Statistics{"done": 5}, st)
}

func (suite *AppTestSuite) TestSubNotExistsStatistics() {
	_, err := suite.app.StatsBySubscription("no-sub")
	suite.ErrorIs(err, app.ErrorSubscriptionNotFound)
}
