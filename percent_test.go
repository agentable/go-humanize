package humanize

import (
	"math"
	"testing"
)

func TestPercent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   float64
		want string
	}{
		// Whole percent points expressed as fractions
		{"zero", 0, "0%"},
		{"half", 0.5, "50%"},
		{"three quarters", 0.75, "75%"},
		{"one", 1, "100%"},
		{"over one is not clamped", 1.25, "125%"},
		{"negative whole", -0.25, "-25%"},

		// One-decimal fractions
		{"one decimal", 0.333, "33.3%"},
		{"another decimal", 0.667, "66.7%"},
		{"small fraction", 0.005, "0.5%"},
		{"negative decimal", -0.125, "-12.5%"},
		{"100/3", (100.0 / 3) / 100, "33.3%"},
		{"200/3", (200.0 / 3) / 100, "66.7%"},

		// Rounding
		{"rounds to one decimal", 0.3333333, "33.3%"},
		{"rounds up", 0.6666666, "66.7%"},
		{"rounds half up", 0.123, "12.3%"},
		{"rounds positive half-tenth up", 0.0005, "0.1%"},
		{"rounds negative half-tenth down", -0.0005, "-0.1%"},
		{"rounds near one up", 0.9995, "100%"},
		{"rounds negative near one down", -0.9995, "-100%"},

		// Sub-decimal collapses to 0%
		{"very small rounds to zero", 0.0001, "0%"},
		{"very small negative rounds to zero", -0.0001, "0%"},

		// Large values
		{"very large", 9999.999, "999999.9%"},
		{"max finite stays finite", math.MaxFloat64, "1.8e+310%"},
		{"negative max finite stays finite", -math.MaxFloat64, "-1.8e+310%"},

		// Special float values
		{"NaN", math.NaN(), "NaN%"},
		{"+Inf", math.Inf(1), "+Inf%"},
		{"-Inf", math.Inf(-1), "-Inf%"},

		// Negative zero
		{"negative zero", math.Copysign(0, -1), "0%"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := Percent(tt.in)
			if got != tt.want {
				t.Errorf("Percent(%f) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestPercentFiniteLargeBoundary(t *testing.T) {
	t.Parallel()

	tests := []float64{
		math.MaxFloat64 / 1000,
		math.Nextafter(math.MaxFloat64/1000, math.Inf(1)),
	}

	for _, in := range tests {
		got := Percent(in)
		if got == "+Inf%" || got == "-Inf%" {
			t.Errorf("Percent(%g) = %q, want finite output", in, got)
		}
	}
}

func BenchmarkPercent(b *testing.B) {
	for b.Loop() {
		Percent(0.333)
	}
}

func BenchmarkPercentWhole(b *testing.B) {
	for b.Loop() {
		Percent(0.75)
	}
}

func BenchmarkPercentSpecial(b *testing.B) {
	for b.Loop() {
		Percent(math.NaN())
	}
}
