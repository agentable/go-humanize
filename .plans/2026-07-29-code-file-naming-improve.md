# Code File Naming Improvement

## Scope And Conventions

- Scope: all 14 tracked, hand-written Go files in the single `humanize` package.
  Gitlink contents under `.agents/` and `.references/` are outside this audit.
- Convention: lowercase ASCII subject names, one production file per formatter
  domain, source-owned tests paired by basename, and Go-reserved documentation
  and example owners preserved.
- No generated Go files or filename-sensitive build, embed, or generate
  directives exist in the main-repository tree.
- Baseline SHA: `840072d35e14c52b72ccc7419512f5343567466a`.
  `task verify` passed with formatting, vet, module tidiness, lint, race tests,
  and vulnerability checks enabled.

## Accepted Renames

None. Every current path is already the shortest truthful subject under the
repository's explicit one-file-per-domain architecture.

## Explicit Keeps

| Path | Primary subject | Reason to keep |
|---|---|---|
| `bytes.go` / `bytes_test.go` | Decimal and binary byte formatting | One cohesive byte-display domain with an exact source/test pair |
| `count.go` / `count_test.go` | Count phrases and singular/plural selection | Names the public `Count` formatter and its directly owned behavior |
| `duration.go` / `duration_test.go` | Duration formatting | Names the public `Duration` formatter and its directly owned behavior |
| `number.go` / `number_test.go` | Reader-facing numeric formatting, including percentages | `Number` and `Percent` share one stable numeric-display concern; splitting is not required for truthful naming |
| `ordinal.go` / `ordinal_test.go` | English ordinal formatting | Names the public `Ordinal` formatter and its directly owned behavior |
| `time.go` / `time_test.go` | Relative-time formatting | The repository explicitly assigns `Relative` and its fixed-time arithmetic to the time domain |
| `doc.go` | Package documentation | Go-reserved package documentation owner |
| `example_test.go` | Package-wide godoc examples | Intentional Go example suite covering the complete public surface |

The names also match the package-shape contract documented in `CLAUDE.md`.
Changing `time.go` to `relative.go`, or splitting `number.go` solely to achieve
one exported function per file, would introduce a second convention without a
clearer table-of-contents entry.

## Routed Structural Work

None. The inventory contains no multi-subject generic bucket, cross-package
move, generated owner, or misleading symbol/API name that must be routed from
this path-only audit.

## Validation

- Run `task verify`, repository-wide coverage, and every benchmark with
  `GOTOOLCHAIN=go1.26.5` and `GOWORK=off`.
- Build for Linux and Windows amd64 with CGO disabled.
- Search local Go modules for real downstream consumers and test any found.
- Confirm `git diff --check` passes and the commit contains only this audit
  evidence.

## Validation Evidence

- `task verify` passed with `GOTOOLCHAIN=go1.26.5` and `GOWORK=off`.
- Repository-wide coverage passed at 97.8% of statements.
- Every benchmark passed under the repository's race-enabled `task bench`.
- Linux and Windows amd64 no-CGO builds passed.
- The three local source consumers (`aq`, `filter`, and `go-filter`) passed
  their complete tests through a temporary `go.work` that selected the audited
  local `go-humanize` tree without modifying consumer manifests.
- `git diff --check` passed; no path changes are present.
