// Используется для валидации данных. Запросы и ответы разделены специально
package validate

import (
	"errors"

	"gitlab.com/thefrol/notty/internal/entity"
)

// CustomerRequest валидирует клиента
func CustomerRequest(c entity.Customer) error {
	var errs []error

	if !phonePattern.MatchString(c.Phone) {
		errs = append(errs, ErrorPhoneValidation)
	}

	if !idPattern.MatchString(c.Id) {
		errs = append(errs, ErrorIdValidation)
	}

	// todo проверка оператора по МТС, билайн, мегафон это цифра или короткий код
	return errors.Join(errs...)
}
