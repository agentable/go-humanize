package humanize

import (
	"math"
	"testing"
)

func TestCount(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		count    int64
		singular string
		plural   string
		want     string
	}{
		{name: "singular positive", count: 1, singular: "item", plural: "items", want: "1 item"},
		{name: "singular negative", count: -1, singular: "item", plural: "items", want: "-1 item"},
		{name: "irregular singular positive", count: 1, singular: "person", plural: "people", want: "1 person"},
		{name: "irregular singular negative", count: -1, singular: "person", plural: "people", want: "-1 person"},
		{name: "plural zero", count: 0, singular: "item", plural: "items", want: "0 items"},
		{name: "plural positive", count: 2, singular: "item", plural: "items", want: "2 items"},
		{name: "plural larger positive", count: 5, singular: "item", plural: "items", want: "5 items"},
		{name: "plural negative", count: -2, singular: "item", plural: "items", want: "-2 items"},
		{name: "plural larger negative", count: -5, singular: "item", plural: "items", want: "-5 items"},
		{name: "plural three digits", count: 100, singular: "item", plural: "items", want: "100 items"},
		{name: "plural negative three digits", count: -100, singular: "item", plural: "items", want: "-100 items"},
		{name: "empty singular", count: 1, singular: "", plural: "items", want: "1"},
		{name: "empty singular negative", count: -1, singular: "", plural: "items", want: "-1"},
		{name: "empty singular zero", count: 0, singular: "", plural: "items", want: "0 items"},
		{name: "empty plural singular case", count: 1, singular: "item", plural: "", want: "1 item"},
		{name: "empty plural zero case", count: 0, singular: "item", plural: "", want: "0"},
		{name: "empty plural general case", count: 2, singular: "item", plural: "", want: "2"},
		{name: "both forms empty negative singular", count: -1, singular: "", plural: "", want: "-1"},
		{name: "both forms empty zero", count: 0, singular: "", plural: "", want: "0"},
		{name: "both forms empty singular", count: 1, singular: "", plural: "", want: "1"},
		{name: "both forms empty plural", count: 2, singular: "", plural: "", want: "2"},
		{name: "irregular plural child", count: 2, singular: "child", plural: "children", want: "2 children"},
		{name: "irregular plural goose", count: 3, singular: "goose", plural: "geese", want: "3 geese"},
		{name: "thousands get commas", count: 1000, singular: "item", plural: "items", want: "1,000 items"},
		{name: "negative thousands get commas", count: -1000, singular: "item", plural: "items", want: "-1,000 items"},
		{name: "max int64", count: math.MaxInt64, singular: "item", plural: "items", want: "9,223,372,036,854,775,807 items"},
		{name: "min int64", count: math.MinInt64, singular: "item", plural: "items", want: "-9,223,372,036,854,775,808 items"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := Count(tt.count, tt.singular, tt.plural)
			if got != tt.want {
				t.Errorf("Count(%d, %q, %q) = %q, want %q",
					tt.count, tt.singular, tt.plural, got, tt.want)
			}
		})
	}
}

func BenchmarkCount(b *testing.B) {
	for b.Loop() {
		Count(5, "item", "items")
	}
}
