//go:build integral

package integralstorage

import (
	"context"

	"gitlab.com/thefrol/notty/internal/entity"
)

func (suite *Storage) TestSpawning() {
	expectedMessageCount := 2
	statusRequested := "TEST_STATUS"

	batch, err := suite.messages.LockedSpawn(context.TODO(), 10, statusRequested)
	suite.NoError(err)

	// проверим, что мы сгенерировали два сообщения
	suite.Equal(expectedMessageCount, len(batch))

	// проверим, что они созданы по правильным рассылкам
	suite.True(In(batch, "anna", "will-send"), "Должен быть в отправленных")
	suite.True(In(batch, "ivan-testov", "will-send"), "Должен быть в отправленных")

	// найдем все сгенерированные сообщения в базе
	for _, m := range batch {
		m, err := suite.messages.Get(context.TODO(), m.Id)
		suite.NoError(err)                     // найдено
		suite.Equal(statusRequested, m.Status) // с правильным статусом созданы
	}

}

// Тут мы проверяем что sqlRepo.Messages нормально резервирует все сообщения
// которые мы просим в нужном количестве
func (suite *Storage) TestReserveFailedMessages() {
	suite.Run("взять правильные сообщения", func() {
		suite.InsertMessageWithStatus("id-1", entity.StatusInitial)
		suite.InsertMessageWithStatus("id-2", entity.StatusSent)
		suite.InsertMessageWithStatus("id-3", entity.StatusFailed)
		suite.InsertMessageWithStatus("id-4", entity.StatusFailed)
		// итого мы можем зарезервировать два сообщения со статусом фейлд.

		// todo надо добавить в constraints NOT NULL иначе мои Scan() падают
		// ну и указать в презентации как это вылезло из тестирования,
		// где я забивал не полную инфу

		got, err := suite.messages.ReserveFromStatus(context.TODO(), 5, entity.StatusFailed)
		suite.NoError(err)
		suite.Equal(2, len(got)) // мы получили два сообщения

		got, err = suite.messages.ReserveFromStatus(context.TODO(), 5, entity.StatusFailed)
		suite.NoError(err)
		suite.Equal(0, len(got)) // теперь мы получим ничего при втором запуске

	})

	suite.Run("достать указанное количество сообщений", func() {

		suite.InsertMessageWithStatus("id-11", entity.StatusFailed)
		suite.InsertMessageWithStatus("id-12", entity.StatusFailed)
		suite.InsertMessageWithStatus("id-13", entity.StatusFailed)
		suite.InsertMessageWithStatus("id-14", entity.StatusFailed)
		suite.InsertMessageWithStatus("id-15", entity.StatusFailed)
		suite.InsertMessageWithStatus("id-16", entity.StatusFailed)

		// запросим 3 хотя в наличе шесть
		requestN := 3
		got, err := suite.messages.ReserveFromStatus(context.TODO(), requestN, entity.StatusFailed)
		suite.NoError(err)
		suite.Equal(3, len(got)) // 3/6

		// теперь запросим пять, хотя осталось только три
		requestN = 5
		got, err = suite.messages.ReserveFromStatus(context.TODO(), requestN, entity.StatusFailed)
		suite.NoError(err)
		suite.Equal(3, len(got)) // 6/6
	})

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

func (s *Storage) InsertMessageWithStatus(id, status string) {
	rs, err := s.db.Exec(`
		INSERT INTO
			Messages(
				id,
				sub_id,
				customer_id,
				message_text,
				phone,
				status,
				sent)
		VALUES
			($1,'sub','cust','msg','+73332221100',$2,NULL)`, id, status)
	s.NoError(err)
	if n, err := rs.RowsAffected(); err != nil && n > int64(0) {
		s.FailNow("cant insert messages %v", err.Error())
	}
}
