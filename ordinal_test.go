package humanize

import (
	"strconv"
	"testing"
)

func TestOrdinal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   int64
		want string
	}{
		// Zero and small values
		{"zero", 0, "0th"},
		{"one", 1, "1st"},
		{"two", 2, "2nd"},
		{"three", 3, "3rd"},
		{"four", 4, "4th"},
		{"five", 5, "5th"},
		{"six", 6, "6th"},
		{"seven", 7, "7th"},
		{"eight", 8, "8th"},
		{"nine", 9, "9th"},
		{"ten", 10, "10th"},

		// Special cases: 11th, 12th, 13th
		{"eleven", 11, "11th"},
		{"twelve", 12, "12th"},
		{"thirteen", 13, "13th"},
		{"fourteen", 14, "14th"},

		// Twenties
		{"twenty", 20, "20th"},
		{"twenty-first", 21, "21st"},
		{"twenty-second", 22, "22nd"},
		{"twenty-third", 23, "23rd"},
		{"twenty-fourth", 24, "24th"},

		// Thirties
		{"thirty-first", 31, "31st"},
		{"thirty-second", 32, "32nd"},
		{"thirty-third", 33, "33rd"},
		{"thirty-fourth", 34, "34th"},

		// Forties
		{"forty-first", 41, "41st"},
		{"forty-second", 42, "42nd"},
		{"forty-third", 43, "43rd"},

		// Hundreds with special cases
		{"one hundred", 100, "100th"},
		{"one hundred first", 101, "101st"},
		{"one hundred second", 102, "102nd"},
		{"one hundred third", 103, "103rd"},
		{"one hundred fourth", 104, "104th"},
		{"one hundred tenth", 110, "110th"},
		{"one hundred eleventh", 111, "111th"},
		{"one hundred twelfth", 112, "112th"},
		{"one hundred thirteenth", 113, "113th"},
		{"one hundred twenty-first", 121, "121st"},

		// Two hundreds with special cases
		{"two hundred eleventh", 211, "211th"},
		{"two hundred twelfth", 212, "212th"},
		{"two hundred thirteenth", 213, "213th"},
		{"two hundred twenty-first", 221, "221st"},

		// Thousands
		{"one thousand", 1000, "1000th"},
		{"one thousand first", 1001, "1001st"},
		{"one thousand eleventh", 1011, "1011th"},
		{"one thousand twenty-first", 1021, "1021st"},

		// Large numbers
		{"ten thousand", 10000, "10000th"},
		{"one hundred thousand", 100000, "100000th"},
		{"one million", 1000000, "1000000th"},

		// Negative values
		{"negative one", -1, "-1st"},
		{"negative two", -2, "-2nd"},
		{"negative three", -3, "-3rd"},
		{"negative four", -4, "-4th"},
		{"negative eleven", -11, "-11th"},
		{"negative twelve", -12, "-12th"},
		{"negative thirteen", -13, "-13th"},
		{"negative twenty-first", -21, "-21st"},
		{"negative one hundred", -100, "-100th"},
		{"negative one hundred eleventh", -111, "-111th"},
		{"min int64", int64(-1 << 63), strconv.FormatInt(int64(-1<<63), 10) + "th"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := Ordinal(tt.in)
			if got != tt.want {
				t.Errorf("Ordinal(%d) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func BenchmarkOrdinal(b *testing.B) {
	for b.Loop() {
		Ordinal(21)
	}
}

func BenchmarkOrdinalSpecialCase(b *testing.B) {
	for b.Loop() {
		Ordinal(11)
	}
}

func BenchmarkOrdinalLarge(b *testing.B) {
	for b.Loop() {
		Ordinal(1000000)
	}
}

func BenchmarkOrdinalNegative(b *testing.B) {
	for b.Loop() {
		Ordinal(-21)
	}
}
