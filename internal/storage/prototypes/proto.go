package protootypes

import (
	"database/sql"
	"log"

	"gitlab.com/thefrol/notty/internal/dto"
	"gitlab.com/thefrol/notty/internal/storage/scan"
)

const (
	StatusCreated = "created"
)

type Prototypes struct {
	db *sql.DB
}

func New(db *sql.DB) Prototypes {
	return Prototypes{
		db: db,
	}
}

// Spawn генерирует n сообщений, которые можно отправить
func (p *Prototypes) Spawn(n int) ([]dto.Prototype, error) {
	// todo думаю можно и без ошибки обойтись- просто закрыть канал
	//id = sql.NullString

	rs, err := p.db.Query(`
			SELECT
				c.id,
				s.id,
				s.msg_text,
				c.phone
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
			LIMIT
				$1
			`, n) // все статусы уже отдельно обработаем, тут выкинем только неотправленные

	if err != nil {
		return nil, err
	}
	defer rs.Close()

	//обрабатываем запрос
	batch := make([]dto.Prototype, 0, n)
	for rs.Next() {
		msg, err := scan.Prototype(rs)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		// todo add constraint
		batch = append(batch, msg)
	}

	if err := rs.Err(); err != nil {
		return nil, err
	}

	return batch, nil
}
