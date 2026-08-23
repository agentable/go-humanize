# Go Humanize Library

Converts machine values into human-readable English strings. Small API, zero
runtime dependencies, stateless formatting surface.

`SPECS/20-api-specs.md` records the current product contract. `CLAUDE.md`
describes coding and workflow rules for agents working in this repo.

## SPECS Index

| Spec | Purpose |
|------|---------|
| `SPECS/20-api-specs.md` | Public API, formatter behavior, coding rules, and testing requirements |
| `SPECS/30-design-decisions.md` | Rejected design expansions and long-term scope boundaries |

Before changing API behavior, formatter rules, relative-time cutovers, or
package scope, read `SPECS/20-api-specs.md` first. Before adding anything that
looks like locale support, custom precision, input parsing, calendar language,
or string utilities, read `SPECS/30-design-decisions.md`.

## References Index

| Reference | Use |
|---|---|
| `.references/go-humanize/` | Go API and formatter comparison evidence |
| `.references/humanize/` | Python behavior and wide-range time comparison evidence |
| `.references/php-humanizer/` | Independent unit-boundary comparison evidence |

```go
import (
	"time"

	"github.com/agentable/go-humanize"
)

humanize.Bytes(1500)                          // "1.5 kB"
humanize.BinaryBytes(1536)                    // "1.5 KiB"
ref := time.Now()
twoHoursAgo := ref.Add(-2 * time.Hour)
humanize.Relative(twoHoursAgo, ref)           // "2 hours ago"
humanize.Number(1234567)                      // "1,234,567"
humanize.Ordinal(3)                           // "3rd"
humanize.Percent(0.333)                       // "33.3%"
humanize.Count(1000, "item", "items")         // "1,000 items"
```

## Agent Workflow

1. Read the relevant SPECS completely before designing or changing behavior.
2. For new behavioral design, inspect at least two relevant projects under
   `.references/` when applicable; treat them as evidence, not a feature list.
3. Establish RED through public observable behavior. Do not manufacture RED
   with tests of source text, private call order, or file layout.
4. Keep the public API small and explicit; do not add options, builders,
   interfaces, aliases, or convenience wrappers.
5. Keep stability coverage in the existing domain test files
   (`bytes_test.go`, `duration_test.go`, and so on). Do not create a separate
   API-wide test file or extra verification task.
6. Use table-driven stdlib tests with visible assertions.
7. Run the narrowest useful test while developing, then run `task verify`
   before handing off broad changes.

## Agent Operating Rules

- Read owning source, tests, SPECS, and relevant references before writing.
- Keep changes surgical and prove the requested user path end to end.
- Apply KISS, DRY, and YAGNI as judgment; share knowledge, not coincidental lines.
- Prefer the standard library and existing project patterns before custom machinery.
- Do not turn prose policies into custom gate scripts or redundant mirror tests.
- Preserve failing evidence and surface blockers explicitly instead of weakening checks.

## Commands

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

```bash
bash scripts/check_changes.sh   # Exit 0 when a release tag is needed
bash scripts/tag_release.sh X.Y.Z  # Create vX.Y.Z and push main + tags
```

## Package Shape

Single package with one file per domain:

```text
go-humanize/
├── bytes.go
├── number.go
├── percent.go
├── duration.go
├── relative.go
├── ordinal.go
├── count.go
└── doc.go
```

One file per domain. No `utils.go`, no `helpers.go`, no `common.go`.

## Design Philosophy

- **Stateless Formatting, No State** — Public formatting functions have no mutable package state.
- **One Obvious Way** — Each formatting need has exactly one function. Variants are separate functions (`Bytes` vs `BinaryBytes`), not option parameters.
- **Graceful on All Inputs** — Every function handles zero, negative, and edge-case values without panicking. Public functions never return errors.
- **Readability Over Precision** — Output optimized for human comprehension. Duration shows at most two significant units. Relative time uses approximate months (30 days) and years (365 days).
- **No Compatibility Tax** — Backward compatibility is not a goal when it conflicts with a smaller and cleaner API surface.
- **Documented Refusals** — Common expansion requests live in `SPECS/30-design-decisions.md` so the package can reject them consistently.

## Public API

Formatting functions never return errors:

- `Bytes`, `BinaryBytes`
- `Number`, `Percent`
- `Duration`, `Relative`
- `Ordinal`, `Count`

## Coding Rules

### Must Follow

- **Go 1.27.0** — use modern builtins (`min`, `max`, `clear`), `slices`/`maps` packages, `for range N`, `b.Loop()` in benchmarks
- **Zero dependencies** — standard library only, no external imports
- **Single package** — flat namespace under `package humanize`, no subpackages
- **English-only** — no i18n/l10n, locale-aware formatting is a different library's job
- **Stateless formatting** — no package-level mutable state, no constructors, no interfaces
- **Graceful edge cases** — format functions never panic or return errors
- **At most one decimal place** — `"1.5 kB"` when needed, `"1 kB"` when not
- **Function naming** — name describes the intent (`Number` formats integers for a reader), no `Format` prefix, no `To` prefix
- **Variants are separate functions** — `Bytes` vs `BinaryBytes`, not `Bytes(n, WithBinary())`
- **Single relative-time API** — `Relative` is the only relative-time entry point; reference time is required
- **`int64` for integer APIs** — bytes, ordinals, and counts all use `int64` to avoid mixed integer surfaces
- **Human-readable counts by default** — `Count` includes comma separators

### Forbidden

- **No errors from public functions** — display formatting is total and panic-free
- **No `fmt.Stringer` types** — return plain `string`, don't wrap in custom types
- **No interfaces** — the API is stateless functions, no `Formatter`, `Humanizer`, or `Renderer` interfaces
- **No functional options** — separate functions for variants, not option parameters
- **No display parsers** — humanized text is output, not a public input grammar
- **No `sync.Pool`** — unless profiling proves allocation pressure
- **No reflection** — type-safe code throughout
- **No sub-second relative time** — `< 1 second` → `"just now"`
- **No calendar math** — days are 24h, months are 30 days, years are 365 days
- **No premature abstraction** — three similar lines are better than a helper used once
- **No documentation masquerading as code** — do not encode prose rules in values no program consumes
- **No policy-only gates or spec-mirror tests** — verify runtime behavior instead of markdown wording or source layout
- **No dependency workarounds** — report dependency defects in `reports/<dependency-name>.md` instead of reimplementing them

## Dependency Issue Reporting

When a tool or dependency is defective, record its name and version, trigger,
expected and actual behavior, relevant output, and a suggested upstream fix in
`reports/<dependency-name>.md`. Continue only with work that does not depend on
the defect; do not replace the dependency with a local imitation.

## Testing

- **Framework:** stdlib assertions only (no testify)
- **Patterns:** table-driven tests, `t.Parallel()` in all tests, `b.Loop()` in benchmarks (Go 1.27.0)
- **Coverage:** all edge cases (zero, negative, `math.MaxInt64`, `math.MinInt64`, `math.NaN()`, `math.Inf(±1)`)
- **Examples:** every public function has at least one `Example*` function for godoc
- **No test helpers that hide assertions** — use `t.Errorf` directly, keep test logic visible

## Agent Skills

| Skill | When to Use |
|---|---|
| `.agents/skills/go-best-practices/` | Go naming, edge cases, tests, and idioms affect an implementation |
| `.agents/skills/tdd-implementing/` | A behavior change can be driven through public RED-GREEN evidence |
| `.agents/skills/code-review/` | Reviewing a completed change for regressions, safety, and missing tests |
| `.agents/skills/improvement-proposing/` | Revising evidence-backed `improve.md` proposals |
| `.agents/skills/spec-writing/` | Updating durable product contracts and acceptance criteria |
| `.agents/skills/committing/` | Staging a narrow verified conventional commit |
| `.agents/skills/releasing/` | Choosing, tagging, pushing, and verifying a Go module release |
