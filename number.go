package humanize

import (
	"math"
	"strconv"
)

// Number formats n as a human-readable English integer with comma separators
// every three digits, such as "1,234,567" and "-1,000,000".
func Number(n int64) string {
	if n == math.MinInt64 {
		// -math.MinInt64 overflows.
		return "-9,223,372,036,854,775,808"
	}

	neg := n < 0
	if neg {
		n = -n
	}

	s := strconv.FormatInt(n, 10)
	sign := ""
	if neg {
		sign = "-"
	}
	if len(s) <= 3 {
		return sign + s
	}

	numCommas := (len(s) - 1) / 3
	resultLen := len(sign) + len(s) + numCommas

	result := make([]byte, resultLen)
	pos := resultLen - 1
	for i := len(s) - 1; i >= 0; i-- {
		result[pos] = s[i]
		pos--

		if (len(s)-i)%3 == 0 && i > 0 {
			result[pos] = ','
			pos--
		}
	}

	copy(result, sign)

	return string(result)
}

// Percent formats fraction as a percentage with at most one decimal place.
// The input is a fraction, so 0.333 renders as "33.3%" and 1 renders as
// "100%". Values above 1 are not clamped: 1.25 renders as "125%".
// NaN and infinities render as "NaN%", "+Inf%", and "-Inf%".
func Percent(fraction float64) string {
	switch {
	case math.IsNaN(fraction):
		return "NaN%"
	case math.IsInf(fraction, 1):
		return "+Inf%"
	case math.IsInf(fraction, -1):
		return "-Inf%"
	}

	pct := math.Round(fraction*1000) / 10
	if pct == 0 {
		return "0%"
	}
	if pct == math.Trunc(pct) {
		return strconv.FormatFloat(pct, 'f', 0, 64) + "%"
	}
	return strconv.FormatFloat(pct, 'f', 1, 64) + "%"
}
