// коллекция разных сырых запросов
package sqlrepo

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gitlab.com/thefrol/notty/internal/entity"
	"gitlab.com/thefrol/notty/internal/storage/sqlrepo/scan"
)

type Messages struct {
	db     *sql.DB
	logger zerolog.Logger
}

func NewMessages(db *sql.DB, logger zerolog.Logger) Messages {
	return Messages{
		db:     db,
		logger: logger,
	}
}

// ReserverFromStatus резервирует n сообщений со статусом status, и устанавливает им
// статус entity.StatusReserved
//
// todo главная проблема этого запроса, что мы хотим выполять его параллельно с
// LockedSpawn, но тот уже заблокировал таблицу
func (m Messages) ReserveFromStatus(n int, status string) ([]entity.Message, error) {
	// тут используется подзапрос, чтобы мы могли получить
	// строго определенное количество сообщений
	// LIMIT с update не работает
	rs, err := m.db.Query(`
		UPDATE
			Messages
		SET
			status=$1
		WHERE
			id IN(
				SELECT
					id
				FROM
					Messages
				WHERE
					status=$2
				LIMIT
					$3
			)
		RETURNING
			id,
			customer_id,
			sub_id,
			message_text,
			phone,
			status,
			sent timestamptz`,
		entity.StatusReserved, status, n)
	if err != nil {
		return nil, err
	}
	defer rs.Close() // todo много такого забыто

	if err := rs.Err(); err != nil {
		return nil, err
	}

	//обрабатываем запрос
	batch := make([]entity.Message, 0, n)
	for rs.Next() {
		msg, err := scan.Message(rs)
		if err != nil {
			m.logger.Error().AnErr("Ошибка обработки результата SQL запроса", err)
			return nil, err
		}
		batch = append(batch, msg)
	}

	if err := rs.Err(); err != nil {
		return nil, err
	}

	return batch, nil
}

// LockedSpawn создает сообщения, блокируя таблицу в базе. Таким образом
// при двух параллельных вызовах этого метода, два разных процесса получат
// разные сообщения гарантированно
func (m Messages) LockedSpawn(n int, status string) ([]entity.Message, error) {
	// Я реально думаю, что транзакция должа быть тут, потому что реализация
	// хранилища не должна вылезать за этот слой. Что если у нас будет
	// не транзакционная БД? Что если у нас вообще будет не БД

	// Операция довольно сложная с возможно несколькими процессами, поэтому
	// логгируем довольно обширно, сразу передаем айдишник всем причастным логгерам
	opLogger := m.logger.With().
		Str("operation", "LockedSpawn").
		Str("table", "Messages").
		Str("operation_id", uuid.NewString()).Logger()

	opLogger.Info().
		Str("status", "started").
		Msg("Сейчас будет заблокирована база данных, чтобы создать сообщения")

	tx, err := m.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		m.logger.Info().
			Str("transaction", "rollback").
			Str("status", "fail").
			Str("lock", "released").
			Msg("Откатываются изменения при генерации сообщений")

		tx.Rollback()
	}()

	// теперь заблокируем таблицу, чтобы никто
	// другой не смог создать такие же сообщения.
	// да и даже не смог бы начать их искать.
	// Просто сейчас нельзя!
	_, err = tx.Exec("LOCK TABLE Messages") // todo выбрать какую-то менее строгую блокировку
	if err != nil {
		return nil, err
	}
	opLogger.Info().
		Str("lock", "aquired").
		Msg("Начало генерации сообщений. Заблокирована таблица")

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
			opLogger.Error().AnErr("Ошибка обработки результата SQL запроса(scan())", err)
			return nil, err
		}
		batch = append(batch, msg)
	}

	if err := rs.Err(); err != nil {
		opLogger.Error().AnErr("Ошибка обработки результата SQL запроса(в строках)", err)
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
	opLogger.Info().
		Str("status", "success").
		Str("lock", "released").
		Msg("Успех. Сообщения сгенерированы. Блокировка снята")

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

func (m Messages) ByStatus(status string, n int) ([]entity.Message, error) {
	rs, err := m.db.Query(`
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
			m.logger.Error().AnErr("Ошибка обработки результата SQL запроса", err)
			return nil, err
		}
		batch = append(batch, msg)
	}

	if err := rs.Err(); err != nil {
		return nil, err
	}
	return batch, nil
}

func (m Messages) Delete(id string) error {
	rs, err := m.db.Exec(`
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

func (m Messages) Update(msg entity.Message) (res entity.Message, err error) {
	r, err := m.db.Exec(`
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

	return m.Get(msg.Id)
}
