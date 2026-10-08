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

Requires **Go 1.27.0**.

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
- Scaled values below 10 use one decimal place; values from 10 use whole
  numbers. Ties round away from zero, trailing `.0` is omitted, and rounding
  into the next unit promotes the suffix: `Bytes(1250)` is `"1.3 kB"` and
  `Bytes(999500)` is `"1 MB"`.
- Rounding uses the exact integer magnitude throughout the `int64` range:
  `Bytes(1249999999999999999)` is `"1.2 EB"`, one byte below the half-point.

### Duration

- `Duration` shows at most two descending units.
- Days are treated as 24 hours.
- Sub-microsecond durations collapse to `"0s"`.

### Relative

- `Relative(target, ref)` always requires an explicit reference time. There
  is no implicit `time.Now()` overload.
- Cutovers: `60s` to minute, `60m` to hour, `24h` to day, `7d` to week,
  `30d` to month, `365d` to year.
- Cutovers are lower bounds and quantities count complete units. For example,
  `359d` is 11 months, `364d` is 12 months, and `365d` is 1 year.
- Months are 30 days. Years are 365 days. There is no calendar math.
- Differences wider than `time.Duration` still render their fixed 365-day
  year count instead of saturating.

### Count

- `Count` uses the singular form for `1` and `-1`, and the plural form for
  every other count.
- If the selected form is empty, `Count` returns only the comma-separated
  number, with no trailing space.

### Percent

- `Percent` takes a fraction: `Percent(0.333)` → `"33.3%"`,
  `Percent(1)` → `"100%"`. Values above `1` are not clamped.
- It rounds the shortest round-trip decimal input with ties away from zero:
  `Percent(0.5005)` → `"50.1%"`.
- `NaN` and infinities render as `"NaN%"`, `"+Inf%"`, `"-Inf%"`.
- Rounded percentage magnitudes from `1e6` use scientific notation:
  `Percent(10000)` returns `"1e+6%"`. Finite inputs do not render as infinity.

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
task lint          # Run golangci-lint v2.13.1 + go mod tidy check
task fmt           # Format code
task vet           # Run go vet
task bench         # Run benchmark baseline
task vuln          # Run govulncheck
task verify        # Full verification: deps, fmt, vet, lint, test, vuln
task deps          # Download and tidy dependencies
task clean         # Clean build artifacts and caches
```

## Release

After `task verify` passes, commit the release changes on `main`:

```bash
bash scripts/check_changes.sh     # 0: release needed; 1: no relevant committed changes
bash scripts/tag_release.sh X.Y.Z # Create vX.Y.Z and atomically push main + this tag
```

Change detection compares the latest local `v*` tag with `HEAD`, including
the complete `go.mod` and `go.sum`. The release script binds the tag to the
captured `main` commit and pushes only that commit and tag. A rejected push
leaves the local tag intact; resolve the rejection before retrying the same
atomic push. See [the release contract](./SPECS/20-api-specs.md#release-scripts).

## License

Released under the [MIT License](./LICENSE).
