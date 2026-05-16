package humanize

import (
	"strconv"
)

// Ordinal formats n with an English ordinal suffix such as "1st", "23rd",
// and "-3rd".
func Ordinal(n int64) string {
	numStr := strconv.FormatInt(n, 10)

	mod100 := n % 100
	if mod100 < 0 {
		mod100 = -mod100
	}
	if mod100 >= 11 && mod100 <= 13 {
		return numStr + "th"
	}

	mod10 := n % 10
	if mod10 < 0 {
		mod10 = -mod10
	}
	suffix := ordinalSuffix(mod10)

	return numStr + suffix
}

func ordinalSuffix(mod10 int64) string {
	switch mod10 {
	case 1:
		return "st"
	case 2:
		return "nd"
	case 3:
		return "rd"
	default:
		return "th"
	}
}
