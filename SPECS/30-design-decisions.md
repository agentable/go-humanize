# Design Decisions

This file records stable product decisions that are likely to be requested
again. `SPECS/20-api-specs.md` remains the normative behavior contract; this
file explains why common expansions stay out of scope.

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

## Rejected: Display Text Parsing

`go-humanize` does not parse humanized display text back into machine values.

Reason: even strict parsing turns display output into a protocol-shaped input
language. Callers should keep raw integers, `time.Time`, and `time.Duration`
for protocols, schedulers, billing, and audit, and use this package only at
the display boundary.

## Rejected: Pluralization Engine

`Count` accepts explicit singular and plural forms. It does not infer English
word forms or model irregular nouns.

Reason: pluralization engines become dictionary and locale systems. The caller
knows the noun forms at the call site.

## Rejected: Generic String Humanizers

The package does not provide title casing, truncation, slug formatting,
underscore splitting, or generic text cleanup.

Reason: generic string utilities would turn the package into a broad text
toolkit. `go-humanize` stays focused on numeric and time display.
