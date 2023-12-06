// тут тестируется сгенерированный код и хендлеры
package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
	"gitlab.com/thefrol/notty/internal/api"
	"gitlab.com/thefrol/notty/internal/app"
	"gitlab.com/thefrol/notty/internal/entity"
	"gitlab.com/thefrol/notty/internal/mock"
)

// тут больно интегральные тесты, их бы в отдельную папочку
// todo

const (
	existingUserId = "api-test-user"
)

type ApiTestSuite struct {
	suite.Suite
	api               api.Server
	customersRepoMock *mock.CustomerereMock
	handlers          http.Handler
}

func (suite *ApiTestSuite) SetupTest() {
	//make mock cutomers service

	mc := minimock.NewController(suite.T())
	customersRepo := mock.NewCustomerereMock(mc)
	suite.customersRepoMock = customersRepo

	//existing user

	existingUser := entity.Customer{
		Id:   existingUserId,
		Name: "Апий Тестировальный",
	}
	//preparing customer repo
	customersRepo.GetMock.Set(func(ctx context.Context, s1 string) (c1 entity.Customer, err error) {
		if s1 == existingUser.Id {
			return entity.Customer{
				Id:   existingUserId,
				Name: "Апий Тестировальный",
			}, nil
		} else {
			return entity.Customer{}, app.ErrorCustomerNotFound
		}
	})

	customersRepo.UpdateMock.Set(func(ctx context.Context, c1 entity.Customer) (entity.Customer, error) {
		if c1.Id == existingUserId {
			return existingUser, nil
		} else {
			suite.Fail("Запрошен неизвестный айди %s для обновления", c1.Id)
			return entity.Customer{}, fmt.Errorf("не удалось обновить")
		}
	})

	customersRepo.DeleteMock.Set(func(ctx context.Context, id string) (err error) {
		if id == existingUserId {
			return nil
		} else {
			suite.Fail("Запрошен неизвестный айди %s для обновления", id)
			return fmt.Errorf("не удалось обновить")
		}
	})

	customersRepo.CreateMock.Set(func(ctx context.Context, c entity.Customer) (entity.Customer, error) {
		if c.Id == existingUser.Id {
			return entity.Customer{}, app.ErrorCustomerExists
		}
		return c, nil
	})

	// app
	app := app.New(suite.customersRepoMock, nil, nil, nil)
	//api
	suite.api = api.New(app, log.Logger)
	suite.handlers = suite.api.OpenAPI()
}

func (suite *ApiTestSuite) TestCustomerGetById() {
	tests := []struct {
		name             string
		id               string
		ReturnCode       int
		CustomerName     string
		validateResponse bool // ошибка при анмаршалинге
	}{
		{
			name:             "апий позитивный",
			id:               existingUserId,
			CustomerName:     "Апий Тестировальный",
			ReturnCode:       200,
			validateResponse: true,
		},
		{
			name:             "Если полльзователя не существует, то вернуть 404",
			id:               "api-test-user-not-exist",
			CustomerName:     "не существеющее имя",
			ReturnCode:       404,
			validateResponse: false,
		},
		// а ещё если корявый айдишник - вернуть ошибку валидации
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()

			req := httptest.NewRequest(
				http.MethodGet,
				"/customer/"+tt.id,
				nil)

			suite.handlers.ServeHTTP(rec, req)

			// вадидируем запрос
			suite.Equal(tt.ReturnCode, rec.Code)

			// валидируем ответ
			c := entity.Customer{}
			err := json.NewDecoder(rec.Body).Decode(&c)

			if !tt.validateResponse {
				return
			}
			suite.NoError(err)

			suite.Equal(c.Id, tt.id)
			suite.Equal(tt.CustomerName, c.Name)
		})
	}

}

func (suite *ApiTestSuite) TestDeleteCustomer() {
	tests := []struct {
		name       string
		id         string
		ReturnCode int
	}{
		{
			name:       "апий позитивный",
			id:         existingUserId,
			ReturnCode: 200,
		},
		{
			name:       "Если полльзователя не существует, то вернуть 404",
			id:         "api-test-user-not-exist",
			ReturnCode: 404,
		},
		// а ещё если корявый айдишник - вернуть ошибку валидации
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()

			req := httptest.NewRequest(
				http.MethodDelete,
				"/customer/"+tt.id,
				nil)

			suite.handlers.ServeHTTP(rec, req)

			// вадидируем запрос
			suite.Equal(tt.ReturnCode, rec.Code)
		})
	}

}

func (suite *ApiTestSuite) TestUpdateExistingCustomer() {
	rec := httptest.NewRecorder()
	id := existingUserId

	c := entity.Customer{
		Id:    id,
		Name:  "Новоименный Данил",
		Phone: "+72223334455",
	}

	data, err := json.Marshal(&c)
	suite.Require().NoError(err)

	req := httptest.NewRequest(
		http.MethodPost,
		"/customer/"+id,
		bytes.NewBuffer(data))

	suite.handlers.ServeHTTP(rec, req)

	bodyString := rec.Body.String()

	suite.Equalf(http.StatusOK, rec.Code, "неправильный код, тело сообщения: %s", bodyString)

	getCallsCount := len(suite.customersRepoMock.GetMock.Calls())
	updateCallsCount := len(suite.customersRepoMock.UpdateMock.Calls())
	suite.Equal(1, getCallsCount)
	suite.Equal(1, updateCallsCount)
}

func (suite *ApiTestSuite) TestUpdateNotExistingCustomer() {
	rec := httptest.NewRecorder()

	c := entity.Customer{
		Id:    "some-inexistant-url",
		Name:  "Новоименный Данил",
		Phone: "+72223334455",
	}

	data, err := json.Marshal(&c)
	suite.Require().NoError(err)

	req := httptest.NewRequest(
		http.MethodPost,
		"/customer/"+c.Id,
		bytes.NewBuffer(data))

	suite.handlers.ServeHTTP(rec, req)

	bodyString := rec.Body.String()

	suite.Equalf(http.StatusNotFound, rec.Code, "неправильный код, тело сообщения: %s", bodyString)

	getCallsCount := len(suite.customersRepoMock.GetMock.Calls())
	updateCallsCount := len(suite.customersRepoMock.UpdateMock.Calls())
	suite.Equal(1, getCallsCount)
	suite.Equal(0, updateCallsCount)
}

// тут мы задаем id, нужно чтобы он считал с базы
func (suite *ApiTestSuite) TestCreateNotExistingCustomer() {
	rec := httptest.NewRecorder()

	c := entity.Customer{
		Id:    "new-id",
		Name:  "Новоименный Данил",
		Phone: "+72223334455",
	}

	data, err := json.Marshal(&c)
	suite.Require().NoError(err)

	req := httptest.NewRequest(
		http.MethodPost,
		"/customer",
		bytes.NewBuffer(data))

	suite.handlers.ServeHTTP(rec, req)

	bodyString := rec.Body.String()

	suite.Equalf(http.StatusCreated, rec.Code, "неправильный код, тело сообщения: %s", bodyString)

	getCallsCount := len(suite.customersRepoMock.GetMock.Calls())
	suite.Equal(1, getCallsCount)
}

func (suite *ApiTestSuite) TestCreateNoId() {
	rec := httptest.NewRecorder()

	c := entity.Customer{
		Name:  "Новоименный Данил",
		Phone: "+72223334455",
	}

	data, err := json.Marshal(&c)
	suite.Require().NoError(err)

	req := httptest.NewRequest(
		http.MethodPost,
		"/customer",
		bytes.NewBuffer(data))

	suite.handlers.ServeHTTP(rec, req)

	bodyString := rec.Body.String()

	suite.Equalf(http.StatusCreated, rec.Code, "неправильный код, тело сообщения: %s", bodyString)

	getCallsCount := len(suite.customersRepoMock.GetMock.Calls())
	suite.Equal(0, getCallsCount, "количество запросов гет в репозиторий должно быть 0")
}

func (suite *ApiTestSuite) TestCreateCustomerExists() {
	rec := httptest.NewRecorder()

	c := entity.Customer{
		Id:    existingUserId,
		Name:  "Новоименный Данил",
		Phone: "+72223334455",
	}

	data, err := json.Marshal(&c)
	suite.Require().NoError(err)

	req := httptest.NewRequest(
		http.MethodPost,
		"/customer",
		bytes.NewBuffer(data))

	suite.handlers.ServeHTTP(rec, req)

	bodyString := rec.Body.String()

	suite.Equalf(http.StatusConflict, rec.Code, "неправильный код, тело сообщения: %s", bodyString)

}

func TestCustomerTestSuite(t *testing.T) {
	suite.Run(t, new(ApiTestSuite))
}
