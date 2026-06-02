package humanize

import (
	"math"
	"testing"
	"time"
)

func TestRelative(t *testing.T) {
	t.Parallel()

	now := time.Date(2024, 2, 24, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name string
		t    time.Time
		now  time.Time
		want string
	}{
		// Just now (< 1 second)
		{"just now - 0ms", now, now, "just now"},
		{"just now - 500ms ago", now.Add(-500 * time.Millisecond), now, "just now"},
		{"just now - 999ms ago", now.Add(-999 * time.Millisecond), now, "just now"},
		{"just now - 500ms future", now.Add(500 * time.Millisecond), now, "just now"},
		{"just now - 999ms future", now.Add(999 * time.Millisecond), now, "just now"},

		// Seconds
		{"1 second ago", now.Add(-time.Second), now, "1 second ago"},
		{"2 seconds ago", now.Add(-2 * time.Second), now, "2 seconds ago"},
		{"30 seconds ago", now.Add(-30 * time.Second), now, "30 seconds ago"},
		{"59 seconds ago", now.Add(-59 * time.Second), now, "59 seconds ago"},
		{"in 1 second", now.Add(time.Second), now, "in 1 second"},
		{"in 30 seconds", now.Add(30 * time.Second), now, "in 30 seconds"},

		// Minutes
		{"1 minute ago", now.Add(-time.Minute), now, "1 minute ago"},
		{"2 minutes ago", now.Add(-2 * time.Minute), now, "2 minutes ago"},
		{"30 minutes ago", now.Add(-30 * time.Minute), now, "30 minutes ago"},
		{"59 minutes ago", now.Add(-59 * time.Minute), now, "59 minutes ago"},
		{"in 1 minute", now.Add(time.Minute), now, "in 1 minute"},
		{"in 45 minutes", now.Add(45 * time.Minute), now, "in 45 minutes"},

		// Hours
		{"1 hour ago", now.Add(-time.Hour), now, "1 hour ago"},
		{"2 hours ago", now.Add(-2 * time.Hour), now, "2 hours ago"},
		{"12 hours ago", now.Add(-12 * time.Hour), now, "12 hours ago"},
		{"23 hours ago", now.Add(-23 * time.Hour), now, "23 hours ago"},
		{"in 1 hour", now.Add(time.Hour), now, "in 1 hour"},
		{"in 6 hours", now.Add(6 * time.Hour), now, "in 6 hours"},

		// Days
		{"1 day ago", now.Add(-24 * time.Hour), now, "1 day ago"},
		{"2 days ago", now.Add(-48 * time.Hour), now, "2 days ago"},
		{"3 days ago", now.Add(-3 * 24 * time.Hour), now, "3 days ago"},
		{"6 days ago", now.Add(-6 * 24 * time.Hour), now, "6 days ago"},
		{"in 1 day", now.Add(24 * time.Hour), now, "in 1 day"},
		{"in 5 days", now.Add(5 * 24 * time.Hour), now, "in 5 days"},

		// Weeks
		{"1 week ago", now.Add(-7 * 24 * time.Hour), now, "1 week ago"},
		{"2 weeks ago", now.Add(-14 * 24 * time.Hour), now, "2 weeks ago"},
		{"3 weeks ago", now.Add(-21 * 24 * time.Hour), now, "3 weeks ago"},
		{"4 weeks ago", now.Add(-28 * 24 * time.Hour), now, "4 weeks ago"},
		{"in 1 week", now.Add(7 * 24 * time.Hour), now, "in 1 week"},
		{"in 2 weeks", now.Add(14 * 24 * time.Hour), now, "in 2 weeks"},

		// Months (30 days)
		{"1 month ago", now.Add(-30 * 24 * time.Hour), now, "1 month ago"},
		{"2 months ago", now.Add(-60 * 24 * time.Hour), now, "2 months ago"},
		{"6 months ago", now.Add(-180 * 24 * time.Hour), now, "6 months ago"},
		{"11 months ago", now.Add(-330 * 24 * time.Hour), now, "11 months ago"},
		{"almost 1 year ago", now.Add(-350 * 24 * time.Hour), now, "1 year ago"},
		{"in 1 month", now.Add(30 * 24 * time.Hour), now, "in 1 month"},
		{"in 3 months", now.Add(90 * 24 * time.Hour), now, "in 3 months"},
		{"almost 1 year in future", now.Add(350 * 24 * time.Hour), now, "in 1 year"},

		// Years (365 days)
		{"1 year ago", now.Add(-365 * 24 * time.Hour), now, "1 year ago"},
		{"2 years ago", now.Add(-730 * 24 * time.Hour), now, "2 years ago"},
		{"5 years ago", now.Add(-5 * 365 * 24 * time.Hour), now, "5 years ago"},
		{"10 years ago", now.Add(-10 * 365 * 24 * time.Hour), now, "10 years ago"},
		{"in 1 year", now.Add(365 * 24 * time.Hour), now, "in 1 year"},
		{"in 3 years", now.Add(3 * 365 * 24 * time.Hour), now, "in 3 years"},

		// Edge cases
		{"negative zero", now, now, "just now"},
		{"very large past", now.Add(-100 * 365 * 24 * time.Hour), now, "100 years ago"},
		{"very large future", now.Add(100 * 365 * 24 * time.Hour), now, "in 100 years"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := Relative(tt.t, tt.now)
			if got != tt.want {
				t.Errorf("Relative(%v, %v) = %q, want %q", tt.t, tt.now, got, tt.want)
			}
		})
	}
}

func TestRelativeEdgeCases(t *testing.T) {
	t.Parallel()

	ref := time.Date(2024, 2, 24, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name   string
		target time.Time
		want   string
	}{
		{name: "zero time", target: time.Time{}, want: "292 years ago"},
		{name: "unix epoch", target: time.Unix(0, 0).UTC(), want: "54 years ago"},
		{name: "far future", target: time.Unix(253402300799, 0).UTC(), want: "in 292 years"}, // Year 9999
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := Relative(tt.target, ref)
			if got != tt.want {
				t.Errorf("Relative(%v, %v) = %q, want %q", tt.target, ref, got, tt.want)
			}
		})
	}
}

func TestRelativeBoundaries(t *testing.T) {
	t.Parallel()

	now := time.Date(2024, 2, 24, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name string
		t    time.Time
		want string
	}{
		// Boundary between seconds and minutes (60 seconds)
		{"59 seconds", now.Add(-59 * time.Second), "59 seconds ago"},
		{"60 seconds", now.Add(-60 * time.Second), "1 minute ago"},
		{"61 seconds", now.Add(-61 * time.Second), "1 minute ago"},

		// Boundary between minutes and hours (60 minutes)
		{"59 minutes", now.Add(-59 * time.Minute), "59 minutes ago"},
		{"60 minutes", now.Add(-60 * time.Minute), "1 hour ago"},
		{"61 minutes", now.Add(-61 * time.Minute), "1 hour ago"},

		// Boundary between hours and days (24 hours)
		{"23 hours", now.Add(-23 * time.Hour), "23 hours ago"},
		{"24 hours", now.Add(-24 * time.Hour), "1 day ago"},
		{"25 hours", now.Add(-25 * time.Hour), "1 day ago"},

		// Boundary between days and weeks (7 days)
		{"6 days", now.Add(-6 * 24 * time.Hour), "6 days ago"},
		{"7 days", now.Add(-7 * 24 * time.Hour), "1 week ago"},
		{"8 days", now.Add(-8 * 24 * time.Hour), "1 week ago"},

		// Boundary between weeks and months (30 days)
		{"29 days", now.Add(-29 * 24 * time.Hour), "4 weeks ago"},
		{"30 days", now.Add(-30 * 24 * time.Hour), "1 month ago"},
		{"31 days", now.Add(-31 * 24 * time.Hour), "1 month ago"},

		// Boundary between months and years (345-day cutover)
		{"344 days", now.Add(-344 * 24 * time.Hour), "11 months ago"},
		{"345 days", now.Add(-345 * 24 * time.Hour), "1 year ago"},
		{"364 days", now.Add(-364 * 24 * time.Hour), "1 year ago"},
		{"365 days", now.Add(-365 * 24 * time.Hour), "1 year ago"},
		{"366 days", now.Add(-366 * 24 * time.Hour), "1 year ago"},

		// Future cutovers mirror past cutovers.
		{"future 60 seconds", now.Add(60 * time.Second), "in 1 minute"},
		{"future 24 hours", now.Add(24 * time.Hour), "in 1 day"},
		{"future 7 days", now.Add(7 * 24 * time.Hour), "in 1 week"},
		{"future 30 days", now.Add(30 * 24 * time.Hour), "in 1 month"},
		{"future 345 days", now.Add(345 * 24 * time.Hour), "in 1 year"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := Relative(tt.t, now)
			if got != tt.want {
				t.Errorf("Relative(%v, %v) = %q, want %q", tt.t, now, got, tt.want)
			}
		})
	}
}

func TestRelativeSymmetry(t *testing.T) {
	t.Parallel()

	now := time.Date(2024, 2, 24, 12, 0, 0, 0, time.UTC)
	past := now.Add(-2 * time.Hour)
	future := now.Add(2 * time.Hour)

	gotPast := Relative(past, now)
	gotFuture := Relative(future, now)

	if gotPast != "2 hours ago" {
		t.Errorf("Relative(past, now) = %q, want '2 hours ago'", gotPast)
	}
	if gotFuture != "in 2 hours" {
		t.Errorf("Relative(future, now) = %q, want 'in 2 hours'", gotFuture)
	}
}

func TestRelativeSaturatedRange(t *testing.T) {
	t.Parallel()

	ancient := time.Time{}
	farFuture := time.Unix(253402300799, 0).UTC()

	if diff := ancient.Sub(farFuture); diff != time.Duration(math.MinInt64) {
		t.Fatalf("ancient.Sub(farFuture) = %v, want %v", diff, time.Duration(math.MinInt64))
	}
	if diff := farFuture.Sub(ancient); diff != time.Duration(math.MaxInt64) {
		t.Fatalf("farFuture.Sub(ancient) = %v, want %v", diff, time.Duration(math.MaxInt64))
	}
	if got := Relative(farFuture, ancient); got != "in 292 years" {
		t.Errorf("Relative(farFuture, ancient) = %q, want %q", got, "in 292 years")
	}
	if got := Relative(ancient, farFuture); got != "292 years ago" {
		t.Errorf("Relative(ancient, farFuture) = %q, want %q", got, "292 years ago")
	}
}

func BenchmarkRelative(b *testing.B) {
	now := time.Now()
	times := []time.Time{
		now.Add(-time.Second),
		now.Add(-time.Minute),
		now.Add(-time.Hour),
		now.Add(-24 * time.Hour),
		now.Add(-7 * 24 * time.Hour),
		now.Add(-30 * 24 * time.Hour),
		now.Add(-365 * 24 * time.Hour),
	}

	for b.Loop() {
		for _, t := range times {
			_ = Relative(t, now)
		}
	}
}
