# Certification Annotations

## Verifier added annotations

The chart-verifier tool adds annotations to a generated report, for example:

```
 verifier-version: 1.0.0
 chart-uri: /charts
 digest: sha256:28801d8d72d838da1ee05f0809fcc5a2d2b9c6cd27ba3e84c477e76f8916aaa1
 lastCertifiedTime: 2021-04-22 10:49:29.714918174 +0000 UTC m=+1.932279870
 certifiedOpenShiftVersions: "3.7.5"
```

| Annotation                 | Description                                                                                                                                                                                            |
| -------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| verifier-version           | The version of the chart-verifier which generated the report.                                                                                                                                          |
| chart-uri                  | The location of the chart specified to the chart-verifier. For report-only submissions, this must be the public url of the chart.                                                                      |
| digest                     | sha:256 value of the chart from which the report was generated. When submitting a report, this value must match the value generated as part of the submission process.                                 |
| lastCertifiedTime          | The time when the report was generated.                                                                                                                                                                |
| certifiedOpenShiftVersions | The version of OCP that `chart-testing` check was performed on. If the role of the logged-in user prevents this from being accessed, the value must be specified using the `--openshift-version` flag. |

> **_NOTE:_** If the digest in the report does not match the digest of the submitted chart, the submission will fail.

> **_NOTE:_** If the certifiedOpenShiftVersions is not set to a valid OpenShift version, the submission will fail.

## Provider annotations

The chart provider can also include annotations in `Chart.yaml`, which may be used when displaying the chart in the catalog, for example:

```
annotations:
 charts.openshift.io/provider: acme
 charts.openshift.io/name: get your awesomeness here
 charts.openshift.io/supportURL: http://acme-help-is-here.com/
 charts.openshift.io/archs: x86_64
```

| Annotation                     | Description                                                                |
| ------------------------------ | -------------------------------------------------------------------------- |
| charts.openshift.io/provider   | Name of chart provider (e.g., Red Hat), ready to be displayed in UI.       |
| charts.openshift.io/name       | Human readable chart name, ready to be displayed in UI.                    |
| charts.openshift.io/supportURL | Where users can find information about the chart provider's support.       |
| charts.openshift.io/archs      | Comma separated list of supported architectures (e.g., x86_64, s390x, ...) |
