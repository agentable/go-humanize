package humanize

import (
	"math"
	"testing"
)

func TestNumber(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   int64
		want string
	}{
		// Zero and small values
		{"zero", 0, "0"},
		{"one", 1, "1"},
		{"nine", 9, "9"},
		{"ten", 10, "10"},
		{"ninety-nine", 99, "99"},
		{"one hundred", 100, "100"},
		{"nine hundred ninety-nine", 999, "999"},

		// Thousands
		{"one thousand", 1000, "1,000"},
		{"ten thousand", 10000, "10,000"},
		{"one hundred thousand", 100000, "100,000"},
		{"one million", 1000000, "1,000,000"},

		// Various patterns
		{"1,234,567", 1234567, "1,234,567"},
		{"10,100,000", 10100000, "10,100,000"},
		{"10,010,000", 10010000, "10,010,000"},
		{"10,001,000", 10001000, "10,001,000"},
		{"123,456,789", 123456789, "123,456,789"},

		// Negative values
		{"-1", -1, "-1"},
		{"-10", -10, "-10"},
		{"-100", -100, "-100"},
		{"-1,000", -1000, "-1,000"},
		{"-10,000", -10000, "-10,000"},
		{"-100,000", -100000, "-100,000"},
		{"-1,000,000", -1000000, "-1,000,000"},
		{"-10,000,000", -10000000, "-10,000,000"},
		{"-10,100,000", -10100000, "-10,100,000"},
		{"-10,010,000", -10010000, "-10,010,000"},
		{"-10,001,000", -10001000, "-10,001,000"},
		{"-123,456,789", -123456789, "-123,456,789"},

		// Edge cases
		{"math.MaxInt64", math.MaxInt64, "9,223,372,036,854,775,807"},
		{"math.MinInt64", math.MinInt64, "-9,223,372,036,854,775,808"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := Number(tt.in)
			if got != tt.want {
				t.Errorf("Number(%d) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func BenchmarkNumber(b *testing.B) {
	for b.Loop() {
		Number(1234567)
	}
}

func BenchmarkNumberLarge(b *testing.B) {
	for b.Loop() {
		Number(math.MaxInt64)
	}
}

func BenchmarkNumberSmall(b *testing.B) {
	for b.Loop() {
		Number(999)
	}
}

func BenchmarkNumberNegative(b *testing.B) {
	for b.Loop() {
		Number(-123456789)
	}
}
