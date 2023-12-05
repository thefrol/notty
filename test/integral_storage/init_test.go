//go:build integral

package integralstorage

import (
	"database/sql"
	"os"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
	"gitlab.com/thefrol/notty/internal/storage/postgres"
	"gitlab.com/thefrol/notty/internal/storage/sqlrepo"
)

type Storage struct {
	suite.Suite
	messages sqlrepo.Messages
	db       *sql.DB
}

func (suite *Storage) SetupTest() {
	str := os.Getenv("NOTTY_TEST_DB")
	db := postgres.MustConnect(str)
	suite.db = db

	err := Purge(db)
	if err != nil {
		suite.FailNowf("ошибка инициализации", "Не могу очистить тестовую базу: %v", err)
	}

	err = Fill(db) // todo выделить как-то отдельно лол
	if err != nil {
		suite.FailNowf("ошибка инициализации", "Не могу наполнить тестовую базу: %v", err)
	}

	// создаем мок отправщика
	suite.messages = sqlrepo.NewMessages(db, log.Logger)
}

func TestIntergralSQLRepositories(t *testing.T) {
	suite.Run(t, new(Storage))
}
