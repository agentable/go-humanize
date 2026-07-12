# Design Decisions

This file records stable product decisions that are likely to be requested
again. `SPECS/20-api-specs.md` records the current behavior contract; this file
explains why common expansions stay out of scope.

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

Evidence: `time.go` and `time_test.go` own the local behavior. Comparison
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
