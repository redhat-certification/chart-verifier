# Profile Configuration

All paths below are repository-root-relative.

## Responsibility
Defines which checks and annotations apply for each vendor type (partner, redhat, community) and profile version. Loaded by `internal/chartverifier/profiles` and consulted by `cmd/verify.go`.

## Behavioral Model
- Inputs: YAML files in `internal/profileconfig/profiles/`, selected at runtime by `--profile` and vendor type.
- Outputs: the set of `v<version>/<name>` check identifiers plus annotation requirements that the verifier resolves against the check registry.
- Flow: `cmd/verify.go` loads the profile via this package, filters the default registry, and runs the resulting checks.
- Failure behavior: an unknown check name in a profile is **silently dropped** by `Profile.FilterChecks` (`internal/chartverifier/profiles/profile.go`) - the check never runs; nothing errors at load or run time. Validate names after edits.

## Boundaries and Invariants
- Check entries must reference names registered in `internal/chartverifier/verifierbuilder.go`; a typo silently breaks the profile.
- Profiles are versioned (`v1.0` … `v1.3`). Never edit an existing version in place once released - create the next version file for all three vendor types and change the default in `internal/chartverifier/profiles/default.go`.
- `type:` (Mandatory/Optional/Recommended) drives certification outcomes; see `docs/helm-chart-checks.md` for semantics.
- Profile **developer-console** is a deprecated profile. The primary profiles are **partner**, **redhat**, and **community**

## Change Map
- To include a new check: add a `checks:` entry to **each current-version** profile YAML (partner, redhat, community), then document it in `docs/helm-chart-checks.md` "Default set of checks".
- To evolve a profile: copy the newest version file, bump `version:`, apply changes, and update `internal/chartverifier/profiles/default.go`.
- Tests live in `internal/chartverifier/profiles/profile_test.go`.

## Local Verification
- See the root `Makefile` test target; run `make test` from the repository root.
