// тут тестируется сгенерированный код и хендлеры
package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/suite"
	"gitlab.com/thefrol/notty/internal/api"
	"gitlab.com/thefrol/notty/internal/app"
	"gitlab.com/thefrol/notty/internal/app/server"
	"gitlab.com/thefrol/notty/internal/entity"
	"gitlab.com/thefrol/notty/internal/service"
	"gitlab.com/thefrol/notty/internal/service/mock"
)

// тут больно интегральные тесты, их бы в отдельную папочку
// todo

const (
	existingSubscriptionId = "api-test-sub"
)

type SubscriptionTestSuite struct {
	suite.Suite
	api                   api.Server
	subscriptionsRepoMock *mock.SubscriptionRepositoryMock
	handlers              http.Handler
}

func (suite *SubscriptionTestSuite) SetupTest() {
	//make mock cutomers service

	mc := minimock.NewController(suite.T())
	subscriptionsRepo := mock.NewSubscriptionRepositoryMock(mc)
	suite.subscriptionsRepoMock = subscriptionsRepo
	//preparing subscription repo
	subscriptionsRepo.GetMock.Set(func(s1 string) (c1 entity.Subscription, err error) {
		if s1 == existingSubscriptionId {
			return entity.Subscription{
				Id:   existingSubscriptionId,
				Desc: "рассылочка моя",
			}, nil
		} else {
			return entity.Subscription{}, app.ErrorSubscriptionNotFound
		}
	})

	subscriptionsRepo.UpdateMock.Set(func(c1 entity.Subscription) (err error) {
		if c1.Id == existingSubscriptionId {
			return nil
		} else {
			suite.Fail("Запрошен неизвестный айди %s для обновления", c1.Id)
			return fmt.Errorf("не удалось обновить")
		}
	})

	subscriptionsRepo.DeleteMock.Set(func(id string) (err error) {
		if id == existingSubscriptionId {
			return nil
		} else {
			suite.Fail("Запрошен неизвестный айди %s для обновления", id)
			return fmt.Errorf("не удалось обновить")
		}
	})

	subscriptionsRepo.CreateMock.Set(func(c entity.Subscription) (err error) {
		if c.Id == existingSubscriptionId {
			return fmt.Errorf("такой пользователь существует")
		}
		return nil
	})

	// services
	subscriptionService := service.NewSubscriptions(subscriptionsRepo)

	// app
	app := server.New(nil, subscriptionService, nil)
	//api
	suite.api = api.New(app)
	suite.handlers = suite.api.OpenAPI()
}

func (suite *SubscriptionTestSuite) TestSubscriptionGetById() {
	tests := []struct {
		name             string
		id               string
		ReturnCode       int
		SubscriptionDesc string
		validateResponse bool // ошибка при анмаршалинге
	}{
		{
			name:             "апий позитивный",
			id:               existingSubscriptionId,
			SubscriptionDesc: "рассылочка моя",
			ReturnCode:       200,
			validateResponse: true,
		},
		{
			name:             "Если полльзователя не существует, то вернуть 404",
			id:               "api-test-subscription-not-exist",
			SubscriptionDesc: "не существеющее имя",
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
				"/sub/"+tt.id,
				nil)

			suite.handlers.ServeHTTP(rec, req)

			// вадидируем запрос
			suite.Equal(tt.ReturnCode, rec.Code)

			// валидируем ответ
			c := entity.Subscription{}
			err := json.NewDecoder(rec.Body).Decode(&c)

			if !tt.validateResponse {
				return
			}
			suite.NoError(err)

			suite.Equal(c.Id, tt.id)
			suite.Equal(tt.SubscriptionDesc, c.Desc)
		})
	}

}

func (suite *SubscriptionTestSuite) TestDeleteSubscription() {
	tests := []struct {
		name       string
		id         string
		ReturnCode int
	}{
		{
			name:       "существующая рассылка",
			id:         existingSubscriptionId,
			ReturnCode: 200,
		},
		{
			name:       "Если рассылки не существует, то вернуть 404",
			id:         "api-test-subscription-not-exist",
			ReturnCode: 404,
		},
		// а ещё если корявый айдишник - вернуть ошибку валидации
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()

			req := httptest.NewRequest(
				http.MethodDelete,
				"/sub/"+tt.id,
				nil)

			suite.handlers.ServeHTTP(rec, req)

			// вадидируем запрос
			suite.Equal(tt.ReturnCode, rec.Code)
		})
	}

}

func (suite *SubscriptionTestSuite) TestUpdateExistingSubscription() {
	rec := httptest.NewRecorder()
	id := existingSubscriptionId

	c := entity.Subscription{
		Id:   id,
		Desc: "новое описание",
	}

	data, err := json.Marshal(&c)
	if err != nil {
		log.Fatal(err)
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"/sub/"+id,
		bytes.NewBuffer(data))

	suite.handlers.ServeHTTP(rec, req)

	bodyString := rec.Body.String()

	suite.Equalf(http.StatusOK, rec.Code, "неправильный код, тело сообщения: %s", bodyString)

	getCallsCount := len(suite.subscriptionsRepoMock.GetMock.Calls())
	updateCallsCount := len(suite.subscriptionsRepoMock.UpdateMock.Calls())
	suite.Equal(2, getCallsCount)
	suite.Equal(1, updateCallsCount)
}

func (suite *SubscriptionTestSuite) TestUpdateNotExistingSubscription() {
	rec := httptest.NewRecorder()

	c := entity.Subscription{
		Id: "some-inexistant-url",
	}

	data, err := json.Marshal(&c)
	if err != nil {
		log.Fatal(err)
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"/sub/"+c.Id,
		bytes.NewBuffer(data))

	suite.handlers.ServeHTTP(rec, req)

	bodyString := rec.Body.String()

	suite.Equalf(http.StatusNotFound, rec.Code, "неправильный код, тело сообщения: %s", bodyString)

	getCallsCount := len(suite.subscriptionsRepoMock.GetMock.Calls())
	updateCallsCount := len(suite.subscriptionsRepoMock.UpdateMock.Calls())
	suite.Equal(1, getCallsCount)
	suite.Equal(0, updateCallsCount)
}

// тут мы задаем id, нужно чтобы он считал с базы
func (suite *SubscriptionTestSuite) TestCreateNotExistingSubscription() {
	rec := httptest.NewRecorder()

	c := entity.Subscription{
		Id: "new-id",
	}

	data, err := json.Marshal(&c)
	if err != nil {
		log.Fatal(err)
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"/sub",
		bytes.NewBuffer(data))

	suite.handlers.ServeHTTP(rec, req)

	bodyString := rec.Body.String()

	suite.Equalf(http.StatusCreated, rec.Code, "неправильный код, тело сообщения: %s", bodyString)

	getCallsCount := len(suite.subscriptionsRepoMock.GetMock.Calls())
	suite.Equal(1, getCallsCount)
}

func (suite *SubscriptionTestSuite) TestCreateNoId() {
	rec := httptest.NewRecorder()

	c := entity.Subscription{}

	data, err := json.Marshal(&c)
	if err != nil {
		log.Fatal(err)
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"/sub",
		bytes.NewBuffer(data))

	suite.handlers.ServeHTTP(rec, req)

	bodyString := rec.Body.String()

	suite.Equalf(http.StatusCreated, rec.Code, "неправильный код, тело сообщения: %s", bodyString)

	getCallsCount := len(suite.subscriptionsRepoMock.GetMock.Calls())
	suite.Equal(0, getCallsCount)
}

func (suite *SubscriptionTestSuite) TestCreateSubscriptionExists() {
	rec := httptest.NewRecorder()

	c := entity.Subscription{
		Id: existingSubscriptionId,
	}

	data, err := json.Marshal(&c)
	if err != nil {
		log.Fatal(err)
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"/sub",
		bytes.NewBuffer(data))

	suite.handlers.ServeHTTP(rec, req)

	bodyString := rec.Body.String()

	suite.Equalf(http.StatusConflict, rec.Code, "неправильный код, тело сообщения: %s", bodyString)

}

func TestSubscriptionsSuite(t *testing.T) {
	suite.Run(t, new(SubscriptionTestSuite))
}
