//go:build integral

package integralstorage

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	"gitlab.com/thefrol/notty/internal/storage/postgres"
	"gitlab.com/thefrol/notty/internal/storage/sqlrepo"
)

type Storage struct {
	suite.Suite
	messages sqlrepo.Messages
}

func (suite *Storage) SetupTest() {
	str := os.Getenv("NOTTY_TEST_DB")
	db := postgres.MustConnect(str)

	err := Purge(db)
	if err != nil {
		suite.FailNowf("ошибка инициализации", "Не могу очистить тестовую базу: %v", err)
	}

	err = Fill(db) // todo выделить как-то отдельно лол
	if err != nil {
		suite.FailNowf("ошибка инициализации", "Не могу наполнить тестовую базу: %v", err)
	}

	// создаем мок отправщика
	suite.messages = sqlrepo.NewMessages(db)
}

func TestIntergralSQLRepositories(t *testing.T) {
	suite.Run(t, new(Storage))
}
