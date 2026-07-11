# Go Humanize

Small, strict Go helpers for turning machine values into human-readable
English display strings.

`go-humanize` stays intentionally narrow: bytes, durations, relative time,
numbers, ordinals, percentages, and count phrases. It has no runtime
dependencies, no builders, no option structs, and no locale layer.

> Humanized text is for display only. Never feed humanized text back into
> protocols, schedulers, billing, or audit. Use raw integers, `time.Time`,
> or `time.Duration` for those.

## Installation

```bash
go get github.com/agentable/go-humanize
```

Requires **Go 1.26.5**.

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
	twoHoursAgo := ref.Add(-2 * time.Hour)

	fmt.Println(humanize.Bytes(1500))                  // "1.5 kB"
	fmt.Println(humanize.BinaryBytes(1536))            // "1.5 KiB"
	fmt.Println(humanize.Relative(twoHoursAgo, ref))   // "2 hours ago"
	fmt.Println(humanize.Number(1234567))              // "1,234,567"
	fmt.Println(humanize.Ordinal(3))                   // "3rd"
	fmt.Println(humanize.Percent(0.333))               // "33.3%"
	fmt.Println(humanize.Count(1000, "item", "items")) // "1,000 items"
}
```

## API

| Function | Purpose | Example |
|---|---|---|
| `Bytes(int64) string` | Decimal byte units | `1500` → `"1.5 kB"` |
| `BinaryBytes(int64) string` | IEC binary byte units | `1536` → `"1.5 KiB"` |
| `Number(int64) string` | Comma-separated integers | `1234567` → `"1,234,567"` |
| `Percent(float64) string` | Fraction to percentage | `0.333` → `"33.3%"` |
| `Duration(time.Duration) string` | Up to two descending units | `90*time.Minute` → `"1h 30m"` |
| `Relative(target, ref time.Time) string` | Relative time from an explicit reference | `"2 hours ago"` |
| `Ordinal(int64) string` | English ordinal suffix | `3` → `"3rd"` |
| `Count(int64, singular, plural string) string` | Count phrase with comma separators | `1000` → `"1,000 items"` |

## Behavior

### Bytes

- `Bytes` uses decimal units: `kB`, `MB`, `GB`, `TB`, `PB`, `EB`.
- `BinaryBytes` uses IEC units: `KiB`, `MiB`, `GiB`, `TiB`, `PiB`, `EiB`.
- Output uses at most one decimal place and drops trailing `.0`.

### Duration

- `Duration` shows at most two descending units.
- Days are treated as 24 hours.
- Sub-microsecond durations collapse to `"0s"`.

### Relative

- `Relative(target, ref)` always requires an explicit reference time. There
  is no implicit `time.Now()` overload.
- Cutovers: `60s` to minute, `60m` to hour, `24h` to day, `7d` to week,
  `30d` to month, `345d` to year.
- Months are 30 days. Years are 365 days. There is no calendar math.

### Percent

- `Percent` takes a fraction: `Percent(0.333)` → `"33.3%"`,
  `Percent(1)` → `"100%"`. Values above `1` are not clamped.
- `NaN` and infinities render as `"NaN%"`, `"+Inf%"`, `"-Inf%"`.
- Extremely large finite percentages may use scientific notation, but they do
  not render as infinity.

## Non-Goals

- Locale-aware formatting
- Alternate separator APIs
- Precision knobs like `BytesN`
- Friendly calendar words like `yesterday` or `tomorrow`
- Text utilities, truncation helpers, or string humanizers
- Builders, config structs, or formatting options
- Parsing humanized display text

See [`SPECS/30-design-decisions.md`](./SPECS/30-design-decisions.md) for the
standing rejection record behind these boundaries.

## Why It Stays Small

Other humanize libraries grow into broad utility bundles. This package does
the opposite on purpose: one obvious entry point per need, a short list of
exported symbols, and behavior that can be explained on one screen.

`SPECS/20-api-specs.md` is the product contract. New public API must pass the
admission rules in
[`SPECS/20-api-specs.md`](./SPECS/20-api-specs.md).

## Development

```bash
task test          # Run all tests with race detection
task lint          # Run golangci-lint v2.12.2 + go mod tidy check
task fmt           # Format code
task vet           # Run go vet
task bench         # Run benchmark baseline
task vuln          # Run govulncheck
task verify        # Full verification: deps, fmt, vet, lint, test, vuln
task deps          # Download and tidy dependencies
task clean         # Clean build artifacts and caches
```

## License

Released under the [MIT License](./LICENSE).
