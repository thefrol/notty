package customers

import (
	"database/sql"
	"fmt"

	"gitlab.com/thefrol/notty/internal/entity"
)

// Customers это репозиторий для сущности Customer то есть для нашего клиента
type Customers struct {
	db *sql.DB
}

func New(db *sql.DB) Customers {
	return Customers{
		db: db,
	}
}

func (c Customers) Get(id string) (res entity.Customer, err error) {
	r := c.db.QueryRow(`
		SELECT
			id,
			name,
			phone,
			operator,
			tag
		FROM
			Customer
		WHERE
			id=$1`, id)

	err = r.Scan(&res.Id, &res.Name, &res.Phone, &res.Operator, &res.Tag)
	return
}

func (c Customers) Delete(id string) error {
	rs, err := c.db.Exec(`
		DELETE
		FROM
			Customer
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

func (c Customers) Update(cl entity.Customer) (res entity.Customer, err error) {
	r, err := c.db.Exec(`
		UPDATE
			Customer
		SET
			name=$2,
			phone=$3,
			operator=$4,
			tag=$5,
		WHERE
			id=$1`, cl.Id, cl.Name, cl.Phone, cl.Operator, cl.Tag)
	if err != nil {
		return entity.Customer{}, err
	}

	if rs, err := r.RowsAffected(); err != nil && rs != int64(1) {
		return entity.Customer{}, fmt.Errorf("ошибка апдейта клиента %w", err)
	}

	return c.Get(cl.Id)
}

func (c Customers) Create(cl entity.Customer) (res entity.Customer, err error) {
	r, err := c.db.Exec(`
		INSERT INTO
			Customer(
				id,
				name,
				phone,
				operator,
				tag
			)
		VALUES($1,$2,$3,$4,$5)`, cl.Id, cl.Name, cl.Phone, cl.Operator, cl.Tag)
	if err != nil {
		return entity.Customer{}, err
	}

	if rs, err := r.RowsAffected(); err != nil && rs != int64(1) {
		return entity.Customer{}, fmt.Errorf("ошибка создания клиента %w", err)
	}

	return c.Get(cl.Id)
}
