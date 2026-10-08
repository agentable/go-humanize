package humanize

import "strconv"

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
	size uint64
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
// It preserves the sign and rounds the exact integer magnitude, with ties
// away from zero.
func Bytes(b int64) string {
	scales := decimalByteScales()
	return formatBytes(b, &scales)
}

// BinaryBytes formats b with IEC binary byte units such as "1.5 KiB" and
// "2 MiB". It preserves the sign and rounds the exact integer magnitude,
// with ties away from zero.
func BinaryBytes(b int64) string {
	scales := binaryByteScales()
	return formatBytes(b, &scales)
}

func formatBytes(b int64, scales *byteScales) string {
	var magnitude uint64
	sign := ""
	if b < 0 {
		// Avoid signed overflow when b is MinInt64.
		magnitude = uint64(-(b + 1)) + 1 //nolint:gosec // For b < 0, -(b+1) is in [0, MaxInt64].
		sign = "-"
	} else {
		magnitude = uint64(b)
	}

	if magnitude < scales[1].size {
		return sign + strconv.FormatUint(magnitude, 10) + " B"
	}

	exp := 0
	for exp+1 < len(scales) && magnitude >= scales[exp+1].size {
		exp++
	}

	size := scales[exp].size
	value, remainder := magnitude/size, magnitude%size
	tenths := value < 10
	if tenths {
		// The largest size is 2^60, so remainder*10 fits in uint64.
		value = value*10 + remainder*10/size
		remainder = remainder * 10 % size
	}
	if remainder >= (size+1)/2 {
		value++
	}

	if !tenths && exp+1 < len(scales) && value >= scales[1].size {
		return sign + "1 " + scales[exp+1].unit
	}
	number := strconv.FormatUint(value, 10)
	if tenths {
		number = strconv.FormatUint(value/10, 10)
		if value%10 != 0 {
			number += "." + strconv.FormatUint(value%10, 10)
		}
	}
	return sign + number + " " + scales[exp].unit
}
