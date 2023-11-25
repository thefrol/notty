//go:build integral

package integralstorage

import "gitlab.com/thefrol/notty/internal/entity"

func (suite *Storage) TestSpawning() {
	expectedMessageCount := 2
	statusRequested := "TEST_STATUS"

	batch, err := suite.messages.Spawn(10, statusRequested)
	suite.NoError(err)

	// проверим, что мы сгенерировали два сообщения
	suite.Equal(expectedMessageCount, len(batch))

	// проверим, что они созданы по правильным рассылкам
	suite.True(In(batch, "anna", "will-send"), "Должен быть в отправленных")
	suite.True(In(batch, "ivan-testov", "will-send"), "Должен быть в отправленных")

	// найдем все сгенерированные сообщения в базе
	for _, m := range batch {
		m, err := suite.messages.Get(m.Id)
		suite.NoError(err)                     // найдено
		suite.Equal(statusRequested, m.Status) // с правильным статусом созданы
	}

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
