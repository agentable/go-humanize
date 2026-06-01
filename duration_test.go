package humanize

import (
	"errors"
	"math"
	"testing"
	"time"
)

func TestDuration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   time.Duration
		want string
	}{
		// Zero
		{"zero", 0, "0s"},

		// Microseconds
		{"1 microsecond", time.Microsecond, "1µs"},
		{"500 microseconds", 500 * time.Microsecond, "500µs"},
		{"999 microseconds", 999 * time.Microsecond, "999µs"},

		// Milliseconds
		{"1 millisecond", time.Millisecond, "1ms"},
		{"500 milliseconds", 500 * time.Millisecond, "500ms"},
		{"1ms 500µs", 1500 * time.Microsecond, "1ms 500µs"},
		{"999 milliseconds", 999 * time.Millisecond, "999ms"},

		// Seconds
		{"1 second", time.Second, "1s"},
		{"30 seconds", 30 * time.Second, "30s"},
		{"1s 500ms", 1500 * time.Millisecond, "1s 500ms"},
		{"59 seconds", 59 * time.Second, "59s"},

		// Minutes
		{"1 minute", time.Minute, "1m"},
		{"30 minutes", 30 * time.Minute, "30m"},
		{"1m 30s", 90 * time.Second, "1m 30s"},
		{"59 minutes", 59 * time.Minute, "59m"},

		// Hours
		{"1 hour", time.Hour, "1h"},
		{"1h 30m", 90 * time.Minute, "1h 30m"},
		{"2h 15m", 2*time.Hour + 15*time.Minute, "2h 15m"},
		{"23 hours", 23 * time.Hour, "23h"},

		// Days
		{"1 day", 24 * time.Hour, "1d"},
		{"2 days", 48 * time.Hour, "2d"},
		{"1d 5h", 29 * time.Hour, "1d 5h"},
		{"2d 12h", 60 * time.Hour, "2d 12h"},
		{"7 days", 7 * 24 * time.Hour, "7d"},

		// Mixed units (two significant units max)
		{"1h 30m 45s shows only 1h 30m", time.Hour + 30*time.Minute + 45*time.Second, "1h 30m"},
		{"2d 5h 30m shows only 2d 5h", 2*24*time.Hour + 5*time.Hour + 30*time.Minute, "2d 5h"},
		{"1m 30s 500ms shows only 1m 30s", time.Minute + 30*time.Second + 500*time.Millisecond, "1m 30s"},
		{"1h 1s skips zero minute unit", time.Hour + time.Second, "1h 1s"},
		{"1d 1µs skips empty lower units", 24*time.Hour + time.Microsecond, "1d 1µs"},

		// Negative durations
		{"-1 second", -time.Second, "-1s"},
		{"-1 minute", -time.Minute, "-1m"},
		{"-1h 30m", -90 * time.Minute, "-1h 30m"},
		{"-2d 5h", -53 * time.Hour, "-2d 5h"},
		{"-500ms", -500 * time.Millisecond, "-500ms"},

		// Edge cases
		{"sub-microsecond rounds to 0s", time.Nanosecond, "0s"},
		{"999 nanoseconds rounds to 0s", 999 * time.Nanosecond, "0s"},
		{"max duration", time.Duration(math.MaxInt64), "106751d 23h"},
		{"min duration", time.Duration(math.MinInt64), "-106751d 23h"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := Duration(tt.in)
			if got != tt.want {
				t.Errorf("Duration(%v) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestParseDuration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		in      string
		want    time.Duration
		wantErr bool
	}{
		// Valid formats
		{"1 second", "1s", time.Second, false},
		{"30 seconds", "30s", 30 * time.Second, false},
		{"1 minute", "1m", time.Minute, false},
		{"1m 30s", "1m 30s", 90 * time.Second, false},
		{"1 hour", "1h", time.Hour, false},
		{"1h 30m", "1h 30m", 90 * time.Minute, false},
		{"2h 15m", "2h 15m", 2*time.Hour + 15*time.Minute, false},
		{"sparse hour-second pair", "1h 1s", time.Hour + time.Second, false},

		// Days (extension over stdlib)
		{"1 day", "1d", 24 * time.Hour, false},
		{"2 days", "2d", 48 * time.Hour, false},
		{"1d 5h", "1d 5h", 29 * time.Hour, false},
		{"2d 12h", "2d 12h", 60 * time.Hour, false},
		{"7 days", "7d", 7 * 24 * time.Hour, false},
		{"sparse day-microsecond pair", "1d 1µs", 24*time.Hour + time.Microsecond, false},
		{"largest canonical day-hour pair", "106751d 23h", 106751*24*time.Hour + 23*time.Hour, false},
		{"largest canonical negative day-hour pair", "-106751d 23h", -(106751*24*time.Hour + 23*time.Hour), false},

		// Milliseconds and microseconds
		{"500ms", "500ms", 500 * time.Millisecond, false},
		{"1s 500ms", "1s 500ms", 1500 * time.Millisecond, false},
		{"1µs", "1µs", time.Microsecond, false},
		{"500µs", "500µs", 500 * time.Microsecond, false},
		{"1ms 500µs", "1ms 500µs", 1500 * time.Microsecond, false},

		// Negative durations
		{"-1s", "-1s", -time.Second, false},
		{"-1m", "-1m", -time.Minute, false},
		{"-1h 30m", "-1h 30m", -90 * time.Minute, false},
		{"-2d 5h", "-2d 5h", -53 * time.Hour, false},

		// Zero
		{"0s", "0s", 0, false},

		// Error cases
		{"empty string", "", 0, true},
		{"whitespace only", "   ", 0, true},
		{"invalid format", "invalid", 0, true},
		{"missing value", "h", 0, true},
		{"no unit", "123", 0, true},
		{"invalid unit", "1x", 0, true},
		{"leading space", " 1h", 0, true},
		{"trailing space", "1h ", 0, true},
		{"no space between parts", "1h30m", 0, true},
		{"extra spaces", "1h  30m", 0, true},
		{"tab separator", "1h\t30m", 0, true},
		{"plus sign", "+1s", 0, true},
		{"double minus", "--1s", 0, true},
		{"fractional hour", "1.5h", 0, true},
		{"fractional day", "2.5d", 0, true},
		{"us suffix", "500us", 0, true},
		{"leading zero hour", "01h", 0, true},
		{"negative zero", "-0s", 0, true},
		{"zero in wrong unit", "0m", 0, true},
		{"zero mixed with value", "0h 30m", 0, true},
		{"zero second unit", "1d 0h", 0, true},
		{"non-canonical seconds", "60s", 0, true},
		{"non-canonical second unit", "1h 60m", 0, true},
		{"non-canonical minutes", "60m", 0, true},
		{"non-canonical hours", "24h", 0, true},
		{"non-canonical milliseconds", "1000ms", 0, true},
		{"non-canonical microseconds", "1000µs", 0, true},
		{"overflow day count", "106752d", 0, true},
		{"overflow combined duration", "106751d 24h", 0, true},
		{"overflow int64 unit count", "9223372036854775808µs", 0, true},
		{"ascending units", "30m 1h", 0, true},
		{"duplicate units", "1h 1h", 0, true},
		{"too many parts", "1h 30m 5s", 0, true},
		{"invalid day format", "1d2d", 0, true},
		{"non-numeric", "abch", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := ParseDuration(tt.in)
			if tt.wantErr {
				if !errors.Is(err, ErrInvalid) {
					t.Errorf("ParseDuration(%q) error = %v, want ErrInvalid", tt.in, err)
				}
				return
			}
			if err != nil {
				t.Errorf("ParseDuration(%q) unexpected error: %v", tt.in, err)
				return
			}
			if got != tt.want {
				t.Errorf("ParseDuration(%q) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

func TestDurationDisplaySyntaxParses(t *testing.T) {
	t.Parallel()

	tests := []time.Duration{
		time.Second,
		time.Minute,
		time.Hour,
		24 * time.Hour,
		90 * time.Minute,
		2*time.Hour + 15*time.Minute,
		29 * time.Hour,
		500 * time.Millisecond,
		1500 * time.Microsecond,
	}

	for _, d := range tests {
		t.Run(d.String(), func(t *testing.T) {
			t.Parallel()

			formatted := Duration(d)
			parsed, err := ParseDuration(formatted)
			if err != nil {
				t.Errorf("ParseDuration(%q) error: %v", formatted, err)
				return
			}
			if parsed == 0 && d != 0 {
				t.Errorf("Duration(%v) = %q, ParseDuration = %v (lost all precision)", d, formatted, parsed)
			}
		})
	}
}

func FuzzParseDurationRejectsNonCanonical(f *testing.F) {
	seeds := []string{
		"0s",
		"1s",
		"1m",
		"1h 30m",
		"2d 5h",
		"500ms",
		"1ms 500µs",
		"-1h 30m",
		"",
		"60s",
		"24h",
		"1000ms",
		"1000µs",
		"01h",
		"-0s",
		"500us",
		"1h 30m 5s",
	}
	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, s string) {
		got, err := ParseDuration(s)
		if err != nil {
			if !errors.Is(err, ErrInvalid) {
				t.Errorf("ParseDuration(%q) error = %v, want ErrInvalid", s, err)
			}
			return
		}
		if formatted := Duration(got); formatted != s {
			t.Errorf("ParseDuration(%q) = %v, accepted non-canonical duration form %q", s, got, formatted)
		}
	})
}

func FuzzParseDurationRoundTripCanonical(f *testing.F) {
	seeds := []int64{
		0,
		int64(time.Nanosecond),
		int64(time.Microsecond),
		int64(time.Millisecond),
		int64(time.Second),
		int64(time.Minute),
		int64(time.Hour),
		int64(24 * time.Hour),
		int64(90 * time.Minute),
		math.MaxInt64,
		math.MinInt64,
	}
	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, nanos int64) {
		value := time.Duration(nanos)
		formatted := Duration(value)
		parsed, err := ParseDuration(formatted)
		if err != nil {
			t.Errorf("ParseDuration(%q) unexpected error: %v", formatted, err)
			return
		}
		if reformatted := Duration(parsed); reformatted != formatted {
			t.Errorf("round-trip mismatch: %q -> %v -> %q", formatted, parsed, reformatted)
		}
	})
}

func BenchmarkDuration(b *testing.B) {
	durations := []time.Duration{
		time.Second,
		time.Minute,
		time.Hour,
		24 * time.Hour,
		90 * time.Minute,
		2*time.Hour + 15*time.Minute,
		500 * time.Millisecond,
	}

	for b.Loop() {
		for _, d := range durations {
			_ = Duration(d)
		}
	}
}

func BenchmarkParseDuration(b *testing.B) {
	inputs := []string{
		"1s",
		"1m",
		"1h",
		"1d",
		"1h 30m",
		"2d 5h",
		"500ms",
	}

	for b.Loop() {
		for _, s := range inputs {
			if _, err := ParseDuration(s); err != nil {
				b.Fatalf("ParseDuration(%q) error = %v", s, err)
			}
		}
	}
}
