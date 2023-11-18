package stats

import "database/sql"

// Service отвечает за сбор статистики
// по сообщениям. Тут не хочется разделять
// на несколько слоев эту логи со статистикой, так
// что по сути у нас сервис-репозиторий
type Service struct {
	db *sql.DB
}

func New(db *sql.DB) *Service {
	return &Service{
		db: db,
	}
}

// All возвращает статистику по всем вообщениям,
func (r Service) All() (Statistics, error) {
	rs, err := r.db.Query(`
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

// BySubscription возвращает статистику по рассылкам
func (r Service) Filter(subId string, customerId string, status string) (Statistics, error) {
	rs, err := r.db.Query(`
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
