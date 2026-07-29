# Code Surface-Lock Pruning Classification

Detection was run against all Go test files using the three signal families
from the Go language overlay. The surface-shape and error-mechanics families
returned no hits. The dependency/boundary family returned one name-based hit.

| File:line | Test/check | Signal | Claimed invariant | Observable contract | Classification | Follow-up owner |
|---|---|---|---|---|---|---|
| `number_test.go:132` | `TestPercentFiniteLargeBoundary` | Test name contains `Boundary` | Finite inputs around the large-percentage formatting threshold must not produce infinity output | Calling `Percent` with either adjacent finite input returns a finite percentage string rather than `+Inf%` or `-Inf%`; this directly enforces the documented finite-values-stay-finite runtime contract | `keep-runtime` | None |

No tests or exclusively owned scaffolding qualify for deletion. The retained
signal exercises public runtime behavior and does not inventory declarations,
imports, dependencies, or policy names.
