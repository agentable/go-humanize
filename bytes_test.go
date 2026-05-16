package humanize

import (
	"errors"
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

		// KB range
		{"1 KB", 1000, "1 KB"},
		{"1.5 KB", 1500, "1.5 KB"},
		{"10 KB", 10000, "10 KB"},
		{"999 KB", 999000, "999 KB"},

		// Round-trip: values that round to an integer at one decimal
		// place must drop the trailing zero so ParseBytes can round-trip
		// the formatter's own output.
		{"1024 rounds to 1 KB", 1024, "1 KB"},
		{"1049 rounds to 1 KB", 1049, "1 KB"},
		{"999.5 KB promotes to 1 MB", 999500, "1 MB"},

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
		{"-1 KB", -1000, "-1 KB"},
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

func TestParseBytes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		in      string
		want    int64
		wantErr bool
	}{
		// Canonical byte units
		{"bytes", "42 B", 42, false},
		{"zero bytes", "0 B", 0, false},
		{"upper byte range", "999 B", 999, false},
		{"binary byte decimal boundary", "1000 B", 1000, false},
		{"binary byte upper range", "1023 B", 1023, false},

		// Decimal units
		{"KB", "1 KB", 1000, false},
		{"MB", "42 MB", 42000000, false},
		{"MB decimal", "1.5 MB", 1500000, false},
		{"GB", "2 GB", 2000000000, false},
		{"TB", "1 TB", terabyte, false},
		{"PB", "1 PB", petabyte, false},
		{"EB", "1 EB", exabyte, false},

		// Binary units (IEC)
		{"KiB", "1 KiB", 1024, false},
		{"MiB", "42 MiB", 44040192, false},
		{"MiB decimal", "1.5 MiB", 1572864, false},
		{"GiB", "2 GiB", 2147483648, false},
		{"TiB", "1 TiB", tebibyte, false},
		{"PiB", "1 PiB", pebibyte, false},
		{"EiB", "1 EiB", exbibyte, false},

		// Negative values
		{"negative KB", "-1 KB", -1000, false},
		{"negative MiB", "-1.5 MiB", -1572864, false},

		// Error cases
		{"empty string", "", 0, true},
		{"missing unit", "42", 0, true},
		{"missing space", "42B", 0, true},
		{"missing space before binary unit", "1KiB", 0, true},
		{"leading space", " 1 KB", 0, true},
		{"trailing space", "1 KB ", 0, true},
		{"lowercase unit", "1 kib", 0, true},
		{"short unit", "42 K", 0, true},
		{"extra spaces", "42  MB", 0, true},
		{"comma in number", "1,005.03 MB", 0, true},
		{"leading zero", "01 KB", 0, true},
		{"negative zero", "-0 B", 0, true},
		{"zero non-byte unit", "0 KB", 0, true},
		{"fractional bytes", "1.5 B", 0, true},
		{"sub-unit decimal", "0.5 KB", 0, true},
		{"trailing zero decimal", "1.0 KB", 0, true},
		{"too many decimals", "1.55 MB", 0, true},
		{"decimal too large", "10.5 MB", 0, true},
		{"byte overflow form", "9223372036854775808 B", 0, true},
		{"non-canonical byte unit", "1024 B", 0, true},
		{"scaled positive overflow", "9.3 EB", 0, true},
		{"scaled negative overflow", "-9.3 EB", 0, true},
		{"unit overflow form", "1000 KB", 0, true},
		{"binary unit overflow form", "1024 KiB", 0, true},
		{"invalid unit", "42 XB", 0, true},
		{"invalid number", "abc KB", 0, true},
		{"not a number", "NaN KB", 0, true},
		{"leading plus", "+1 KB", 0, true},
		{"positive infinity", "+Inf KB", 0, true},
		{"negative infinity", "-Inf KB", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := ParseBytes(tt.in)
			if tt.wantErr {
				if !errors.Is(err, ErrInvalid) {
					t.Errorf("ParseBytes(%q) error = %v, want ErrInvalid", tt.in, err)
				}
				return
			}
			if err != nil {
				t.Errorf("ParseBytes(%q) unexpected error: %v", tt.in, err)
				return
			}
			if got != tt.want {
				t.Errorf("ParseBytes(%q) = %d, want %d", tt.in, got, tt.want)
			}
		})
	}
}

func TestParseBytesAcceptsFormatterBoundaryOutput(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		in     int64
		format func(int64) string
	}{
		{name: "binary max int64", in: math.MaxInt64, format: BinaryBytes},
		{name: "binary min int64", in: math.MinInt64, format: BinaryBytes},
		{name: "decimal max int64", in: math.MaxInt64, format: Bytes},
		{name: "decimal min int64", in: math.MinInt64, format: Bytes},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			formatted := tt.format(tt.in)
			got, err := ParseBytes(formatted)
			if err != nil {
				t.Errorf("ParseBytes(%q) unexpected error: %v", formatted, err)
				return
			}
			if reformatted := tt.format(got); reformatted != formatted {
				t.Errorf("ParseBytes(%q) = %d, reformatted as %q", formatted, got, reformatted)
			}
		})
	}
}

func TestBytesDisplaySyntaxParses(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		format func(int64) string
		values []int64
	}{
		{"Bytes", Bytes, []int64{
			1,
			1000,
			1024, // cross-base: previously formatted as "1.0 KB" and broke round-trip
			1049, // rounds to "1 KB" at one decimal
			1500,
			1000 * 1000,
			1500 * 1000,
			1000 * 1000 * 1000,
			terabyte,
			petabyte,
		}},
		{"BinaryBytes", BinaryBytes, []int64{
			1,
			1024,
			1075, // cross-base: 1075/1024 rounds to "1 KiB" at one decimal
			1536,
			1024 * 1024,
			1024 * 1024 * 1024,
			tebibyte,
			pebibyte,
		}},
	}

	for _, tc := range cases {
		for _, val := range tc.values {
			t.Run(tc.name+"/"+Number(val), func(t *testing.T) {
				t.Parallel()

				formatted := tc.format(val)
				parsed, err := ParseBytes(formatted)
				if err != nil {
					t.Errorf("ParseBytes(%q) unexpected error: %v", formatted, err)
					return
				}
				if reformatted := tc.format(parsed); reformatted != formatted {
					t.Errorf("round-trip mismatch: %q -> %d -> %q", formatted, parsed, reformatted)
				}
			})
		}
	}
}

func FuzzParseBytesRejectsNonCanonical(f *testing.F) {
	seeds := []string{
		"0 B",
		"1 KB",
		"1.5 MB",
		"1 KiB",
		"1.5 MiB",
		"-1 KB",
		"9.2 EB",
		"8 EiB",
		"",
		"1KB",
		"1.0 KB",
		"1000 KB",
		"1024 KiB",
		"0 KB",
		"1.5 B",
		"NaN KB",
		"+Inf KB",
	}
	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, s string) {
		got, err := ParseBytes(s)
		if err != nil {
			if !errors.Is(err, ErrInvalid) {
				t.Errorf("ParseBytes(%q) error = %v, want ErrInvalid", s, err)
			}
			return
		}
		if Bytes(got) != s && BinaryBytes(got) != s {
			t.Errorf("ParseBytes(%q) = %d, accepted non-canonical byte form", s, got)
		}
	})
}

func FuzzParseBytesRoundTripCanonical(f *testing.F) {
	seeds := []int64{
		0,
		1,
		999,
		1000,
		1024,
		1500,
		1000 * 1000,
		1024 * 1024,
		math.MaxInt64,
		math.MinInt64,
	}
	for _, seed := range seeds {
		f.Add(seed, false)
		f.Add(seed, true)
	}

	f.Fuzz(func(t *testing.T, value int64, binary bool) {
		format := Bytes
		if binary {
			format = BinaryBytes
		}

		formatted := format(value)
		parsed, err := ParseBytes(formatted)
		if err != nil {
			t.Errorf("ParseBytes(%q) unexpected error: %v", formatted, err)
			return
		}
		if reformatted := format(parsed); reformatted != formatted {
			t.Errorf("round-trip mismatch: %q -> %d -> %q", formatted, parsed, reformatted)
		}
	})
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

func BenchmarkParseBytes(b *testing.B) {
	for b.Loop() {
		_, _ = ParseBytes("2.5 GB")
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
