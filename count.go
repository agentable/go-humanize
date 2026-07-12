package humanize

func pluralForm(count int64, singular, plural string) string {
	if count == 1 || count == -1 {
		return singular
	}
	return plural
}

// Count formats count with Number's comma separators and chooses singular
// when count is ±1.
func Count(count int64, singular, plural string) string {
	form := pluralForm(count, singular, plural)
	if form == "" {
		return Number(count)
	}
	return Number(count) + " " + form
}
