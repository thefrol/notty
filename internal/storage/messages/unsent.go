package messages

import (
	"log"

	"gitlab.com/thefrol/notty/internal/entity"
	"gitlab.com/thefrol/notty/internal/storage/scan"
)

// NewMessages генерирует сообщения, которые можно отправить
func (m *Messages) NewMessages(buffer int) (chan entity.Message, error) {
	// todo думаю можно и без ошибки обойтись- просто закрыть канал
	//id = sql.NullString

	rs, err := m.db.Query(`
			SELECT
				COALESCE(m.id,''),
				c.id,
				s.id,
				s.msg_text,
				c.phone,
				m.status,
				m.sent
			FROM
				subscription s,
				customer c
			LEFT JOIN
				messages m
				ON m.customer_id=c.id
			WHERE
				c.tag ILIKE s.tag_filter
				AND c.operator ILIKE s.operator_filter
				AND m.status IS NULL
			`) // все статусы уже отдельно обработаем, тут выкинем только неотправленные

	if err != nil {
		return nil, err
	}
	in := make(chan entity.Message, buffer)

	go func() {
		defer rs.Close()

		defer close(in)

		for rs.Next() {
			m, err := scan.Message(rs)
			if err != nil {
				log.Println(err)
				return
			}
			// Должны из этой функции уже выходить как created или не выходить,
			// и вообще мы должны дождаться прежде чем читать заново, что
			// мы прошли до конца уже
			in <- m

		}

		if err := rs.Err(); err != nil {
			log.Printf("Ошибка при чтении запроса %v\n", err)
		}

	}()

	return in, nil
}
