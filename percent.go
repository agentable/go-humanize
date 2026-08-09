package humanize

import (
	"math"
	"strconv"
)

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

	if math.Abs(fraction) > math.MaxFloat64/1000 {
		return largePercent(fraction)
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

func largePercent(fraction float64) string {
	sign := ""
	if fraction < 0 {
		sign = "-"
	}

	value := math.Abs(fraction)
	exp := int(math.Floor(math.Log10(value)))
	mantissa := value / math.Pow10(exp)
	mantissa = math.Round(mantissa*10) / 10
	if mantissa >= 10 {
		mantissa /= 10
		exp++
	}

	precision := 1
	if mantissa == math.Trunc(mantissa) {
		precision = 0
	}

	percentExp := exp + 2
	expSign := "+"
	if percentExp < 0 {
		expSign = ""
	}

	return sign + strconv.FormatFloat(mantissa, 'f', precision, 64) + "e" + expSign + strconv.Itoa(percentExp) + "%"
}
