# Public Chart Verifier API

All paths below are repository-root-relative.

## Responsibility
The importable API of chart-verifier: check names, report and report-summary types, the verifier builder, version info, and golden samples. Consumed by `cmd/` internally and by external projects - treat everything here as a compatibility surface.

## Behavioral Model
- Inputs: builder calls on `pkg/chartverifier/verifier.NewVerifier()` (enable/disable checks, set values, flags) and a chart URI at `Run`.
- Outputs: `Report` objects serializable to YAML/JSON; report summaries via `pkg/chartverifier/reportsummary`; check-name constants.
- Flow: `cmd/verify.go` builds a verifier through this package, which delegates execution to the engine in `internal/chartverifier`.
- Failure behavior: `Run` returns an error only for invalid configuration (unknown check, bad key name); check failures are data in the report.

## Boundaries and Invariants
- Backwards compatibility is the rule: additive changes only. Renaming or removing exported identifiers breaks external importers (certification tooling imports `pkg/chartverifier`).
- Check-name constants here are the single source of truth referenced by profile YAMLs and `internal/` - update dependents in the same change.
- `pkg/chartverifier/reportsummary` output strings are compared by scripts and tests; change them only via a `proposals/` entry.
- Keep implementation out of this tree: logic belongs in `internal/`; `pkg/` defines types and delegates.

## Change Map
- To add a check name: add the constant in `pkg/chartverifier/checks/checks.go`, then update the registry in `internal/chartverifier/verifierbuilder.go` and profile YAMLs in the same change.
- To change report output: start in `pkg/chartverifier/report/types.go` and `pkg/chartverifier/reportsummary/`, then check `pkg/chartverifier/samples/` and `pkg/chartverifier/verifier/` fixtures that pin the format.
- Tests: `pkg/chartverifier/verifier/verifier_test.go`, `pkg/chartverifier/reportsummary/` tests, and `pkg/chartverifier/samples/` tests.

## Local Verification
- See the root `Makefile` test target; run `make test` from the repository root.
