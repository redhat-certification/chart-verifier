# Chart Checks

All paths below are repository-root-relative.

## Responsibility
Implements every verification check the verifier can run, plus the fixture charts used to test them. Does not own orchestration (see `internal/chartverifier/verifier.go`) or profile membership (see `internal/profileconfig/`).

## Behavioral Model
- Inputs: a `CheckOptions` struct - chart URI, Viper config, user values, Helm env settings, timeout, GPG public keys for signed charts.
- Outputs: `Result{Ok, Skipped, Reason}` - checks never abort a run; failures are outcomes, not errors. Return an `error` only for unusable inputs.
- Flow: each check is a standalone `CheckFunc` registered by name+version in `internal/chartverifier/verifierbuilder.go` (`DefaultRegistry`). The verifier resolves the profile's check list against that registry and runs each entry independently; there is no guaranteed execution order and checks must not depend on each other.
- Failure behavior: use `NewSkippedResult` when the check does not apply (e.g. unsigned chart); `NewResult(false, reason)` when it applies and fails. Reasons are user-facing strings - keep them actionable.

## Boundaries and Invariants
- Check names are a public contract: defined in `pkg/chartverifier/checks/checks.go`, referenced by profile YAMLs and downstream tooling. Add a new `CheckName` constant there, never rename an existing one.
- To change check behavior incompatibly, register a new version (`"v1.1"`, function `MyCheck_V1_1`) and keep the old function - the pattern used by `HasKubeVersion` / `ImagesAreCertified`.
- Checks may shell out through `internal/tool` (helm, kubectl, gpg); do not invoke binaries directly with `exec.Command` from here.
- Fixture charts (`chart-0.1.0-v3.*.tgz`, `psql-service*`) are committed artifacts named for the condition they exercise; add a new fixture rather than editing an existing one.

## Change Map
- To add a check: implement the `CheckFunc` in `internal/chartverifier/checks/checks.go`, add the `CheckName` in `pkg/chartverifier/checks/checks.go`, register it in `internal/chartverifier/verifierbuilder.go`, add it to the profile YAMLs in `internal/profileconfig/profiles/`, and document it in `docs/helm-chart-checks.md`. The `add-chart-check` skill walks through this.
- Tests for check behavior live in `internal/chartverifier/checks/checks_test.go` (and `internal/chartverifier/checks/charttesting_test.go`, `internal/chartverifier/checks/helm_test.go` for chart-testing-based checks).
- Do not edit registry semantics in `internal/chartverifier/checks/registry.go` without reading how `FilteredRegistry` in `internal/chartverifier/verifierbuilder.go` resolves names and versions.

## Local Verification
- See the root `Makefile` test target; run `make test` from the repository root.

## Conditional References
- For what each check means to certification, read `docs/helm-chart-checks.md`.
- For check-configuration propagation design, read `proposals/P99999-propagate-user-defined-configuration-to-checks.md`.
