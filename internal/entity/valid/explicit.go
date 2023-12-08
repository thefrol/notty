package valid

func HasNoExplicitLanguage(text string) error {
	if ExplicitPattern.MatchString(text) {
		return ErrorExplicitLanguage
	}
	return nil
}
