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

type byteScale struct {
	unit string
	size int64
}

type byteScales [7]byteScale

func decimalByteScales() byteScales {
	return byteScales{
		{unit: "B", size: 1},
		{unit: "KB", size: kilobyte},
		{unit: "MB", size: megabyte},
		{unit: "GB", size: gigabyte},
		{unit: "TB", size: terabyte},
		{unit: "PB", size: petabyte},
		{unit: "EB", size: exabyte},
	}
}

func binaryByteScales() byteScales {
	return byteScales{
		{unit: "B", size: 1},
		{unit: "KiB", size: kibibyte},
		{unit: "MiB", size: mebibyte},
		{unit: "GiB", size: gibibyte},
		{unit: "TiB", size: tebibyte},
		{unit: "PiB", size: pebibyte},
		{unit: "EiB", size: exbibyte},
	}
}

func (scales *byteScales) find(unit string) (byteScale, bool) {
	for _, scale := range scales {
		if scale.unit == unit {
			return scale, true
		}
	}
	return byteScale{}, false
}

// Bytes formats b with decimal byte units such as "1.5 KB" and "2 MB".
// It preserves the sign and saturates math.MinInt64 to math.MaxInt64 when
// taking the absolute value.
func Bytes(b int64) string {
	scales := decimalByteScales()
	return formatBytes(b, &scales)
}

// BinaryBytes formats b with IEC binary byte units such as "1.5 KiB" and
// "2 MiB". It preserves the sign and saturates math.MinInt64 to math.MaxInt64
// when taking the absolute value.
func BinaryBytes(b int64) string {
	scales := binaryByteScales()
	return formatBytes(b, &scales)
}

func formatBytes(b int64, scales *byteScales) string {
	base := scales[1].size
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

	return sign + number + " " + scales[exp].unit
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

// ParseBytes parses s when it is a canonical Bytes or BinaryBytes result.
// It accepts only the units B, KB, MB, GB, TB, PB, EB, KiB, MiB, GiB, TiB,
// PiB, and EiB, with a single space between the number and unit.
// It returns ErrInvalid for malformed or non-canonical input. The parsed
// value represents the nearest byte described by the display text, not the
// original source value before formatting.
func ParseBytes(s string) (int64, error) {
	if s == "" {
		return 0, fmt.Errorf("empty string: %w", ErrInvalid)
	}
	if s != strings.TrimSpace(s) {
		return 0, fmt.Errorf("leading or trailing whitespace: %w", ErrInvalid)
	}

	numberPart, unitPart, ok := strings.Cut(s, " ")
	if !ok || numberPart == "" || unitPart == "" || strings.Contains(unitPart, " ") {
		return 0, fmt.Errorf("expected \"<number> <unit>\": %w", ErrInvalid)
	}
	multiplier, format, ok := byteUnit(unitPart)
	if !ok {
		return 0, fmt.Errorf("unknown unit %q: %w", unitPart, ErrInvalid)
	}

	f, err := strconv.ParseFloat(numberPart, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid number %q: %w: %w", numberPart, err, ErrInvalid)
	}

	result := math.Round(f * float64(multiplier))
	var value int64
	switch {
	case math.IsNaN(result):
		return 0, fmt.Errorf("invalid number: %w", ErrInvalid)
	case result >= float64(math.MaxInt64):
		value = math.MaxInt64
	case result <= float64(math.MinInt64):
		value = math.MinInt64
	default:
		value = int64(result)
	}
	if unitPart == "B" {
		if Bytes(value) != s && BinaryBytes(value) != s {
			return 0, fmt.Errorf("non-canonical byte form: %w", ErrInvalid)
		}
		return value, nil
	}
	if format(value) != s {
		return 0, fmt.Errorf("non-canonical byte form: %w", ErrInvalid)
	}

	return value, nil
}

func byteUnit(unit string) (int64, func(int64) string, bool) {
	decimalScales := decimalByteScales()
	if scale, ok := decimalScales.find(unit); ok {
		return scale.size, Bytes, true
	}
	binaryScales := binaryByteScales()
	if scale, ok := binaryScales.find(unit); ok {
		return scale.size, BinaryBytes, true
	}
	return 0, nil, false
}
