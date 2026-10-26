---
name: add-chart-check
description: Add a new verification check to chart-verifier end to end - implementation, registration, profiles, docs, and tests. Use when a task introduces or modifies a chart check.
---

# Add a Chart Check

## Preconditions
- Read go.mod for information on project golang version.
- Know the check's kebab-case name, its certification type (Mandatory/Optional), and which vendor profiles include it.
- Confirm the name does not already exist in `pkg/chartverifier/checks/checks.go`.
- Determine the implementation details of the check.

## Procedure
1. Add the `CheckName` constant in `pkg/chartverifier/checks/checks.go`.
2. Implement the check as a `CheckFunc` in `internal/chartverifier/checks/checks.go`:
   - signature `func MyCheck(opts *CheckOptions) (Result, error)`;
   - return `NewSkippedResult` when the check does not apply, `NewResult` with an actionable reason otherwise;
   - shell out only through `internal/tool`.
3. Register it in `internal/chartverifier/verifierbuilder.go` `init()`: `defaultRegistry.Add(apiChecks.MyCheck, "v1.0", checks.MyCheck)`. If it modifies an existing check incompatibly, add a `"v1.1"` entry with a `MyCheck_V1_1` function and keep the old one.
4. Add the check to each current-version profile YAML in `internal/profileconfig/profiles/` with its type. Do not edit already released profile versions.
5. Add a test in `internal/chartverifier/checks/checks_test.go` using a fixture chart from `internal/chartverifier/checks/` - add a new fixture `.tgz` if no existing one exercises the condition.
6. Document the check in `docs/helm-chart-checks.md` (description, type, how to satisfy it).
7. If the default profile has been modified, update `internal/chartverifier/profiles/default.go` to ensure that it is served when no profile is specified.

These steps modify source files only; nothing here pushes, publishes, or touches remote systems.

## Decision Points
| Condition | Action |
| --- | --- |
| Check applies only to some charts (e.g. signed) | Return `NewSkippedResult` with a clear reason instead of failing |
| Behavior change breaks old reports | Register a new check version; never change an existing version's semantics |
| Check needs a cluster or helm install | Use chart-testing via the existing `ChartTesting` helpers; do not add direct kubectl calls |
| New fixture needed | Ask for guidance on how to generate the needed fixture. State your parameters. Then build it from a source chart and commit the `.tgz`; never edit an existing fixture in place |

## Completion Checks
- [ ] Run `make test` from the repository root using the root `Makefile` test target; record any failures.
- [ ] New check appears in all current-version profile YAMLs and in `docs/helm-chart-checks.md`
- [ ] `make lint` reports 0 issues
- [ ] `make fmt` reports 0 issues.

## Conditional References
- For check semantics and certification types, read `docs/agentrefs/chart-checks.md` and `docs/helm-chart-checks.md`.
- For profile versioning rules, read `docs/agentrefs/profile-configuration.md`.
