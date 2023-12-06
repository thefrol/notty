// Тут содержится репозиторий клиентов.

package sqlrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"gitlab.com/thefrol/notty/internal/app"
	"gitlab.com/thefrol/notty/internal/entity"
	"gitlab.com/thefrol/notty/internal/storage/sqlrepo/scan"
)

var _ app.Customerere = (*Customers)(nil)

// Customers это репозиторий для сущности Customer то есть для нашего клиента
type Customers struct {
	db *sql.DB
}

func NewCustomers(db *sql.DB) Customers {
	return Customers{
		db: db,
	}
}

// presentation вот этот формат довольно красивый можно об этом рассказать
// отдельно запросы, отдельно функции

const getCustomer = `
SELECT
	id,
	name,
	phone,
	operator,
	tag
FROM
	Customer
WHERE
	id=$1`

// Get возвращает клиента с указанным id, или вощвращаем
// app.ErrorCustomerNotFound - ошибку
func (c Customers) Get(ctx context.Context, id string) (res entity.Customer, err error) {
	r := c.db.QueryRowContext(ctx, getCustomer, id)
	if err := r.Err(); err != nil {
		return entity.Customer{}, err
	}

	s, err := scan.Client(r)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Customer{}, app.ErrorCustomerNotFound
		}
		return entity.Customer{}, err
	}
	return s, nil
}

const deleteCustomer = `
DELETE
FROM
	Customer
WHERE
	id=$1`

func (c Customers) Delete(ctx context.Context, id string) error {
	rs, err := c.db.ExecContext(ctx, deleteCustomer, id)
	if err != nil {
		return err
	}

	if rs, err := rs.RowsAffected(); err != nil && rs != int64(1) {
		return fmt.Errorf("ошибка удаления клиента, считаем что он ненайден %w", err)
	}

	return nil
}

const updateCustomer = `
UPDATE
	Customer
SET
	name=$2,
	phone=$3,
	operator=$4,
	tag=$5
WHERE
	id=$1
RETURNING *`

func (c Customers) Update(ctx context.Context, cl entity.Customer) (entity.Customer, error) {
	r := c.db.QueryRowContext(ctx, updateCustomer, cl.Id, cl.Name, cl.Phone, cl.Operator, cl.Tag)
	if err := r.Err(); err != nil {
		return entity.Customer{}, err
	}

	s, err := scan.Client(r)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) { // todo правда это или нет хз
			return entity.Customer{}, app.ErrorCustomerNotFound
		}
		return entity.Customer{}, err
	}
	return s, nil
}

const createCustomer = `
INSERT INTO
	Customer(
		id,
		name,
		phone,
		operator,
		tag
	)
VALUES($1,$2,$3,$4,$5)
RETURNING *`

func (c Customers) Create(ctx context.Context, cl entity.Customer) (entity.Customer, error) {
	r := c.db.QueryRowContext(ctx, createCustomer, cl.Id, cl.Name, cl.Phone, cl.Operator, cl.Tag)
	if err := r.Err(); err != nil {
		return entity.Customer{}, nil //todo ??? ошибка же
	}

	return scan.Client(r)
}
