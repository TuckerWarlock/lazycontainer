package i18n

// GetTranslationSet returns the English translation set
func GetTranslationSet() *TranslationSet {
	set := englishSet()
	return &set
}
