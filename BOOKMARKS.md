# Bookmarks

Curated references for task-specific context. Use the descriptions to decide what to load; these references are not intended to be read all at once. Paths are relative to the repository root.

## Architecture and behavior

- [Verifier engine](docs/agentrefs/verifier-engine.md) - verification flow, registry/profile resolution, report invariants, and cross-directory signing behavior. Read when changing orchestration, report generation, or chart signing.
- [Chart checks](docs/agentrefs/chart-checks.md) - check result semantics, compatibility/versioning, fixtures, and the end-to-end change map. Read when adding or modifying a check.
- [Profile configuration](docs/agentrefs/profile-configuration.md) - profile membership, vendor types, and versioning rules. Read when changing profile data or defaults.

## API and design contracts

- [Public API](docs/agentrefs/public-api.md) - exported compatibility constraints and the report/check-name change map. Read before changing `pkg/chartverifier/` or report output.

## Development and release

- [Check documentation](docs/helm-chart-checks.md) - user-facing check semantics and certification requirements. Consult when changing check behavior or documentation.
- [Python tooling README](scripts/README.md) and [`ruff.toml`](ruff.toml) - scripts setup and lint conventions. Read when modifying Python tooling.
