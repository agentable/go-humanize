# Go Humanize Spec

`go-humanize` formats machine values into human-readable English.

This file records the current accepted product contract. When code, tests,
examples, or explicit project goals disagree, first settle the intended
behavior in code and tests, then update this file to match.

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

- `Bytes(int64) string`
- `BinaryBytes(int64) string`
- `Number(int64) string`
- `Percent(float64) string`
- `Duration(time.Duration) string`
- `Relative(target, ref time.Time) string`
- `Ordinal(int64) string`
- `Count(int64, singular, plural string) string`

> **Why**: Each formatting need has one obvious entry point. Display text is
> output, not an input grammar.
>
> **Rejected**: Builder types, config structs, interfaces, aliases, or
> option-driven variants of the same formatter. Also rejected: public parsers
> for humanized display text.

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
- It does not parse humanized display text back into machine values.
- It has a clear reason it belongs in this package instead of at the call site.
- It has table-driven tests, edge-case tests, and a godoc example.
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

- `Bytes` uses decimal units: `kB`, `MB`, `GB`, `TB`, `PB`, `EB`
- `BinaryBytes` uses IEC binary units: `KiB`, `MiB`, `GiB`, `TiB`, `PiB`, `EiB`
- Output uses at most one decimal place
- Integers do not show trailing `.0`
- The package does not export byte unit constants

> **Why**: Decimal is the reader-facing default for storage quantities; binary
> is the developer/sysadmin variant for values that are naturally measured in
> powers of two.

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
- Finite values stay finite; extremely large percentages may use scientific notation

> **Why**: Fraction input matches established percent-formatting convention and
> separates formatting from the multiply-by-100 conversion.

### Duration

- Output uses at most two descending units
- Days are treated as `24h`
- Nanoseconds below `1µs` collapse to `0s`
- Microseconds are written as `µs`

> **Why**: Duration text should optimize for readability first and preserve a
> single display form.
>
> **Rejected**: Unlimited precision, alternate unit spellings, or parsing
> display text back into `time.Duration`.

### Relative Time

- `Relative(target, ref)` is the canonical relative-time API
- The reference time is required; there is no implicit `time.Now()` overload
- Durations under `1s` return `just now`
- Cutovers:
  `60s -> minute`, `60m -> hour`, `24h -> day`, `7d -> week`,
  `30d -> month`, `365d -> year`
- Every cutover is the inclusive lower bound for its unit
- Quantities count complete units; partial units are discarded
- Months are `30d`
- Years are `365d`
- `345d..359d` render as `11 months`, `360d..364d` render as `12 months`,
  and `365d` renders as `1 year`
- Past and future differences use the same magnitude and unit; only the
  direction phrase changes
- Differences wider than `time.Duration` must retain their fixed 365-day year
  quotient instead of saturating at `time.Time.Sub`'s duration limit

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
- Counts `1` and `-1` select the caller-provided singular form; every other
  count selects the caller-provided plural form
- If the selected form is empty, `Count` returns exactly `Number(count)` and
  does not append separator whitespace
- Non-empty forms are used as provided; `Count` does not trim, infer, or
  normalize noun forms

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
- Parsing humanized display text
- Implicit `time.Now()` overloads

## Coding Standards

- Use Go 1.26.5
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

> **Why**: The package surface is small enough that visible edge-case tests
> communicate behavior better than helper-heavy abstractions.
>
> **Rejected**: External assertion libraries, hidden assertion helpers, or
> layout tests that only guard documentation structure.

## Acceptance Criteria

- Given a relative difference immediately below or at each unit cutover, when
  `Relative` formats it, then the output uses the largest unit whose complete
  lower bound has been reached. Verification path: table-driven tests in
  `relative_test.go`, especially `TestRelativeBoundaries`.
- Given the same elapsed magnitude before and after a reference time, when
  `Relative` formats both directions, then quantity and unit match while only
  `ago` versus `in` changes. Verification path: `TestRelativeSymmetry` and the
  future boundary cases in `relative_test.go`.
- Given two times whose difference exceeds `time.Duration`, when `Relative`
  formats them, then it returns the complete fixed 365-day year quotient in
  both directions. Verification path: `TestRelativeWideRange` in
  `relative_test.go`.
- Given `Count` selects an empty singular or plural form, when it formats the
  phrase, then the result equals `Number(count)` with no trailing space.
  Verification path: empty-form cases in `TestCount` in `count_test.go`.

The owning implementation paths are `relative.go` and `count.go`. Reference
evidence for boundary judgment lives in `.references/go-humanize/times.go`,
`.references/humanize/src/humanize/time.py`,
`.references/humanize/tests/test_time.py`, and
`.references/go-humanize/english/words.go`. Reference behavior is evidence,
not an expansion mandate for this package.

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
- Do not add public parsers for humanized display text. Use raw integers,
  `time.Time`, or `time.Duration` at protocol, scheduler, billing, and audit
  boundaries.
- Do not export byte unit constants. Keep the display language internal to the
  formatters.
- Do not add an implicit `time.Now()` overload to `Relative`. Always require
  an explicit reference time.
