package storages_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/thefrol/notty/internal/app"
	"gitlab.com/thefrol/notty/internal/storage/postgres"
	"gitlab.com/thefrol/notty/internal/storage/sqlrepo"
)

// тут не надо проверять валидные и невалидные поля
// Это уже сделали юнит тесты

func Test_StatsValues(t *testing.T) {
	if !NoSkip && TestDSN == "" {
		t.Skip("Пропускаю тесты на статистику")
	}

	conn := postgres.MustConnect(TestDSN)

	_, err := conn.Exec(`DELETE FROM messages`)
	require.NoError(t, err, "не могу очистить базу сообщений")

	_, err = conn.Exec(`
		INSERT INTO
			messages
		VALUES
			('msg-id-1','sub-id-1','customer-id-1','msg_text','+79163332211','done','2023-11-20 11:25:06.635877+00'),
			('msg-id-2','sub-id-1','customer-id-2','msg_text','+79163332211','done','2023-11-20 11:25:06.635877+00'),
			('msg-id-3','sub-id-2','customer-id-2','msg_text','+79163332211','done','2023-11-20 11:25:06.635877+00')
			;`)

	require.NoError(t, err, "не могу наполнить базу сообщениями")

	stats := sqlrepo.NewStatistics(conn)

	// фильтруем все сообщения
	st, err := stats.Filter(context.TODO(), "%", "%", "%")
	assert.NoError(t, err)
	assert.Equal(t, app.Statistics{"done": 3}, st)

	st, err = stats.All(context.TODO())
	assert.NoError(t, err)
	assert.Equal(t, app.Statistics{"done": 3}, st)

	// только по первой подписке
	st, err = stats.Filter(context.TODO(), "sub-id-1", "%", "%")
	assert.NoError(t, err)
	assert.Equal(t, app.Statistics{"done": 2}, st)

	// только по первому клиенту
	st, err = stats.Filter(context.TODO(), "%", "customer-id-1", "%")
	assert.NoError(t, err)
	assert.Equal(t, app.Statistics{"done": 1}, st)

}

// тут нужно сделать билд тег интеграция и сделать нормальную сюиту котрая теперь даждым тестом чистила базу бы
