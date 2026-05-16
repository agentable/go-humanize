package humanize

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

const (
	kilobyte = 1000
	megabyte = kilobyte * 1000
	gigabyte = megabyte * 1000
	terabyte = gigabyte * 1000
	petabyte = terabyte * 1000
	exabyte  = petabyte * 1000
)

const (
	_        = 1 << (iota * 10)
	kibibyte // 1 << 10
	mebibyte
	gibibyte
	tebibyte
	pebibyte
	exbibyte
)

// Bytes formats b with decimal byte units such as "1.5 KB" and "2 MB".
// It preserves the sign and saturates math.MinInt64 to math.MaxInt64 when
// taking the absolute value.
func Bytes(b int64) string {
	return formatBytes(b, 1000, decimalByteUnit)
}

// BinaryBytes formats b with IEC binary byte units such as "1.5 KiB" and
// "2 MiB". It preserves the sign and saturates math.MinInt64 to math.MaxInt64
// when taking the absolute value.
func BinaryBytes(b int64) string {
	return formatBytes(b, 1024, binaryByteUnit)
}

func formatBytes(b int64, base int64, unitName func(int) string) string {
	neg := b < 0
	if b == math.MinInt64 {
		b = math.MaxInt64
	} else if neg {
		b = -b
	}

	sign := ""
	if neg {
		sign = "-"
	}

	if b < base {
		return sign + strconv.FormatInt(b, 10) + " B"
	}

	exp := 0
	value := float64(b)
	for value >= float64(base) && exp < 6 {
		value /= float64(base)
		exp++
	}

	precision := bytePrecision(value)
	if exp < 6 && roundedByteValue(value, precision) >= float64(base) {
		value /= float64(base)
		exp++
		precision = bytePrecision(value)
	}
	number := strconv.FormatFloat(value, 'f', precision, 64)

	return sign + number + " " + unitName(exp)
}

func bytePrecision(value float64) int {
	rounded := math.Round(value*10) / 10
	if value >= 10 || rounded == math.Trunc(rounded) {
		return 0
	}
	return 1
}

func roundedByteValue(value float64, precision int) float64 {
	if precision == 0 {
		return math.Round(value)
	}
	return math.Round(value*10) / 10
}

func decimalByteUnit(exp int) string {
	switch exp {
	case 0:
		return "B"
	case 1:
		return "KB"
	case 2:
		return "MB"
	case 3:
		return "GB"
	case 4:
		return "TB"
	case 5:
		return "PB"
	default:
		return "EB"
	}
}

func binaryByteUnit(exp int) string {
	switch exp {
	case 0:
		return "B"
	case 1:
		return "KiB"
	case 2:
		return "MiB"
	case 3:
		return "GiB"
	case 4:
		return "TiB"
	case 5:
		return "PiB"
	default:
		return "EiB"
	}
}

// ParseBytes parses s when it is a canonical Bytes or BinaryBytes result.
// It accepts only the units B, KB, MB, GB, TB, PB, EB, KiB, MiB, GiB, TiB,
// PiB, and EiB, with a single space between the number and unit.
// It returns ErrInvalid for malformed or non-canonical input. The parsed
// value represents the nearest byte described by the display text, not the
// original source value before formatting.
func ParseBytes(s string) (int64, error) {
	if s == "" {
		return 0, fmt.Errorf("%w: empty string", ErrInvalid)
	}
	if s != strings.TrimSpace(s) {
		return 0, fmt.Errorf("%w: leading or trailing whitespace", ErrInvalid)
	}

	numberPart, unitPart, ok := strings.Cut(s, " ")
	if !ok || numberPart == "" || unitPart == "" || strings.Contains(unitPart, " ") {
		return 0, fmt.Errorf("%w: expected \"<number> <unit>\"", ErrInvalid)
	}
	multiplier, format, ok := byteUnit(unitPart)
	if !ok {
		return 0, fmt.Errorf("%w: unknown unit %q", ErrInvalid, unitPart)
	}

	f, err := strconv.ParseFloat(numberPart, 64)
	if err != nil {
		return 0, fmt.Errorf("%w: %w", ErrInvalid, err)
	}

	result := math.Round(f * float64(multiplier))
	var value int64
	switch {
	case math.IsNaN(result):
		return 0, fmt.Errorf("%w: invalid number", ErrInvalid)
	case result >= float64(math.MaxInt64):
		value = math.MaxInt64
	case result <= float64(math.MinInt64):
		value = math.MinInt64
	default:
		value = int64(result)
	}
	if unitPart == "B" {
		if Bytes(value) != s && BinaryBytes(value) != s {
			return 0, fmt.Errorf("%w: non-canonical byte form", ErrInvalid)
		}
		return value, nil
	}
	if format(value) != s {
		return 0, fmt.Errorf("%w: non-canonical byte form", ErrInvalid)
	}

	return value, nil
}

func byteUnit(unit string) (int64, func(int64) string, bool) {
	switch unit {
	case "B":
		return 1, Bytes, true
	case "KB":
		return kilobyte, Bytes, true
	case "MB":
		return megabyte, Bytes, true
	case "GB":
		return gigabyte, Bytes, true
	case "TB":
		return terabyte, Bytes, true
	case "PB":
		return petabyte, Bytes, true
	case "EB":
		return exabyte, Bytes, true
	case "KiB":
		return kibibyte, BinaryBytes, true
	case "MiB":
		return mebibyte, BinaryBytes, true
	case "GiB":
		return gibibyte, BinaryBytes, true
	case "TiB":
		return tebibyte, BinaryBytes, true
	case "PiB":
		return pebibyte, BinaryBytes, true
	case "EiB":
		return exbibyte, BinaryBytes, true
	default:
		return 0, nil, false
	}
}
