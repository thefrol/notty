//go:build e2e
// +build e2e

// Это интегральный тест, из базы до отправки. Отправщик мокируется.
// использует переменную окружения NOTTY_TEST_DB для тестирования
package e2esends

import (
	"os"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
	"gitlab.com/thefrol/notty/internal/mock"
	"gitlab.com/thefrol/notty/internal/notifyloop"
	"gitlab.com/thefrol/notty/internal/storage/postgres"
	"gitlab.com/thefrol/notty/internal/storage/sqlrepo"
)

type FromDbToSend struct {
	suite.Suite
	app        *notifyloop.Notifyer
	senderMock *mock.SenderMock
	messages   sqlrepo.Messages
	stats      sqlrepo.Statistics
}

func (suite *FromDbToSend) SetupTest() {
	str := os.Getenv("NOTTY_TEST_DB")
	db := postgres.MustConnect(str)

	err := Purge(db)
	if err != nil {
		suite.FailNowf("ошибка инициализации", "Не могу очистить тестовую базу: %v", err)
	}

	err = Fill(db)
	if err != nil {
		suite.FailNowf("ошибка инициализации", "Не могу наполнить тестовую базу: %v", err)
	}

	// создаем мок отправщика
	mc := minimock.NewController(suite.T())
	suite.senderMock = mock.NewSenderMock(mc)

	mr := sqlrepo.NewMessages(db, log.Logger)

	suite.app = notifyloop.NewNotifyer(mr, suite.senderMock)

	suite.messages = mr

	// создаем сервис статистики
	suite.stats = *sqlrepo.NewStatistics(db)
}

func (suite *FromDbToSend) TestNoop() {
	suite.Equal(1, 1)
}

func TestE2ESending(t *testing.T) {
	suite.Run(t, new(FromDbToSend))
}
