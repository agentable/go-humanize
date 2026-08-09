package humanize

import (
	"math"
	"strconv"
	"time"
)

const relativeYearSeconds int64 = 365 * 24 * 60 * 60

// Relative formats target relative to ref using approximate 30-day months and
// 365-day years. Durations under a second return "just now". The output
// describes target's position relative to ref: target before ref renders as
// "N units ago", and target after ref renders as "in N units".
func Relative(target, ref time.Time) string {
	future := target.After(ref)
	earlier, later := target, ref
	if future {
		earlier, later = ref, target
	}

	diff := later.Sub(earlier)
	if diff < time.Second {
		return "just now"
	}

	const (
		minute = 60 * time.Second
		hour   = 60 * minute
		day    = 24 * hour
		week   = 7 * day
		month  = 30 * day
		year   = time.Duration(relativeYearSeconds) * time.Second
	)

	var value int64
	var unit string
	switch {
	case diff == time.Duration(math.MaxInt64):
		value = relativeYears(earlier, later)
		unit = "year"
	case diff < minute:
		value = int64(diff / time.Second)
		unit = "second"
	case diff < hour:
		value = int64(diff / minute)
		unit = "minute"
	case diff < day:
		value = int64(diff / hour)
		unit = "hour"
	case diff < week:
		value = int64(diff / day)
		unit = "day"
	case diff < month:
		value = int64(diff / week)
		unit = "week"
	case diff < year:
		value = int64(diff / month)
		unit = "month"
	default:
		value = int64(diff / year)
		unit = "year"
	}

	unit = pluralForm(value, unit, unit+"s")
	valueStr := strconv.FormatInt(value, 10)
	if future {
		return "in " + valueStr + " " + unit
	}
	return valueStr + " " + unit + " ago"
}

func relativeYears(earlier, later time.Time) int64 {
	earlierYears, earlierRemainder := splitRelativeYear(earlier.Unix())
	laterYears, laterRemainder := splitRelativeYear(later.Unix())

	years := laterYears - earlierYears
	if laterRemainder < earlierRemainder ||
		laterRemainder == earlierRemainder && later.Nanosecond() < earlier.Nanosecond() {
		years--
	}
	return years
}

func splitRelativeYear(seconds int64) (years, remainder int64) {
	years = seconds / relativeYearSeconds
	remainder = seconds % relativeYearSeconds
	if remainder < 0 {
		years--
		remainder += relativeYearSeconds
	}
	return years, remainder
}
