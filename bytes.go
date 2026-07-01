package humanize

import (
	"math"
	"strconv"
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
		{unit: "kB", size: kilobyte},
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

// Bytes formats b with decimal byte units such as "1.5 kB" and "2 MB".
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
