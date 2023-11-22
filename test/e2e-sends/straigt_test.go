//go:build e2e
// +build e2e

package e2esends

import "gitlab.com/thefrol/notty/internal/entity"

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
}
