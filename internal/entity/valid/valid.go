// Содержит примитивы валидации
package valid

import "errors"

// If собирает несколько валидаций вместе, это такая
// функция помошник, чтобы получить семантику кода прикольную,
// можно гордится
//
// valid.If(
//
//	valid.PhoneNumber(c.Phone),
//	valid.Id(c.Id)
//
// )
func If(errs ...error) error {
	return errors.Join(errs...)
}
