// коллекция разных сырых запросов
package sqlrepo

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
	"gitlab.com/thefrol/notty/internal/entity"
	"gitlab.com/thefrol/notty/internal/storage/sqlrepo/scan"
)

type Messages struct {
	db *sql.DB
}

func NewMessages(db *sql.DB) Messages {
	return Messages{
		db: db,
	}
}

// deprecated
// todo delete
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

// Spawn создает сообщения, блокируя таблицу в базе. Таким образом
// при двух параллельных вызовах этого метода, два разных процесса получат
// разные сообщения гарантированно
//
// Классная функция, конечно, дофига работы выполняет. Но насколько она гибкая?
// средне. И тестировать её просто офигеешь. То есть это прям надо несколько раз
// забивать базу, и проверять, чтобы все работало
func (m Messages) Spawn(n int, status string) ([]entity.Message, error) {
	// Я реально думаю, что транзакция должа быть тут, потому что реализация
	// хранилища не должна вылезать за этот слой. Что если у нас будет
	// не транзакционная БД? Что если у нас вообще будет не БД
	tx, err := m.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// теперь заблокируем таблицу, чтобы никто
	// другой не смог создать такие же сообщения.
	// да и даже не смог бы начать их искать.
	// Просто сейчас нельзя!
	_, err = tx.Exec("LOCK TABLE Messages") // todo выбрать какую-то менее строгую блокировку
	if err != nil {
		return nil, err
	}

	rs, err := tx.Query(`
		WITH
			proto AS(
				SELECT
					c.id AS customer_id,
					s.id AS subscription_id,
					s.msg_text AS message_text,
					c.phone AS phone
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
					AND s.sub_start < now() 
					AND now() < s.sub_end 
				LIMIT
					$1)

		INSERT INTO
			Messages(
				id,
				customer_id,
				sub_id,
				message_text,
				phone,
				status,
				sent
			)
		SELECT
			gen_random_uuid(),
			customer_id,
			subscription_id,
			message_text,
			phone,
			$2,
			NULL
		FROM 
			proto
		RETURNING
			id,
			customer_id,
			sub_id,
			message_text,
			phone,
			status,
			sent ;`, n, status)
	if err != nil {
		return nil, err
	}

	if err := rs.Err(); err != nil {
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

	// коммитим
	// комит должен стоять в самом конце,
	// иначе он тупо отменит контекст запроса
	// и rs.Err() выдаст ошибку
	// insight
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return batch, nil
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
