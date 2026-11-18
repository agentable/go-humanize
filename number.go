package humanize

import (
	"strconv"
	"strings"
)

// Number formats n as a human-readable English integer with comma separators
// every three digits, such as "1,234,567" and "-1,000,000".
func Number(n int64) string {
	text := strconv.FormatInt(n, 10)
	digits, negative := strings.CutPrefix(text, "-")
	if len(digits) <= 3 {
		return text
	}

	numCommas := (len(digits) - 1) / 3
	result := make([]byte, len(text)+numCommas)
	pos := len(result) - 1
	for i := len(digits) - 1; i >= 0; i-- {
		result[pos] = digits[i]
		pos--

		if (len(digits)-i)%3 == 0 && i > 0 {
			result[pos] = ','
			pos--
		}
	}

	if negative {
		result[0] = '-'
	}

	return string(result)
}
