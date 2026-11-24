# Repository Agent Guide

## Purpose
chart-verifier is a Go CLI that certifies Helm charts for Red Hat OpenShift. It runs a configurable set of checks (via `verify`), emits a YAML/JSON report, and can summarize existing reports (via `report`). Check sets are driven by versioned profiles (partner, redhat, community).

## Working in This Repository
- Prerequisites: Go 1.26 (see `go.mod`); `gpg` on PATH for signed-chart tests; Python 3 for `scripts/` and functional tests.
- Setup: `make tidy`
- Build: `make bin` → `./out/chart-verifier` (CGO disabled, version stamped via ldflags).
- Tests: see the root `Makefile` test target; run `make test` from the repository root.
- Lint: `make lint` (installs pinned golangci-lint via `go install`).
- Format: `make fmt` (gofumpt)

## Repository-Wide Constraints
- Do not break the public API in `pkg/` - it is imported by external projects (e.g. certification tooling). Additive changes only unless a proposal says otherwise.
- Report structure is an interoperability contract with chart submission tooling.
- Check names (`v<version>/<name>`) appear in profile YAMLs, docs, and downstream tooling; never rename an existing check - add a new version instead (pattern: `HasKubeVersion` v1.0 → v1.1).
- Do not edit anything under `tests/charts/` fixture `.tgz` files by hand - they are inputs checked in as-is; regenerate from source charts when a fixture must change.
- `make tidy` and `make fmt` may rewrite files and then require a clean diff; run them deliberately.
- Container images build with `podman` by default (`IMAGE_BUILDER` overridable); never run `make push-image` - publishing is a release activity.

## Task-Specific Context
- Before task-specific code or design work, consult [BOOKMARKS.md](BOOKMARKS.md). It is the curated, on-demand index for architecture, development, API-compatibility, and design or user-facing references. Load only entries relevant to the task rather than duplicating their details here.
