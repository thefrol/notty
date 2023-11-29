//go:build e2e
// +build e2e

package e2esends

import (
	"context"
	"time"

	"gitlab.com/thefrol/notty/internal/dto"
	"gitlab.com/thefrol/notty/internal/entity"
	"gitlab.com/thefrol/notty/internal/service"
)

// TestSendingByTag проверяет что отправляются два сообщения, которые
// подходят по фильтру тегов moscow.%, при этом отправляются только сообщения
// с активной рассылки, рассылка из прошлого не срабатывает,
// как не срабатывает рассылка из будущего(которая пока не реализована)
// Так же не отправляется сообщение человеку из другого города с тегом
// rostov.best
func (suite *FromDbToSend) TestSendingByTag() {
	// вот сколько сообщений я хочу отправить
	sendWant := 2

	// тут мок будет просто считать количество запросов на отправку
	sended := 0
	suite.senderMock.SendMock.Set(func(m1 entity.Message) (err error) {
		// просто считаем количество
		sended++
		return nil
	})

	// сконфигурируем сервис отправки
	worker := service.Worker{
		Notifyer:  suite.app,
		Timeout:   0,
		BatchSize: 50,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // сервер остановится при выходе из теста

	// запустим его в горуине
	go func() {
		worker.FetchAndSend(ctx) // todo надо как-то этого парня вырубать)
	}()

	// дадим ему время
	time.Sleep(100 * time.Millisecond)
	cancel() // можно сервис уже тут останавливать

	// проверим

	suite.Equal(sendWant, sended)

	// а ещё нам нужно проверить, что они записаны в базу нормально
	// хотя это не по этому тесту конечно вообще

	msgs, err := suite.messages.ByStatus("done", 10)
	suite.NoError(err)

	suite.True(In(msgs, "anna", "will-send"), "Должен быть в отправленных")
	suite.True(In(msgs, "ivan-testov", "will-send"), "Должен быть в отправленных")

	// и заодно тогда статистику проверим
	anySub := "%"
	anyClient := "%"
	anyStatus := "%"

	// по всем
	s, err := suite.stats.All()
	suite.NoError(err)
	suite.Equal(dto.Statistics{"done": 2}, s)

	// по рассылкам
	s, err = suite.stats.Filter("will-send", anyClient, anyStatus)
	suite.NoError(err)
	suite.Equal(dto.Statistics{"done": 2}, s)

	s, err = suite.stats.Filter("no-send", anyClient, anyStatus)
	suite.NoError(err)
	suite.Equal(dto.Statistics{}, s)

	// по клиентам

	client := "anna"
	sub := anySub
	status := anyStatus

	s, err = suite.stats.Filter(sub, client, status)
	suite.NoError(err)
	suite.Equal(dto.Statistics{"done": 1}, s)

	client = "ivan-testov"
	sub = anySub
	status = anyStatus

	s, err = suite.stats.Filter(anySub, client, anyStatus)
	suite.NoError(err)
	suite.Equal(dto.Statistics{"done": 1}, s)

}

// In проверяет что сообщение с таким отправителем и
// рассылкой существует в батче сообщений
func In(msgs []entity.Message, custId, subId string) bool {
	for _, v := range msgs {
		if v.CustomerId == custId && v.SubscriptionId == subId {
			return true
		}
	}
	return false
}
