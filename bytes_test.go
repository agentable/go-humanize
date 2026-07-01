package humanize

import (
	"math"
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

		// Values that round to an integer at one decimal place must drop
		// the trailing zero.
		{"1024 rounds to 1 kB", 1024, "1 kB"},
		{"1049 rounds to 1 kB", 1049, "1 kB"},
		{"999.5 kB promotes to 1 MB", 999500, "1 MB"},

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

		// Negative values
		{"-1 byte", -1, "-1 B"},
		{"-1 kB", -1000, "-1 kB"},
		{"-1.5 MB", -1500000, "-1.5 MB"},
		{"-5.5 GB", int64(-5.5 * float64(gigabyte)), "-5.5 GB"},

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

		// Round-trip: values that round to an integer at one decimal
		// place must drop the trailing zero.
		{"1075 rounds to 1 KiB", 1075, "1 KiB"},
		{"1023.5 KiB promotes to 1 MiB", 1024*1024 - 512, "1 MiB"},

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

		// Negative values
		{"-1 byte", -1, "-1 B"},
		{"-1 KiB", -1024, "-1 KiB"},
		{"-1.5 MiB", -1572864, "-1.5 MiB"},
		{"-5.5 GiB", int64(-5.5 * float64(gibibyte)), "-5.5 GiB"},

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
