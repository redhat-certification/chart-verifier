# Flags and side effects

Determine which [checks](checks.md) the user wants before running, and confirm side effects with them.

## Cluster access and resources

- `chart-verifier` picks up cluster access from the KUBECONFIG value, or the `--kubeconfig`, `--kube-token`, and related flags can be used instead. Prefer the environment variable or `--kubeconfig` value to avoid emitting contents.
- Do not enable cluster-based checks unless the user explicitly agrees to use a live cluster. Use `--disable chart-testing` for verification without a cluster. This should only be done to shorten feedback loops, as this check is required for most profiles.
- `chart-testing` creates and normally deletes resources in the cluster. `--skip-cleanup` keeps those resources; use it only when requested. Warn that leftover resources may remain if the run is interrupted.

## Signed charts

- `--pgp-public-key` supplies the public key file for `signature-is-valid`. The key file is user-supplied; do not generate or transmit key material.

## Network and file side effects

- `images-are-certified` makes outbound requests to registries and Red Hat certified software catalog APIs. Disable it with `--disable images-are-certified` if those requests are not wanted. This should only be done to shorten feedback loops as this check is required for most profiles.
- `--write-to-file` writes `chartverifier/report.yaml` (or whichever format is required).
- An error log is written to `chartverifier/verifier-<timestamp>.log` by default; suppress with `--suppress-error-log` if the user prefers no file writes.
