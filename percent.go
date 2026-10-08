package humanize

import (
	"math"
	"strconv"
	"strings"
)

// Percent formats fraction as a percentage with at most one decimal place.
// The input is a fraction, so 0.333 renders as "33.3%" and 1 renders as
// "100%". Values above 1 are not clamped: 1.25 renders as "125%".
// Rounded percentages of magnitude 1e6 or more use scientific notation.
// NaN and infinities render as "NaN%", "+Inf%", and "-Inf%".
// Rounding uses the input's shortest round-trip decimal representation,
// with ties away from zero: 0.5005 renders as "50.1%".
func Percent(fraction float64) string {
	switch {
	case math.IsNaN(fraction):
		return "NaN%"
	case math.IsInf(fraction, 1):
		return "+Inf%"
	case math.IsInf(fraction, -1):
		return "-Inf%"
	}

	// Scientific form exposes the decimal digits without scaling the float.
	mantissa, exponent, _ := strings.Cut(strconv.FormatFloat(math.Abs(fraction), 'e', -1, 64), "e")
	digits := strings.ReplaceAll(mantissa, ".", "")
	// FormatFloat produces valid digits and an exponent; the shortest float64
	// representation has at most 17 significant digits, which fit in uint64.
	coefficient, _ := strconv.ParseUint(digits, 10, 64)
	exp, _ := strconv.Atoi(exponent)
	if exp < -4 || coefficient == 0 {
		return "0%"
	}

	var tenths uint64
	if exp < 4 {
		tenths = roundPercent(coefficient, exp+4-len(digits))
	}
	suffix := "%"
	if exp >= 4 || tenths >= 10_000_000 {
		// Round the original digits, not the fixed-format rounded result.
		tenths = roundPercent(coefficient, 2-len(digits))
		if tenths == 100 {
			tenths = 10
			exp++
		}
		suffix = "e+" + strconv.Itoa(exp+2) + "%"
	}
	if tenths == 0 {
		return "0%"
	}

	sign := ""
	if fraction < 0 {
		sign = "-"
	}
	number := strconv.FormatUint(tenths/10, 10)
	if tenths%10 != 0 {
		number += "." + strconv.FormatUint(tenths%10, 10)
	}
	return sign + number + suffix
}

// roundPercent shifts decimal digits and rounds once to an integer.
// Callers bound the power to 1e17 and the result to 1e7.
func roundPercent(coefficient uint64, shift int) uint64 {
	places := shift
	if places < 0 {
		places = -places
	}
	power := uint64(1)
	for range places {
		power *= 10
	}
	if shift >= 0 {
		return coefficient * power
	}
	return (coefficient + power/2) / power
}
