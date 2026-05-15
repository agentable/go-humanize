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

var decimalUnits = []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}
var binaryUnits = []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB"}

type byteUnitSpec struct {
	multiplier int64
	format     func(int64) string
}

var byteUnits = map[string]byteUnitSpec{
	"B":   {multiplier: 1, format: Bytes},
	"KB":  {multiplier: kilobyte, format: Bytes},
	"MB":  {multiplier: megabyte, format: Bytes},
	"GB":  {multiplier: gigabyte, format: Bytes},
	"TB":  {multiplier: terabyte, format: Bytes},
	"PB":  {multiplier: petabyte, format: Bytes},
	"EB":  {multiplier: exabyte, format: Bytes},
	"KiB": {multiplier: kibibyte, format: BinaryBytes},
	"MiB": {multiplier: mebibyte, format: BinaryBytes},
	"GiB": {multiplier: gibibyte, format: BinaryBytes},
	"TiB": {multiplier: tebibyte, format: BinaryBytes},
	"PiB": {multiplier: pebibyte, format: BinaryBytes},
	"EiB": {multiplier: exbibyte, format: BinaryBytes},
}

// Bytes formats b with decimal byte units such as "1.5 KB" and "2 MB".
// It preserves the sign and saturates math.MinInt64 to math.MaxInt64 when
// taking the absolute value.
func Bytes(b int64) string {
	return formatBytes(b, 1000, decimalUnits)
}

// BinaryBytes formats b with IEC binary byte units such as "1.5 KiB" and
// "2 MiB". It preserves the sign and saturates math.MinInt64 to math.MaxInt64
// when taking the absolute value.
func BinaryBytes(b int64) string {
	return formatBytes(b, 1024, binaryUnits)
}

func formatBytes(b int64, base int64, units []string) string {
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
	for value >= float64(base) && exp < len(units)-1 {
		value /= float64(base)
		exp++
	}

	precision := 1
	rounded := math.Round(value*10) / 10
	if value >= 10 || rounded == math.Trunc(rounded) {
		precision = 0
	}
	number := strconv.FormatFloat(value, 'f', precision, 64)

	return sign + number + " " + units[exp]
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
	unit, ok := byteUnits[unitPart]
	if !ok {
		return 0, fmt.Errorf("%w: unknown unit %q", ErrInvalid, unitPart)
	}

	f, err := strconv.ParseFloat(numberPart, 64)
	if err != nil {
		return 0, fmt.Errorf("%w: %w", ErrInvalid, err)
	}

	result := math.Round(f * float64(unit.multiplier))
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
	if unit.format(value) != s {
		return 0, fmt.Errorf("%w: non-canonical byte form", ErrInvalid)
	}

	return value, nil
}
