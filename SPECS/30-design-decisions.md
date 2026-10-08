# Design Decisions

This file records stable product decisions that are likely to be requested
again. `SPECS/20-api-specs.md` records the current behavior contract; this file
explains why common expansions stay out of scope.

## Accepted: Exact Byte Boundaries

Byte digits and unit promotion derive from the exact integer magnitude,
including `math.MinInt64`. Converting large integers to float before rounding
loses the one-byte distinction around half-points: a value one byte below
`1.25 EB` must still render as `1.2 EB`. Reject epsilon corrections and
compatibility branches; integer quotient and remainder express the boundary
directly.

Evidence: `bytes.go` and `TestByteRoundingBoundaries` in `bytes_test.go`.
Comparison paths: `.references/go-humanize/bytes.go` and
`.references/humanize/src/humanize/filesize.py` illustrate floating display
and repeated-rounding risks; their algorithms are not local contracts.

## Accepted: Decimal Percent Display Input

`Percent` interprets the float's shortest round-trip decimal representation
as its display input, then rounds ties away from zero. This matches ordinary
decimal intent: `Percent(0.5005)` is `50.1%`, while the preceding float is
`50%`. Fixed and scientific formats use the same interpretation. Scientific
mantissas round directly from the original digits, avoiding a second rounding
of the fixed-format result.

Reject exact binary64 fractions as the display interpretation, mixed rules by
magnitude, epsilon adjustments, and selectable rounding modes. The shortest
decimal conversion has a measured runtime and allocation cost; accept that
local cost rather than introduce an approximate fast path with different
boundary behavior.

Evidence: `percent.go` and `TestPercentDecimalBoundaries` in `percent_test.go`.
`.references/humanize/src/humanize/number.py` demonstrates exponent/mantissa
representation; it does not supply this package's input interpretation.

## Rejected: Generic Unit Engine

Bytes round exact integers, Percent rounds a decimal display input, Duration
truncates to two descending units, and Relative counts complete approximate
units. These are different rules. Share the two byte variants' mechanism in
`bytes.go`; keep the other domains in `percent.go`, `duration.go`, and
`relative.go`. A generic engine would hide these decisions behind parameters
without a consumer needing that abstraction.

## Rejected: Locale-Aware Formatting

`go-humanize` is English-only. Locale-aware separators, unit names, plural
rules, and percent formats belong in a different package.

Reason: locale support creates a data and configuration surface that conflicts
with the package's zero-dependency, one-obvious-way design.

## Rejected: Custom Separators

`Number` always uses comma groups of three digits.

Reason: configurable separators turn a reader-focused integer formatter into a
locale or presentation engine. Callers with custom display rules should format
their own strings.

## Rejected: Custom Precision

Byte and percent output use at most one decimal place. There is no `BytesN`,
`DurationN`, or precision option.

Reason: precision knobs make display output harder to scan and create many
near-equivalent API paths.

## Rejected: Implicit Current Time

`Relative` always requires `target` and `ref`. There is no `RelativeNow`.

Reason: explicit reference time keeps relative output testable, replayable, and
stable across time zones and clocks.

## Rejected: Friendly Calendar Words

Relative time does not emit words such as `yesterday`, `tomorrow`, or `last
week`.

Reason: friendly calendar language requires calendar math and cultural rules.
This package uses approximate numeric buckets only.

## Rejected: Early or Rounded Year Cutovers

`Relative` changes from months to years only after a complete fixed 365-day
year. It does not treat `345d`, `360d`, or another near-year threshold as one
year. Complete 30-day months remain visible up to the year boundary, so
`360d..364d` render as 12 months.

Reason: a unit label should mean that the complete unit has elapsed. Early or
nearest-unit cutovers mix rounding with truncation, make the year divisor and
cutover disagree, and weaken past/future boundary reasoning.

Evidence: `relative.go` and `relative_test.go` own the local behavior. Comparison
sources include `.references/go-humanize/times.go`,
`.references/go-humanize/times_test.go`, and
`.references/php-humanizer/src/Coduo/PHPHumanizer/DateTime/Unit/Year.php`;
their alternate thresholds are not contracts for this package.

## Rejected: Display Text Parsing

`go-humanize` does not parse humanized display text back into machine values.

Reason: even strict parsing turns display output into a protocol-shaped input
language. Callers should keep raw integers, `time.Time`, and `time.Duration`
for protocols, schedulers, billing, and audit, and use this package only at
the display boundary.

## Rejected: Pluralization Engine

`Count` accepts explicit singular and plural forms. It does not infer English
word forms or model irregular nouns. It selects singular for `1` and `-1`,
plural otherwise, and returns only the formatted number when the selected form
is empty.

Reason: pluralization engines become dictionary and locale systems. The caller
knows the noun forms at the call site. Treating an empty selected form as
absent also avoids inventing a word or emitting meaningless whitespace.

Evidence: `count.go` and `count_test.go` own the local behavior;
`.references/go-humanize/english/words.go` demonstrates why implicit fallback
word generation is outside this package's boundary.

## Rejected: Generic String Humanizers

The package does not provide title casing, truncation, slug formatting,
underscore splitting, or generic text cleanup.

Reason: generic string utilities would turn the package into a broad text
toolkit. `go-humanize` stays focused on numeric and time display.
