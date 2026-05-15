# Go Humanize Library

`go-humanize` formats machine values into human-readable English.

It stays intentionally narrow: bytes, durations, relative time, numbers,
ordinals, percentages, and count phrases. No builders, no option structs, no
locale layer, and no runtime dependencies.

> Humanized text is for display only. Never feed humanized text back into
> protocols, schedulers, billing, or audit. Use raw integers, `time.Time`,
> or `time.Duration` for those.

## Installation

```bash
go get github.com/agentable/go-humanize
```

Requires **Go 1.26.2**.

## Quick Taste

```go
package main

import (
	"fmt"
	"time"

	"github.com/agentable/go-humanize"
)

func main() {
	ref := time.Date(2024, 2, 24, 12, 0, 0, 0, time.UTC)

	fmt.Println(humanize.Bytes(1500))                        // "1.5 KB"
	fmt.Println(humanize.BinaryBytes(1536))                  // "1.5 KiB"
	fmt.Println(humanize.Relative(ref.Add(-2*time.Hour), ref)) // "2 hours ago"
	fmt.Println(humanize.Number(1234567))                    // "1,234,567"
	fmt.Println(humanize.Percent(0.333))                     // "33.3%"
	fmt.Println(humanize.Count(1000, "item", "items"))       // "1,000 items"
}
```

## API

- `Bytes(int64) string` — decimal units (KB, MB, GB, ...)
- `BinaryBytes(int64) string` — IEC binary units (KiB, MiB, GiB, ...)
- `Number(int64) string` — comma-separated integers
- `Percent(float64) string` — takes a fraction; `0.333` → `"33.3%"`
- `Duration(time.Duration) string` — at most two descending units
- `Relative(target, ref time.Time) string` — relative-time formatter
- `Ordinal(int64) string` — English ordinal suffix
- `Count(int64, singular, plural string) string` — count phrase
- `ParseBytes(string) (int64, error)` — inverse of `Bytes` / `BinaryBytes`
- `ParseDuration(string) (time.Duration, error)` — inverse of `Duration`
- `ErrInvalid` — sentinel returned by parse failures

## Behavior

### Bytes

- `Bytes` is decimal by default because Finder, AWS/GCP consoles, and most
  user-facing products use decimal byte units.
- `BinaryBytes` is for the developer/sysadmin context (memory, block devices).
- Output uses at most one decimal place.
- `ParseBytes` accepts only text that round-trips through `Bytes` or
  `BinaryBytes` unchanged. The parsed value represents the nearest byte
  described by the display text, not the original source value before
  formatting.

### Duration

- `Duration` shows at most two descending units.
- Days are treated as 24 hours.
- Sub-microsecond durations collapse to `"0s"`.
- `ParseDuration` accepts only text that round-trips through `Duration`
  unchanged. Day units (`d`) are supported because `time.ParseDuration`
  cannot parse them.

### Relative

- `Relative(target, ref)` always requires an explicit reference time. There
  is no implicit `time.Now()` overload.
- Cutovers: `60s → minute`, `60m → hour`, `24h → day`, `7d → week`,
  `30d → month`, `345d → year`. Months are `30d`; years are `365d`.

### Percent

- `Percent` takes a fraction: `Percent(0.333)` → `"33.3%"`,
  `Percent(1)` → `"100%"`. Values above `1` are not clamped.
- `NaN` and infinities render as `"NaN%"`, `"+Inf%"`, `"-Inf%"`.

### Errors

- Parse failures wrap `ErrInvalid`. Use
  `errors.Is(err, humanize.ErrInvalid)`.

## Non-Goals

- Locale-aware formatting
- Alternate separator APIs
- Precision knobs like `BytesN`
- Friendly calendar words like `yesterday` or `tomorrow`
- Text utilities, truncation helpers, or string humanizers
- Builders, config structs, or formatting options

## Why It Stays Small

Other humanize libraries grow into broad utility bundles. This package does
the opposite on purpose: one obvious entry point per need, a short list of
exported symbols, and a behavior contract that can be explained on one screen.

## Development

```bash
task test          # Run all tests with race detection
task lint          # Run golangci-lint v2.11.4 + go mod tidy check
task fmt           # Format code
task vet           # Run go vet
task verify        # Full verification: deps, fmt, vet, lint, test
task deps          # Download and tidy dependencies
task clean         # Clean build artifacts and caches
```

## License

Released under the [MIT License](./LICENSE).
