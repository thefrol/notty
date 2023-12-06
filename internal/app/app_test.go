package app_test

import (
	"context"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/suite"
	"gitlab.com/thefrol/notty/internal/app"

	"gitlab.com/thefrol/notty/internal/entity"
	"gitlab.com/thefrol/notty/internal/mock"
)

type AppTestSuite struct {
	suite.Suite
	app app.App
}

func (suite *AppTestSuite) SetupTest() {
	mc := minimock.NewController(suite.T())

	customers := mock.NewCustomerereMock(mc)
	subs := mock.NewSubscripterMock(mc)
	stats := mock.NewStatisterMock(mc)

	// статистика по всем сообщениям
	stats.AllMock.Expect(context.Background()).Return(app.Statistics{
		"done": 4,
	}, nil)

	suite.app = app.New(customers, subs, stats, nil)

	// наполняем моки

	// для общей статистики
	stats.FilterMock.When(context.Background(), "%", "%", "%").Then(app.Statistics{
		"done": 4,
	}, nil)

	// статистика по клиентам
	stats.FilterMock.When(context.Background(), "%", "test-customer", "%").Then(app.Statistics{
		"done": 2,
	}, nil)

	customers.GetMock.
		When(context.Background(), "test-customer").
		Then(entity.Customer{}, nil)

	customers.GetMock.
		When(context.Background(), "no-customer").
		Then(entity.Customer{}, app.ErrorCustomerNotFound)

	// статистика по клиентам
	stats.FilterMock.When(context.Background(), "test-sub", "%", "%").Then(app.Statistics{
		"done": 5,
	}, nil)

	subs.GetMock.
		When(context.Background(), "test-sub").
		Then(entity.Subscription{}, nil)

	subs.GetMock.
		When(context.Background(), "no-sub").
		Then(entity.Subscription{}, app.ErrorSubscriptionNotFound)

}

func TestAppSuite(t *testing.T) {
	suite.Run(t, new(AppTestSuite))
}

// Тестируем статистику

func (suite *AppTestSuite) TestAllStatistics() {
	st, err := suite.app.FullStats(context.Background())
	suite.NoError(err)
	suite.Equal(app.Statistics{"done": 4}, st)
}

func (suite *AppTestSuite) TestCustomerStatistics() {
	st, err := suite.app.StatsByClient(context.Background(), "test-customer")
	suite.NoError(err)
	suite.Equal(app.Statistics{"done": 2}, st)
}

func (suite *AppTestSuite) TestCustomerNotExistsStatistics() {
	_, err := suite.app.StatsByClient(context.Background(), "no-customer")
	suite.ErrorIs(err, app.ErrorCustomerNotFound)
}

func (suite *AppTestSuite) TestSubStatistics() {
	st, err := suite.app.StatsBySubscription(context.Background(), "test-sub")
	suite.NoError(err)
	suite.Equal(app.Statistics{"done": 5}, st)
}

func (suite *AppTestSuite) TestSubNotExistsStatistics() {
	_, err := suite.app.StatsBySubscription(context.Background(), "no-sub")
	suite.ErrorIs(err, app.ErrorSubscriptionNotFound)
}
