package humanize

import (
	"math"
	"strconv"
	"time"
)

type durationUnit struct {
	name string
	size time.Duration
}

type durationUnitSet [6]durationUnit

func durationUnits() durationUnitSet {
	return durationUnitSet{
		{name: "d", size: 24 * time.Hour},
		{name: "h", size: time.Hour},
		{name: "m", size: time.Minute},
		{name: "s", size: time.Second},
		{name: "ms", size: time.Millisecond},
		{name: "µs", size: time.Microsecond},
	}
}

func absNegativeDuration(d time.Duration) time.Duration {
	if d == time.Duration(math.MinInt64) {
		return time.Duration(math.MaxInt64)
	}
	return -d
}

// Duration formats d with up to two significant units such as "1h 30m" and
// "2d 5h". It preserves the sign, treats a day as 24 hours, uses ms and µs
// below a second, and returns "0s" for zero.
func Duration(d time.Duration) string {
	if d == 0 {
		return "0s"
	}

	neg := d < 0
	if neg {
		d = absNegativeDuration(d)
	}

	result := ""
	partsCount := 0
	remaining := d

	for _, u := range durationUnits() {
		if remaining < u.size {
			continue
		}

		count := remaining / u.size
		part := strconv.FormatInt(int64(count), 10) + u.name
		if result == "" {
			result = part
		} else {
			result += " " + part
		}
		partsCount++
		if partsCount == 2 {
			break
		}
		remaining -= count * u.size
	}

	if result == "" {
		return "0s"
	}
	if neg {
		return "-" + result
	}
	return result
}
