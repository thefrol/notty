package validate

import (
	"errors"
	"regexp"
)

var (
	phonePattern         = regexp.MustCompile(`\+\d{10,15}`) // +0 123 45 67 89 и еще пара цифр на случай Беларуса например
	ErrorPhoneValidation = errors.New("телефон должен быть вида +71113334455")

	idPattern         = regexp.MustCompile(`^[a-zA-Z0-9-]+$`)
	ErrorIdValidation = errors.New("идентификатор должен состоять из цифр, латинских букв и тире")

	ExplicitPattern       = regexp.MustCompile(`(пизд\W|хуй|хуя)`) // :))
	ErrorExplicitLanguage = errors.New("содержит мат")

	ErrorInvalidPeriod = errors.New("начало интервала должно быть раньше конца интервала")
)
