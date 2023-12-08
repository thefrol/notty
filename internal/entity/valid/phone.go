package valid

// Phone - проверяет что введен правильный номер телефона
func Phone(phone string) error {
	if !phonePattern.MatchString(phone) {
		return ErrorPhoneValidation
	}
	return nil
}
