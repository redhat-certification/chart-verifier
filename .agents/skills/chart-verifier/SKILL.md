---
name: chart-verifier
description: Validate Helm charts and inspect chart verification reports using the chart-verifier CLI. Use when the user asks to verify, certify, or run checks against a Helm chart, or to read, summarize, or diagnose a chart-verifier report.
license: Apache-2.0
compatibility: See references/install-chart-verifier.md for installation prerequisites. Some checks additionally need network, container registry, or Kubernetes cluster access.
---

# Chart Verifier

Use the `chart-verifier` CLI to verify Helm charts (local paths, remote URLs) and to inspect the verification reports it produces.

## Safety rules

- Never print, copy, or transmit kubeconfig contents, bearer tokens, registry credentials, or private keys. Pass credential material only via the CLI's dedicated flags or files the user already controls. Describe credential errors without exposing credential values.
- Do NOT install with root privileges.
- Do NOT read a KUBECONFIG.

## Procedure

In this skill and its references, `<chart-verifier-bin>` is a placeholder for the discovered or installed executable's path. Replace it with the actual path before running commands, quoting the path if it contains spaces.

### 0. Confirm environmental requirements

Use [check-kubeconfig.sh](scripts/check-kubeconfig.sh) to check for the presence of the kubeconfig file. If it is not present, you must STOP.
If KUBECONFIG is not found, you must tell the user to set this value and configure with sufficient access to install their specified chart.

### 1. Check for an existing installation

If asked to install the binary, skip this step.
Otherwise, use [find-chart-verifier.sh](scripts/find-chart-verifier.sh) to locate an existing binary on `PATH`:

### 2. Download and install Chart Verifier

Follow [references/install-chart-verifier.md](references/install-chart-verifier.md) for scripts and guidance on installing Chart Verifier.

### 3. Confirm the installation

Run `<chart-verifier-bin> version` and `<chart-verifier-bin> --help` (and the `--help` of the subcommands you will use), using the discovered or installed binary. Prefer the CLI's own `--help` output over assuming flags from any particular release.

### 4. Verify a chart

Follow [references/verify.md](references/verify.md) for the happy path from profile selection through verification and report validation.

### 5. Inspect a report

Use `<chart-verifier-bin> report` to summarize report contents. Subcommands, digest validation, and profile options are covered in [references/report-interpretation.md](references/report-interpretation.md).

### 6. Interpret the outcome

Present each check result and what to do next.
Outcomes - positive, negative, skipped, and digest mismatch - are explained with suggested next steps in [references/report-interpretation.md](references/report-interpretation.md).

### 7. Submit for certification

Once the chart passes its profile, submit the same chart source and `chartverifier/report.yaml` in a pull request to [openshift-helm-charts/charts](https://github.com/openshift-helm-charts/charts), following its [submission instructions](https://github.com/openshift-helm-charts/charts/blob/main/docs/README.md).

## Extended use cases

For information on included checks and their environmental prerequisites, see [references/checks.md](references/checks.md).

For extended use cases and additional flags, see [references/flags-and-side-effects.md](references/flags-and-side-effects.md).
