# chart-verifier

`chart-verifier` is a tool that certifies a Helm chart against a configurable list of checks; those checks can be
whitelisted or blacklisted through command line options.

Each check is independent and order is not and will not be guaranteed, and its input will be informed through options in
the command line interface; currently the only input is the required `uri` option.

The following checks are being implemented:

| Name | Description
|---|---
| `is-helm-v3` | Checks whether the given `uri` is a Helm v3 chart.
| `has-readme` | Checks whether the Helm chart contains a `README.md` file.
| `contains-test` | Checks whether the Helm chart contains at least one test file.
| `readme-contains-values-schema` | Checks whether the Helm chart `README.md` file contains a `values` schema section.
| `keywords-are-openshift-categories` | Checks whether the Helm chart's `Chart.yaml` file includes keywords mapped to OpenShift categories.
| `is-commercial-chart` | Checks whether the Helm chart is a Commercial chart.
| `is-community-chart` | Checks whether the Helm chart is a Community chart.
| `has-minkubeversion` | Checks whether the Helm chart's `Chart.yaml` includes the `minKubeVersion` field. 
| `not-contains-crds` | Check whether the Helm chart does not include CRDs.
| `not-contains-infra-plugins-and-drivers` | Check whether the Helm chart does not include infra plugins and drivers (network, storage, hardware, etc)
| `can-be-installed-without-manual-prerequisites` |
| `can-be-installed-without-cluster-admin-privileges` |

## Architecture

This tool is part of a larger process that aims to certify Helm charts, and its sole responsibility is to ingest a Helm
chart URI (`file://`, `https?://`, etc)
and return either a *positive* result indicating the Helm chart has passed all checks, or a *negative* result indicating
which checks have failed and possibly propose solutions.

The application is separated in two pieces: a command line interface and a library. This is handy because the command
line interface is specific to the user interface, and the library can be generic enough to be used to, for example,
inspect Helm chart bytes in flight.

One positive aspect of the command line interface specificity is that its output can be tailored to the methods of
consumption the user expects; in other words, the command line interface can be programmed in such way it can be
represented as either *YAML* or *JSON* formats, in addition to a descriptive representation tailored to human actors.

The interpretation of what is considered a certified Helm chart depends on which checks the chart has been submitted to,
so this information should be present in the certificate as well.

Primitive functions to manipulate the Helm chart should be provided, since most checks involve inspecting the contents
of the chart itself; for example, whether a `README.md` file exists, or whether `README.md` contains the `values`'
specification, implicating in offering a cache API layer is required to avoid downloading and unpacking the charts for
each test.

## Usage

To certify a chart against all available checks:

```text
> chart-verifier --uri ./chart.tgz
> chart-verifier --uri ~/src/chart
> chart-verifier --uri https://www.example.com/chart.tgz
```

To apply only the `is-helm-v3` check:

```text
> chart-verifier --only is-helm-v3 --uri https://www.example.com/chart.tgz
```

To apply all checks except `is-helm-v3`:

```text
> chart-verifier --except is-helm-v3 --uri https://www.example.com/chart.tgz
```
