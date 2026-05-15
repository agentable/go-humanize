package humanize_test

import (
	"errors"
	"fmt"
	"time"

	"github.com/agentable/go-humanize"
)

func ExampleBytes() {
	const gigabyte = 1000 * 1000 * 1000

	fmt.Println(humanize.Bytes(0))
	fmt.Println(humanize.Bytes(1500))
	fmt.Println(humanize.Bytes(1000000))
	fmt.Println(humanize.Bytes(5 * gigabyte))
	// Output:
	// 0 B
	// 1.5 KB
	// 1 MB
	// 5 GB
}

func ExampleBinaryBytes() {
	const gibibyte = 1024 * 1024 * 1024

	fmt.Println(humanize.BinaryBytes(0))
	fmt.Println(humanize.BinaryBytes(1536))
	fmt.Println(humanize.BinaryBytes(1048576))
	fmt.Println(humanize.BinaryBytes(5 * gibibyte))
	// Output:
	// 0 B
	// 1.5 KiB
	// 1 MiB
	// 5 GiB
}

func ExampleParseBytes() {
	kb, _ := humanize.ParseBytes("42 KB")
	fmt.Println(kb)

	mib, _ := humanize.ParseBytes("1.5 MiB")
	fmt.Println(mib)

	gb, _ := humanize.ParseBytes("2.5 GB")
	fmt.Println(gb)

	if _, err := humanize.ParseBytes("garbage"); errors.Is(err, humanize.ErrInvalid) {
		fmt.Println("invalid")
	}
	// Output:
	// 42000
	// 1572864
	// 2500000000
	// invalid
}

func ExampleNumber() {
	fmt.Println(humanize.Number(0))
	fmt.Println(humanize.Number(1234567))
	fmt.Println(humanize.Number(-1000000))
	fmt.Println(humanize.Number(123456789))
	// Output:
	// 0
	// 1,234,567
	// -1,000,000
	// 123,456,789
}

func ExamplePercent() {
	fmt.Println(humanize.Percent(0.75))
	fmt.Println(humanize.Percent(0.333))
	fmt.Println(humanize.Percent(1.0 / 3))
	fmt.Println(humanize.Percent(0))
	// Output:
	// 75%
	// 33.3%
	// 33.3%
	// 0%
}

func ExampleDuration() {
	fmt.Println(humanize.Duration(0))
	fmt.Println(humanize.Duration(500 * time.Millisecond))
	fmt.Println(humanize.Duration(90 * time.Minute))
	fmt.Println(humanize.Duration(29 * time.Hour))
	// Output:
	// 0s
	// 500ms
	// 1h 30m
	// 1d 5h
}

func ExampleParseDuration() {
	d1, _ := humanize.ParseDuration("1h 30m")
	fmt.Println(d1)

	d2, _ := humanize.ParseDuration("2d 5h")
	fmt.Println(d2)

	d3, _ := humanize.ParseDuration("500ms")
	fmt.Println(d3)
	// Output:
	// 1h30m0s
	// 53h0m0s
	// 500ms
}

func ExampleRelative() {
	ref := time.Date(2024, 2, 24, 12, 0, 0, 0, time.UTC)

	fmt.Println(humanize.Relative(ref.Add(-3*24*time.Hour), ref))
	fmt.Println(humanize.Relative(ref.Add(-7*24*time.Hour), ref))
	fmt.Println(humanize.Relative(ref.Add(-30*24*time.Hour), ref))
	fmt.Println(humanize.Relative(ref.Add(2*24*time.Hour), ref))
	// Output:
	// 3 days ago
	// 1 week ago
	// 1 month ago
	// in 2 days
}

func ExampleOrdinal() {
	fmt.Println(humanize.Ordinal(1))
	fmt.Println(humanize.Ordinal(2))
	fmt.Println(humanize.Ordinal(3))
	fmt.Println(humanize.Ordinal(11))
	fmt.Println(humanize.Ordinal(21))
	fmt.Println(humanize.Ordinal(100))
	// Output:
	// 1st
	// 2nd
	// 3rd
	// 11th
	// 21st
	// 100th
}

func ExampleCount() {
	fmt.Println(humanize.Count(1, "item", "items"))
	fmt.Println(humanize.Count(0, "item", "items"))
	fmt.Println(humanize.Count(1000, "item", "items"))
	fmt.Println(humanize.Count(-1, "person", "people"))
	// Output:
	// 1 item
	// 0 items
	// 1,000 items
	// -1 person
}
