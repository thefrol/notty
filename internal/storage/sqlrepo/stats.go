// Этот пакет отвечает за сбор статистики, содержит репозиторий, сервис
// и обьект данных
//
// В этом модуле я просто пробую совершенно другую семантику работы с данными.
// все, что касается одной какой-то сущности или данных будет лежать в одной папке
package sqlrepo

import (
	"context"
	"database/sql"

	"gitlab.com/thefrol/notty/internal/app"
)

// Statistics отвечает за сбор статистики
// по сообщениям. Тут не хочется разделять
// на несколько слоев эту логи со статистикой, так
// что по сути у нас сервис-репозиторий
type Statistics struct {
	db *sql.DB
}

func NewStatistics(db *sql.DB) *Statistics {
	return &Statistics{
		db: db,
	}
}

// All возвращает статистику по всем вообщениям,
func (r Statistics) All(ctx context.Context) (app.Statistics, error) {
	rs, err := r.db.QueryContext(ctx, `
	SELECT
		status,
		COUNT(status)
	FROM
		messages
	WHERE
		status IS NOT NULL
	GROUP BY
		status`)
	if err != nil {
		return nil, err
	}
	defer rs.Close()

	res := make(map[string]int, 10) // todo magic number
	for rs.Next() {
		var (
			status string
			count  int
		)
		err = rs.Scan(&status, &count)
		if err != nil {
			return nil, err
		}
		res[status] = count
	}

	if err := rs.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

// Filters возвращает статистику по рассылкам
func (r Statistics) Filter(ctx context.Context, subId string, customerId string, status string) (app.Statistics, error) {
	rs, err := r.db.QueryContext(ctx, `
	SELECT
		status,
		COUNT(status)
	FROM
		messages
	WHERE
		status IS NOT NULL
		AND sub_id LIKE $1
		AND customer_id LIKE $2
		AND status ILIKE $3
	GROUP BY
		status`, subId, customerId, status)
	if err != nil {
		return nil, err
	}
	defer rs.Close()

	res := make(map[string]int, 10) // todo magic number
	for rs.Next() {
		var (
			status string
			count  int
		)
		err = rs.Scan(&status, &count)
		if err != nil {
			return nil, err
		}
		res[status] = count
	}

	if err := rs.Err(); err != nil {
		return nil, err
	}

	return res, nil
}
