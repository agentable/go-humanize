package humanize

import (
	"math"
	"strconv"
	"testing"
)

func TestBytes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   int64
		want string
	}{
		// Zero and small values
		{"zero", 0, "0 B"},
		{"one byte", 1, "1 B"},
		{"nine bytes", 9, "9 B"},

		// kB range
		{"1 kB", 1000, "1 kB"},
		{"1.5 kB", 1500, "1.5 kB"},
		{"10 kB", 10000, "10 kB"},
		{"999 kB", 999000, "999 kB"},
		{"below half tenth", 1249, "1.2 kB"},
		{"half tenth", 1250, "1.3 kB"},
		{"above half tenth", 1251, "1.3 kB"},
		{"below second half tenth", 2249, "2.2 kB"},
		{"second half tenth", 2250, "2.3 kB"},
		{"above second half tenth", 2251, "2.3 kB"},
		{"below whole tie", 10499, "10 kB"},
		{"whole tie", 10500, "11 kB"},
		{"above whole tie", 10501, "11 kB"},

		// Values that round to an integer at one decimal place must drop
		// the trailing zero.
		{"1024 rounds to 1 kB", 1024, "1 kB"},
		{"1049 rounds to 1 kB", 1049, "1 kB"},
		{"999.5 kB promotes to 1 MB", 999500, "1 MB"},
		{"below MB promotion", 999499, "999 kB"},
		{"below GB promotion", 999499000, "999 MB"},
		{"GB promotion", 999500000, "1 GB"},

		// MB range
		{"1 MB", 1000000, "1 MB"},
		{"1.5 MB", 1500000, "1.5 MB"},
		{"10 MB", 10000000, "10 MB"},
		{"999 MB", 999000000, "999 MB"},

		// GB range
		{"1 GB", 1000000000, "1 GB"},
		{"5.5 GB", int64(5.5 * float64(gigabyte)), "5.5 GB"},
		{"10 GB", 10000000000, "10 GB"},

		// TB range
		{"1 TB", terabyte, "1 TB"},
		{"2.5 TB", int64(2.5 * float64(terabyte)), "2.5 TB"},

		// PB range
		{"1 PB", petabyte, "1 PB"},
		{"1.5 PB", int64(1.5 * float64(petabyte)), "1.5 PB"},

		// EB range
		{"1 EB", exabyte, "1 EB"},
		{"below EB half tenth", 1249999999999999999, "1.2 EB"},
		{"EB half tenth", 1250000000000000000, "1.3 EB"},
		{"below EB promotion", 999499999999999999, "999 PB"},
		{"EB promotion", 999500000000000000, "1 EB"},
		{"negative below EB half tenth", -1249999999999999999, "-1.2 EB"},
		{"negative below EB promotion", -999499999999999999, "-999 PB"},

		// Negative values
		{"-1 byte", -1, "-1 B"},
		{"-1 kB", -1000, "-1 kB"},
		{"-1.5 MB", -1500000, "-1.5 MB"},
		{"-5.5 GB", int64(-5.5 * float64(gigabyte)), "-5.5 GB"},
		{"negative below half tenth", -1249, "-1.2 kB"},
		{"negative half tenth", -1250, "-1.3 kB"},
		{"negative above half tenth", -1251, "-1.3 kB"},
		{"negative below second half tenth", -2249, "-2.2 kB"},
		{"negative second half tenth", -2250, "-2.3 kB"},
		{"negative above second half tenth", -2251, "-2.3 kB"},
		{"negative below whole tie", -10499, "-10 kB"},
		{"negative whole tie", -10500, "-11 kB"},
		{"negative above whole tie", -10501, "-11 kB"},
		{"negative below MB promotion", -999499, "-999 kB"},
		{"negative MB promotion", -999500, "-1 MB"},
		{"negative below GB promotion", -999499000, "-999 MB"},
		{"negative GB promotion", -999500000, "-1 GB"},

		// Edge cases
		{"math.MaxInt64", math.MaxInt64, "9.2 EB"},
		{"math.MinInt64", math.MinInt64, "-9.2 EB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := Bytes(tt.in)
			if got != tt.want {
				t.Errorf("Bytes(%d) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestBinaryBytes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   int64
		want string
	}{
		// Zero and small values
		{"zero", 0, "0 B"},
		{"one byte", 1, "1 B"},
		{"nine bytes", 9, "9 B"},

		// KiB range
		{"1 KiB", 1024, "1 KiB"},
		{"1.5 KiB", 1536, "1.5 KiB"},
		{"10 KiB", 10240, "10 KiB"},
		{"1023 KiB", 1024*1024 - 1024, "1023 KiB"},
		{"below half tenth", 1279, "1.2 KiB"},
		{"half tenth", 1280, "1.3 KiB"},
		{"above half tenth", 1281, "1.3 KiB"},

		// Round-trip: values that round to an integer at one decimal
		// place must drop the trailing zero.
		{"1075 rounds to 1 KiB", 1075, "1 KiB"},
		{"1023.5 KiB promotes to 1 MiB", 1024*1024 - 512, "1 MiB"},
		{"below MiB promotion", 1048063, "1023 KiB"},
		{"below GiB promotion", 1073217535, "1023 MiB"},
		{"GiB promotion", 1073217536, "1 GiB"},

		// MiB range
		{"1 MiB", 1024 * 1024, "1 MiB"},
		{"1.5 MiB", 1572864, "1.5 MiB"},
		{"10 MiB", 10485760, "10 MiB"},
		{"1023 MiB", 1024*1024*1024 - 1024*1024, "1023 MiB"},

		// GiB range
		{"1 GiB", 1024 * 1024 * 1024, "1 GiB"},
		{"5.5 GiB", int64(5.5 * float64(gibibyte)), "5.5 GiB"},
		{"10 GiB", 10737418240, "10 GiB"},

		// TiB range
		{"1 TiB", tebibyte, "1 TiB"},
		{"2.5 TiB", int64(2.5 * float64(tebibyte)), "2.5 TiB"},

		// PiB range
		{"1 PiB", pebibyte, "1 PiB"},
		{"1.5 PiB", int64(1.5 * float64(pebibyte)), "1.5 PiB"},

		// EiB range
		{"1 EiB", exbibyte, "1 EiB"},
		{"below EiB half tenth", 1441151880758558719, "1.2 EiB"},
		{"EiB half tenth", 1441151880758558720, "1.3 EiB"},
		{"below EiB promotion", 1152358554653425663, "1023 PiB"},
		{"EiB promotion", 1152358554653425664, "1 EiB"},
		{"negative below EiB half tenth", -1441151880758558719, "-1.2 EiB"},
		{"negative below EiB promotion", -1152358554653425663, "-1023 PiB"},

		// Negative values
		{"-1 byte", -1, "-1 B"},
		{"-1 KiB", -1024, "-1 KiB"},
		{"-1.5 MiB", -1572864, "-1.5 MiB"},
		{"-5.5 GiB", int64(-5.5 * float64(gibibyte)), "-5.5 GiB"},
		{"negative below half tenth", -1279, "-1.2 KiB"},
		{"negative half tenth", -1280, "-1.3 KiB"},
		{"negative above half tenth", -1281, "-1.3 KiB"},
		{"negative below MiB promotion", -1048063, "-1023 KiB"},
		{"negative MiB promotion", -1048064, "-1 MiB"},
		{"negative below GiB promotion", -1073217535, "-1023 MiB"},
		{"negative GiB promotion", -1073217536, "-1 GiB"},

		// Edge cases
		{"math.MaxInt64", math.MaxInt64, "8 EiB"},
		{"math.MinInt64", math.MinInt64, "-8 EiB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := BinaryBytes(tt.in)
			if got != tt.want {
				t.Errorf("BinaryBytes(%d) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestByteRoundingBoundaries(t *testing.T) {
	t.Parallel()

	type testCase struct {
		in   uint64
		want string
	}

	variants := []struct {
		name   string
		base   uint64
		units  []string
		format func(int64) string
	}{
		{"decimal", 1000, []string{"kB", "MB", "GB", "TB", "PB", "EB"}, Bytes},
		{"binary", 1024, []string{"KiB", "MiB", "GiB", "TiB", "PiB", "EiB"}, BinaryBytes},
	}
	for _, variant := range variants {
		t.Run(variant.name, func(t *testing.T) {
			t.Parallel()
			size := uint64(1)
			for i, unit := range variant.units {
				size *= variant.base
				// 1.25 units is an exact half tenth in both unit systems.
				half := size + size/4
				cases := []testCase{
					{half - 1, "1.2 " + unit},
					{half, "1.3 " + unit},
					{half + 1, "1.3 " + unit},
				}
				// Unscaled bytes do not round into the first unit.
				if i == 0 {
					cases = append(cases, testCase{size - 1, strconv.FormatUint(variant.base-1, 10) + " B"})
				}
				// 1.15 units is between integers for binary units. Check
				// its two neighboring bytes, as well as exact decimal ties.
				floorHalf := size + 3*size/20
				atFloor := "1.1 " + unit
				if 3*size%20 == 0 {
					atFloor = "1.2 " + unit
				}
				cases = append(cases,
					testCase{floorHalf - 1, "1.1 " + unit},
					testCase{floorHalf, atFloor},
					testCase{floorHalf + 1, "1.2 " + unit},
				)
				if size <= math.MaxInt64/11 {
					wholeHalf := 10*size + size/2
					carryFloor := 9*size + 19*size/20
					atCarryFloor := "9.9 " + unit
					if 19*size%20 == 0 {
						atCarryFloor = "10 " + unit
					}
					cases = append(cases, []testCase{
						{wholeHalf - 1, "10 " + unit},
						{wholeHalf, "11 " + unit},
						{wholeHalf + 1, "11 " + unit},
						{carryFloor - 1, "9.9 " + unit},
						{carryFloor, atCarryFloor},
						{carryFloor + 1, "10 " + unit},
					}...)
				}
				if i+1 < len(variant.units) {
					promotion := (variant.base-1)*size + size/2
					cases = append(cases, []testCase{
						{promotion - 1, strconv.FormatUint(variant.base-1, 10) + " " + unit},
						{promotion, "1 " + variant.units[i+1]},
						{promotion + 1, "1 " + variant.units[i+1]},
					}...)
				}
				for _, tt := range cases {
					for _, sign := range []int64{1, -1} {
						in, want := sign*int64(tt.in), tt.want
						if sign < 0 {
							want = "-" + want
						}
						if got := variant.format(in); got != want {
							t.Errorf("%s(%d) = %q, want %q", variant.name, in, got, want)
						}
					}
				}
			}
		})
	}
}

func BenchmarkBytes(b *testing.B) {
	for b.Loop() {
		Bytes(1500000)
	}
}

func BenchmarkBinaryBytes(b *testing.B) {
	for b.Loop() {
		BinaryBytes(1536)
	}
}

func BenchmarkBytesLarge(b *testing.B) {
	for b.Loop() {
		Bytes(5 * terabyte)
	}
}

func BenchmarkBytesNegative(b *testing.B) {
	for b.Loop() {
		Bytes(-1500000)
	}
}
