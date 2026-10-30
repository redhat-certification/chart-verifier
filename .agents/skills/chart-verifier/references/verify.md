# Verify a chart

This document covers the happy path for verifying a chart.
Examples use `<chart-verifier-bin>` as a placeholder for the executable path selected in [the skill procedure](../SKILL.md#procedure).

## 1. Verify the chart

```sh
<chart-verifier-bin> verify --set profile.vendorType=partner --write-to-file ./my-chart
```

For a remote chart, replace `./my-chart` with its download URL, e.g. `"https://github.com/<org>/<repo>/blob/<branch>/<path>/<chart-name>-<chart-version>.tgz?raw=true"`.

Ask the user which profile to use (`partner`, `redhat`, or `community`); if unspecified, assume `partner`.
If needed, add `--chart-values values.yaml`; repeat the flag for additional values files.

For partners, ask whether they distribute through Red Hat's index or distribute the chart themselves. If unsure, have them check their project at [Red Hat Connect](https://connect.redhat.com). For self-distribution configured as **web catalog only**, add `--web-catalog-only` and verify a packaged chart (`.tgz`, local or remote), not a source directory. Otherwise, including when unsure, omit the flag.

## 2. Confirm the profile passes

```sh
<chart-verifier-bin> report all chartverifier/report.yaml
```

This evaluates the report against its recorded profile. Proceed only when `results.failed` is `"0"`; a successful command exit alone does not mean the chart passed. Fix any reported issues and rerun verification.
