// Пакет помощник для создания струкур после запросов селект
//
// Пример:
//
//	for rows.Next(){
//		 entity=scan.Subscription(rows)
//	}
package scan

import "gitlab.com/thefrol/notty/internal/entity"

type Scanner interface {
	Scan(...any) error
}

func Subscription(r Scanner) (entity.Subscription, error) {
	res := entity.Subscription{}

	err := r.Scan(
		&res.Id,
		&res.Text,
		&res.Start,
		&res.End,
		&res.OperatorFilter,
		&res.PhoneFilter,
		&res.TagFilter,
		&res.Desc)
	if err != nil {
		return entity.Subscription{}, err
	}

	return res, nil
}

func Client(r Scanner) (entity.Customer, error) {
	res := entity.Customer{}
	err := r.Scan(&res.Id, &res.Name, &res.Phone, &res.Operator, &res.Tag)
	if err != nil {
		return entity.Customer{}, err
	}

	return res, nil
}
