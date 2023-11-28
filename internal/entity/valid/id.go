package valid

func Id(id string) error {
	if !idPattern.MatchString(id) {
		return ErrorIdValidation
	}
	return nil
}
