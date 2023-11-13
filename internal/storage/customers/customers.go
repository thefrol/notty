package customers

import (
	"database/sql"
	"fmt"
	"log"

	"gitlab.com/thefrol/notty/internal/entity"
	"gitlab.com/thefrol/notty/internal/storage/scan"
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

	return scan.Client(r)
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
			tag=$5
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

// Filter возвращает всех клиентов для определенной рассылки
// поскольку таких сообщений может быть ну очень много, то
// то возвращает канал
func (c Customers) Filter(tag string, operator string, size int) (chan entity.Customer, error) {

	rs, err := c.db.Query(`
		SELECT
			id,
			name,
			phone,
			operator,
			tag	
		FROM
			customer c
		WHERE
			tag ilike $1
			AND operator ilike $2
		`, tag, operator)

	if err != nil {
		return nil, err
	}
	in := make(chan entity.Customer, size)

	go func() {
		defer rs.Close()

		for rs.Next() {
			c, err := scan.Client(rs)
			if err != nil {
				log.Printf("Ошибка при чтении строки запроса %v\n", err)
				return
			}
			in <- c

		}

		if err := rs.Err(); err != nil {
			log.Printf("Ошибка при чтении запроса %v\n", err)
		}
		close(in)
	}()

	return in, nil
}
