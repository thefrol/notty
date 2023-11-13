package subscriptions

import (
	"database/sql"
	"fmt"

	"gitlab.com/thefrol/notty/internal/entity"
	"gitlab.com/thefrol/notty/internal/storage/scan"
)

// todo мда, назавание глупое, надо же рассылка а не подписка

// Subscriptions это репозиторий для сущности Subscription то есть для наших рассылок
type Subscriptions struct {
	db *sql.DB
}

func New(db *sql.DB) Subscriptions {
	return Subscriptions{
		db: db,
	}
}

func (c Subscriptions) Get(id string) (res entity.Subscription, err error) {
	r := c.db.QueryRow(`
		SELECT
			id,
			msg_text, 
			sub_start,
			sub_end,
			operator_filter,
			phone_filter,
			tag_filter,
			description
		FROM
			Subscription
		WHERE
			id=$1`, id)

	return scan.Subscription(r)
}

func (c Subscriptions) Delete(id string) error {
	rs, err := c.db.Exec(`
		DELETE
		FROM
			Subscription
		WHERE
			id=$1`, id)

	if err != nil {
		return err
	}

	if rs, err := rs.RowsAffected(); err != nil && rs != int64(1) {
		return fmt.Errorf("ошибка создания клиента, считаем что он ненайден %w", err)
	}

	return nil
}

func (c Subscriptions) Update(cl entity.Subscription) (res entity.Subscription, err error) {
	r, err := c.db.Exec(`
		UPDATE
			Subscription
		SET	
			msg_text=$2, 
			sub_start=$3,
			sub_end=$4,
			operator_filter=$5,
			phone_filter=$6,
			tag_filter=$7,
			description=$8
		WHERE
			id=$1`, cl.Id, cl.Text, cl.Start,
		cl.End, cl.OperatorFilter, cl.PhoneFilter,
		cl.TagFilter, cl.Desc)
	if err != nil {
		return entity.Subscription{}, err
	}

	if rs, err := r.RowsAffected(); err != nil && rs != int64(1) {
		return entity.Subscription{}, fmt.Errorf("ошибка апдейта клиента %w", err)
	}

	return c.Get(cl.Id)
}

//todo запросы надо скомпилировать

func (c Subscriptions) Create(cl entity.Subscription) (res entity.Subscription, err error) {
	r, err := c.db.Exec(`
		INSERT INTO
			Subscription(
				id,
				msg_text, 
				sub_start,
				sub_end,
				operator_filter,
				phone_filter,
				tag_filter,
				description
			)
		VALUES($1,$2,$3,$4,$5,$6,$7,$8)`,
		cl.Id, cl.Text, cl.Start,
		cl.End, cl.OperatorFilter, cl.PhoneFilter,
		cl.TagFilter, cl.Desc)
	if err != nil {
		return entity.Subscription{}, err
	}

	if rs, err := r.RowsAffected(); err != nil && rs != int64(1) {
		return entity.Subscription{}, fmt.Errorf("ошибка создания клиента %w", err)
	}

	return c.Get(cl.Id)
}

// Active возвращает список активных рассылок
func (s Subscriptions) Active() ([]entity.Subscription, error) {
	var subs []entity.Subscription

	rs, err := s.db.Query(`
		SELECT
			id,
			msg_text, 
			sub_start,
			sub_end,
			operator_filter,
			phone_filter,
			tag_filter,
			description
		FROM 
			subscription
		WHERE
			sub_start<now()
			AND now()<sub_end`)
	if err != nil {
		return nil, err
	}
	defer rs.Close()

	for rs.Next() {
		res, err := scan.Subscription(rs)
		if err != nil {
			return nil, err
		}

		subs = append(subs, res)
	}

	return subs, nil
}
