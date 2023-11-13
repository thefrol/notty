// Используется для валидации данных. Запросы и ответы разделены специально
package validate

import (
	"errors"

	"gitlab.com/thefrol/notty/internal/entity"
)

// Customer валидирует клиента
func SubscriptionRequest(sub entity.Subscription) error {
	var errs []error

	if !idPattern.MatchString(sub.Id) {
		errs = append(errs, ErrorIdValidation)
	}

	if sub.End.Before(sub.Start) {
		errs = append(errs, ErrorInvalidPeriod)
	}

	if ExplicitPattern.MatchString(sub.Text) {
		errs = append(errs, ErrorExplicitLanguage)
	}

	return errors.Join(errs...)
}
