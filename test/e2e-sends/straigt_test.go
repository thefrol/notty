//go:build e2e
// +build e2e

package e2esends

import (
	"gitlab.com/thefrol/notty/internal/dto"
	"gitlab.com/thefrol/notty/internal/entity"
)

// TestSendingByTag проверяет что отправляются два сообщения, которые
// подходят по фильтру тегов moscow.%, при этом отправляются только сообщения
// с активной рассылки, рассылка из прошлого не срабатывает,
// как не срабатывает рассылка из будущего(которая пока не реализована)
// Так же не отправляется сообщение человеку из другого города с тегом
// rostov.best
func (suite *FromDbToSend) TestSendingByTag() {
	// results I want
	sendWant := 2

	// marking all as done
	sended := 0
	suite.senderMock.WorkMock.Set(func(ch1 <-chan entity.Message) (ch2 <-chan entity.Message, err error) {
		ch := make(chan entity.Message)

		go func() {
			for m := range ch1 {
				//тут какой-то дедлок происходит если сендер мок вот так использовать
				// видимо у них там общий лок на мок,
				// поэтому считаем количество отправленных тут!
				// а не где-то там
				// конечно логику надо как-то менять
				// и этот воркер - нетестируемая штука

				//err := suite.senderMock.Send(m)
				//suite.NoError(err)

				sended++ // типа отправлено )))
				m.SentNow()
				ch <- m
			}
			close(ch)
		}()

		return ch, nil
	})

	suite.app.FindAndSend(30, 1)

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
