# Chart Verifier Engine

All paths below are repository-root-relative.

## Responsibility
Orchestration and report generation: builds the verifier, runs profile-selected checks, and assembles the API report. Check implementations live in `internal/chartverifier/checks/`; profile YAMLs in `internal/profileconfig/`; public types in `pkg/chartverifier/`.

## Behavioral Model
- Inputs: a `FilteredRegistry` of selected checks (see below), a chart URI, Viper config, and run options via `VerifierBuilder` (`internal/chartverifier/helmverifier.go`).
- Outputs: `apiReport.Report` (defined in `pkg/chartverifier/report`), optionally rendered as JUnit XML by `internal/chartverifier/junitxml/`.
- Flow: `internal/chartverifier/verifier.go` `Verify(uri)` loads the chart and runs each check independently with a fresh `CheckOptions` (per-check Viper subconfig, `AnnotationHolder` for OpenShift-version annotations), accumulating results via `internal/chartverifier/reportBuilder.go`. `webCatalogOnly` runs require a tarball; `internal/chartverifier/pyxis/` queries the Red Hat container catalog for `images-are-certified`. A nil check function or check error aborts the run; check failures are outcomes recorded in the report.

## Registry Resolution
`internal/chartverifier/verifierbuilder.go` `init()` registers every check in `defaultRegistry` by name + version. `FilteredRegistry` (name → check) is built by `internal/chartverifier/profiles.FilterChecks`, which looks up each profile entry `v<version>/<name>` in the registry. Unknown entries are **silently dropped** - a profile typo removes the check. The map is keyed by name; when a profile lists two versions, the later entry wins.

## Report Contract Invariants
- Digests (`internal/chartverifier/reportBuilder.go`): chart sha (`GenerateSha` over sorted files, excluding verifier output), package digest, public-key digest, and report digest (`SetReportDigest`).
- Metadata (annotations, certified/tested OpenShift versions) is driven by the profile's annotation list, not hardcoded here.
- Report structure should NOT be changed without discussion.
- Summaries and digest verification live in `pkg/chartverifier/reportsummary` (via `cmd/report.go`), which recomputes the report digest to detect tampering; `junitxml` embeds it too. This package generates digests; that one checks them.

## Internal Report Fixtures
`internal/chartverifier/report/` contains only static sample report YAMLs under `internal/chartverifier/report/testreports/v1.0` and `v1.1` - no Go code, and currently unreferenced by tests. Do not confuse it with the public `pkg/chartverifier/report` package that defines the report types.

## Chart Signing (cross-directory routing)
No single directory owns signing - when touching it, read all of:
- `internal/chartverifier/checks/` - the `signature-is-valid` check (`SignatureIsValid`).
- `internal/tool/pgp.go` - GPG keyring creation (`GetKeyRing`), `GetPublicKeyDigest`; key order matters (covered by `pkg/chartverifier/verifier/verifier_test.go`'s `TestSignedChart`).
- `docs/helm-chart-checks.md#signed-charts` - user-facing semantics.
- `internal/chartverifier/verifier.go` sets the public-key digest after the signature check when one key is supplied and the chart verifies as signed.

## Local Verification
- See the root `Makefile` test target; run `make test` from the repository root.

## Conditional References
- `docs/agentrefs/chart-checks.md` (check semantics); `docs/agentrefs/profile-configuration.md` (profiles); `docs/agentrefs/public-api.md` (public API).
