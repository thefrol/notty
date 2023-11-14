// коллекция разных сырых запросов
package messages

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
	"gitlab.com/thefrol/notty/internal/entity"
	"gitlab.com/thefrol/notty/internal/storage/scan"
)

type Messages struct {
	db *sql.DB
}

func New(db *sql.DB) Messages {
	return Messages{
		db: db,
	}
}

func (m Messages) Create(ms entity.Message) (entity.Message, error) {
	if ms.Id != "" {
		log.Println("Создается сообщение с ненулевым айдишником")
	}
	id := uuid.New().String()
	r, err := m.db.Exec(`
		INSERT INTO
			Messages(
				id,
				sub_id,
				customer_id,
				message_text,
				phone,
				status,
				sent
			)
		VALUES($1,$2,$3,$4,$5,$6,$7)`, id,
		ms.SubscriptionId, ms.CustomerId, ms.Text, ms.Phone,
		ms.Status, ms.Sent)
	if err != nil {
		return entity.Message{}, err
	}

	if rs, err := r.RowsAffected(); err != nil && rs != int64(1) {
		return entity.Message{}, fmt.Errorf("ошибка создания сообщения %w", err)
	}

	return m.Get(id)
}

func (m Messages) Get(id string) (res entity.Message, err error) {
	r := m.db.QueryRow(`
		SELECT
			id,
			customer_id,
			sub_id,
			message_text,
			phone,
			status,
			sent timestamptz
		FROM
			Messages
		WHERE
			id=$1`, id)

	return scan.Message(r)
}

func (c Messages) ByStatus(status string, n int) ([]entity.Message, error) {
	rs, err := c.db.Query(`
	SELECT
		id,
		customer_id,
		sub_id,
		message_text,
		phone,
		status,
		sent timestamptz
	FROM
		Messages
	WHERE
		status=$1`, status)
	if err != nil {
		return nil, err
	}

	//обрабатываем запрос
	batch := make([]entity.Message, 0, n)
	for rs.Next() {
		msg, err := scan.Message(rs)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		batch = append(batch, msg)
	}

	if err := rs.Err(); err != nil {
		return nil, err
	}
	return batch, nil
}

func (c Messages) Delete(id string) error {
	rs, err := c.db.Exec(`
		DELETE
		FROM
			Messages
		WHERE
			id=$1`, id)

	if err != nil {
		return err
	}

	if rs, err := rs.RowsAffected(); err != nil && rs != int64(1) {
		return fmt.Errorf("ошибка удаления сообщения %w", err)
	}

	return nil
}

func (c Messages) Update(msg entity.Message) (res entity.Message, err error) {
	r, err := c.db.Exec(`
		UPDATE
			Messages
		SET
			sub_id =$2,
			customer_id =$3,
			message_text=$4,
			phone=$5,
			status=$6,
			sent=$7
		WHERE
			id=$1`, msg.Id, msg.SubscriptionId, msg.CustomerId,
		msg.Text, msg.Phone, msg.Status, msg.Sent)
	if err != nil {
		return entity.Message{}, err
	}

	if rs, err := r.RowsAffected(); err != nil && rs != int64(1) {
		return entity.Message{}, fmt.Errorf("ошибка апдейта сообщения %w", err)
	}

	return c.Get(msg.Id)
}
