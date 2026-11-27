package humanize

import "strconv"

// Number formats n as a human-readable English integer with comma separators
// every three digits, such as "1,234,567" and "-1,000,000".
func Number(n int64) string {
	if n > -1000 && n < 1000 {
		return strconv.FormatInt(n, 10)
	}

	// An int64 needs at most 19 digits, a sign, and six separators.
	var decimal [20]byte
	var formatted [26]byte
	digits := strconv.AppendInt(decimal[:0], n, 10)
	result := formatted[:0]
	if digits[0] == '-' {
		result = append(result, '-')
		digits = digits[1:]
	}

	firstGroup := (len(digits)-1)%3 + 1
	result = append(result, digits[:firstGroup]...)
	for i := firstGroup; i < len(digits); i += 3 {
		result = append(result, ',')
		result = append(result, digits[i:i+3]...)
	}

	return string(result)
}
