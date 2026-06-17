package humanize

import (
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
