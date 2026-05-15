package humanize

import (
	"strconv"
	"time"
)

// Relative formats target relative to ref using approximate 30-day months and
// 365-day years. Durations under a second return "just now". The output
// describes target's position relative to ref: target before ref renders as
// "N units ago", and target after ref renders as "in N units".
func Relative(target, ref time.Time) string {
	diff := ref.Sub(target)
	if diff < time.Second && diff > -time.Second {
		return "just now"
	}

	future := diff < 0
	if future {
		diff = absNegativeDuration(diff)
	}

	const (
		minute = 60 * time.Second
		hour   = 60 * minute
		day    = 24 * hour
		week   = 7 * day
		month  = 30 * day
		year   = 365 * day
		// Prefer "1 year" over "12 months" near a full year.
		yearCutover = 345 * day
	)

	var value int64
	var unit string
	switch {
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
	case diff < yearCutover:
		value = int64(diff / month)
		unit = "month"
	default:
		value = max(int64(1), int64(diff/year))
		unit = "year"
	}

	unit = pluralForm(value, unit, unit+"s")
	valueStr := strconv.FormatInt(value, 10)
	if future {
		return "in " + valueStr + " " + unit
	}
	return valueStr + " " + unit + " ago"
}
