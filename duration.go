package humanize

import (
	"fmt"
	"math"
	"strconv"
	"strings"
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

// ParseDuration parses s when it is a canonical Duration result.
// It accepts day units in addition to time.ParseDuration units and returns
// ErrInvalid for malformed or non-canonical input.
func ParseDuration(s string) (time.Duration, error) {
	if s == "" {
		return 0, fmt.Errorf("empty string: %w", ErrInvalid)
	}
	if s != strings.TrimSpace(s) {
		return 0, fmt.Errorf("leading or trailing whitespace: %w", ErrInvalid)
	}
	original := s

	neg := false
	if rest, ok := strings.CutPrefix(s, "-"); ok {
		neg = true
		s = rest
	} else if _, ok := strings.CutPrefix(s, "+"); ok {
		return 0, fmt.Errorf("leading plus sign is not supported: %w", ErrInvalid)
	}

	first, second, hasSecond := strings.Cut(s, " ")
	if hasSecond && (second == "" || strings.Contains(second, " ")) {
		return 0, fmt.Errorf("invalid duration: %w", ErrInvalid)
	}

	total := int64(0)
	prevRank := -1
	parsePart := func(part string) error {
		nanos, rank, err := parseDurationPart(part)
		if err != nil {
			return err
		}
		if nanos == 0 && part != "0s" {
			return fmt.Errorf("zero value must be expressed as 0s: %w", ErrInvalid)
		}
		if prevRank >= 0 && rank <= prevRank {
			return fmt.Errorf("units must be in descending order: %w", ErrInvalid)
		}
		if nanos > math.MaxInt64-total {
			return fmt.Errorf("value out of range: %w", ErrInvalid)
		}
		total += nanos
		prevRank = rank
		return nil
	}

	if err := parsePart(first); err != nil {
		return 0, err
	}
	if hasSecond {
		if err := parsePart(second); err != nil {
			return 0, err
		}
	}

	value := time.Duration(total)
	if neg {
		value = -value
	}
	if Duration(value) != original {
		return 0, fmt.Errorf("non-canonical duration form: %w", ErrInvalid)
	}

	return value, nil
}

func parseDurationPart(part string) (int64, int, error) {
	split := len(part)
	for i := range len(part) {
		if part[i] < '0' || part[i] > '9' {
			split = i
			break
		}
	}
	if split == 0 || split == len(part) {
		return 0, 0, fmt.Errorf("unknown unit in %q: %w", part, ErrInvalid)
	}

	valueStr := part[:split]
	unit := part[split:]
	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid value in %q: %s: %w", part, err.Error(), ErrInvalid)
	}

	for rank, u := range durationUnits() {
		if u.name != unit {
			continue
		}

		sizeNanos := int64(u.size)
		if value > math.MaxInt64/sizeNanos {
			return 0, 0, fmt.Errorf("value out of range: %w", ErrInvalid)
		}
		return value * sizeNanos, rank, nil
	}

	return 0, 0, fmt.Errorf("unknown unit in %q: %w", part, ErrInvalid)
}
