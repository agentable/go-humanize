# Go Humanize Spec

`go-humanize` formats machine values into human-readable English.

This file is the normative product contract. If examples, comments, or older commits
disagree with this file, this file wins.

> **Why**: The package is intentionally small, so one contract file is enough to define
> what belongs in the API and what does not.
>
> **Rejected**: Splitting contract rules across multiple design docs or treating the
> README as the source of truth.

## Positioning

- Small formatting library, not a general-purpose humanization toolkit
- English-only
- Zero runtime dependencies
- Stateless formatting surface
- Breaking changes are allowed when they make the API smaller, clearer, and
  more coherent
- Humanized text is for display only; never round-trip it through protocols,
  schedulers, billing, or audit

> **Why**: The package should stay narrow, predictable, and easy to explain on one
> screen.
>
> **Rejected**: Expanding into a broad text-utility bundle with multiple ways to
> express the same formatting job.

## Canonical API

Formatting:

- `Bytes(int64) string`
- `BinaryBytes(int64) string`
- `Number(int64) string`
- `Percent(float64) string`
- `Duration(time.Duration) string`
- `Relative(target, ref time.Time) string`
- `Ordinal(int64) string`
- `Count(int64, singular, plural string) string`

Parsing:

- `ParseBytes(string) (int64, error)`
- `ParseDuration(string) (time.Duration, error)`

Errors:

- Parse failures wrap `ErrInvalid`
- Callers should use `errors.Is(err, humanize.ErrInvalid)`

> **Why**: Each formatting need has one obvious entry point, and parsing stays
> explicit and opt-in.
>
> **Rejected**: Builder types, config structs, interfaces, aliases, or
> option-driven variants of the same formatter.

## Public API Admission

New exported symbols are allowed only when they preserve the package's small,
stateless shape. A proposed public API must satisfy all of these rules:

- It can be explained to a first-time caller in one sentence.
- It is not a precision knob, option-driven variant, alias, or convenience
  wrapper around an existing function.
- It does not require locale data, calendar math, time zone rules, business
  dictionaries, or caller-provided configuration.
- Its output is display-only and unsuitable for protocols, schedulers, billing,
  or audit records.
- It has a clear reason it belongs in this package instead of at the call site.
- It has table-driven tests, edge-case tests, fuzz coverage when parsing is
  involved, and a godoc example.
- The full public API still fits in the README API section without becoming a
  catalog.

> **Why**: Stable small APIs are protected by admission rules, not by adding
> options after every new request.
>
> **Rejected**: Adding a public function because it is convenient in one product
> flow, mirrors another library's broad surface, or could be configured into an
> existing formatter.

## Behavior Contract

### Bytes

- `Bytes` uses decimal units: `KB`, `MB`, `GB`, `TB`, `PB`, `EB`
- `BinaryBytes` uses IEC binary units: `KiB`, `MiB`, `GiB`, `TiB`, `PiB`, `EiB`
- Output uses at most one decimal place
- Integers do not show trailing `.0`
- The package does not export byte unit constants

> **Why**: Decimal is what user-facing products (Finder, AWS/GCP consoles,
> consumer storage) display; binary is the developer/sysadmin variant.

### ParseBytes

- Input shape is exactly `"<number> <unit>"`
- Leading or trailing whitespace is rejected
- Units must be canonical units emitted by the formatter
- Input is accepted only if it round-trips through `Bytes` or `BinaryBytes`
  unchanged
- The parsed value represents the nearest byte described by the display text,
  not the original source value before formatting
- `ParseBytes` parses display syntax, not loose user input

Examples:

- Accepted: `1 KB`, `1.5 MB`, `42 MiB`, `0 B`
- Rejected: `1KB`, `1.0 KB`, `1000 KB`, `1024 KiB`, `0 KB`, `1.5 B`

> **Why**: Byte parsing is for the library's own display language, not for
> cleaning up arbitrary user input.
>
> **Rejected**: Permissive parsing that accepts alternate spacing, aliases, or
> non-canonical values.

### Number

- `Number(int64)` formats with comma separators every three digits
- Negative values keep the leading `-`
- `math.MinInt64` is rendered explicitly because `-MinInt64` overflows

> **Why**: A function name should describe intent ("format a number for a
> reader") rather than implementation ("add commas").
>
> **Rejected**: Locale-aware separators, width parameters, or float overloads.

### Percent

- `Percent(float64)` takes a fraction: `0.333` → `"33.3%"`, `1` → `"100%"`
- Values above `1` are not clamped: `1.25` → `"125%"`
- Output uses at most one decimal place; integer percents drop the decimal
- Sub-decimal fractions that round to zero render as `"0%"`
- `NaN` renders as `"NaN%"`; `±Inf` render as `"+Inf%"` and `"-Inf%"`

> **Why**: Fraction input matches `NSNumberFormatter.percentStyle`,
> `Intl.NumberFormat({style:"percent"})`, and Excel's percent format. It also
> separates formatting from the multiply-by-100 conversion.

### Duration

- Output uses at most two descending units
- Days are treated as `24h`
- Nanoseconds below `1µs` collapse to `0s`
- Microseconds are written as `µs`

### ParseDuration

- Input is accepted only if it round-trips through `Duration` unchanged
- Leading or trailing whitespace is rejected
- Day units are supported because `time.ParseDuration` does not support them
- `ParseDuration` parses display syntax, not loose user input

Examples:

- Accepted: `1h 30m`, `2d 5h`, `500ms`, `0s`
- Rejected: `60s`, `24h`, `1000ms`, `1000µs`, `01h`, `-0s`, `500us`

> **Why**: Duration text should optimize for readability first and preserve a
> single canonical display form.
>
> **Rejected**: Unlimited precision, alternate unit spellings, or parser
> behavior that accepts values the formatter never emits.

### Relative Time

- `Relative(target, ref)` is the canonical relative-time API
- The reference time is required; there is no implicit `time.Now()` overload
- Durations under `1s` return `just now`
- Cutovers:
  `60s -> minute`, `60m -> hour`, `24h -> day`, `7d -> week`,
  `30d -> month`, `345d -> year`
- Months are `30d`
- Years are `365d`

> **Why**: Relative time should stay consistent, approximate, and easy to scan
> instead of pretending to do calendar math. Requiring an explicit reference
> keeps it testable, replayable, and timezone-stable.
>
> **Rejected**: Multiple public relative-time entry points, friendly
> vocabulary such as `yesterday` and `tomorrow`, or implicit `Now()` defaults.

### Counts

- `Count` includes comma separators by default
- `Count` is the only public count/pluralization entry point
- Human-readable output is the default path, not a composition pattern

> **Why**: Counts should read naturally without asking callers to compose
> lower-level helpers.
>
> **Rejected**: Separate public count-formatting layers or alternate count APIs
> that split formatting from pluralization.

## Non-goals

- Locale-aware formatting
- Alternate separator APIs
- Precision knobs such as `BytesN`
- Friendly calendar vocabulary such as `yesterday` or `tomorrow`
- Text utilities, truncation helpers, or string humanizers
- Builders, config structs, interfaces, or formatting options
- Implicit `time.Now()` overloads

## Coding Standards

- Use Go 1.26.3
- Use the standard library only
- Keep a single flat `humanize` package with no subpackages
- Keep formatting functions stateless
- Keep formatting functions panic-free and error-free
- Keep byte, ordinal, and count integer APIs on `int64`
- Keep output at at most one decimal place
- Name functions by the action they perform, not the implementation detail
- Model variants as separate functions, not options

> **Why**: A tiny stateless library stays maintainable when the surface area,
> naming, and package shape all remain obvious.
>
> **Rejected**: Stateful formatters, mixed integer surfaces, functional
> options, or package splits that add ceremony without new capability.

## Testing Requirements

- Use the standard library only
- Keep tests table-driven and mark tests with `t.Parallel()`
- Use `b.Loop()` in benchmarks
- Cover zero, negative, `math.MaxInt64`, `math.MinInt64`, `math.NaN()`, and
  `math.Inf(±1)` where relevant
- Keep assertions visible in the test body
- Provide at least one `Example*` function for every public function
- Keep parser fuzz tests for canonical round-trips and non-canonical rejection

> **Why**: The package surface is small enough that visible edge-case tests
> communicate behavior better than helper-heavy abstractions.
>
> **Rejected**: External assertion libraries, hidden assertion helpers, or
> layout tests that only guard documentation structure.

## Forbidden

- Do not add locale-aware formatting. Put locale behavior in a different library.
- Do not add alternate separator APIs. Keep `Number` as the integer separator path.
- Do not add precision knobs such as `BytesN`. Keep at most one decimal place.
- Do not add friendly calendar vocabulary such as `yesterday` or `tomorrow`.
  Keep relative time numeric and approximate.
- Do not add text utilities, truncation helpers, or generic string humanizers.
  Keep the package focused on numeric and time formatting.
- Do not add builders, config structs, interfaces, or formatting options. Add
  a separate function only when the behavior is genuinely distinct.
- Do not export byte unit constants. Keep the display language internal to the
  formatters and parsers.
- Do not add an implicit `time.Now()` overload to `Relative`. Always require
  an explicit reference time.
