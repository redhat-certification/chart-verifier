# Helm chart checks for Red Hat OpenShift certification

Helm chart checks are a set of checks against which the Red Hat Helm chart-verifier tool verifies and validates whether a Helm chart is qualified for a certification from Red Hat. These checks contain metadata that have certain parameters and values with which a Helm chart must comply to be certified by Red Hat. A Red Hat-certified Helm chart qualifies in terms of readiness for distribution in the [OpenShift Helm Chart Repository](https://github.com/openshift-helm-charts).

## Key features
- You can execute all of the mandatory checks.
  - Include or exclude individual checks during the verification process by using the command line options.
- Each check is independent, and there is no specific execution order.
- Executing the checks includes running the `helm lint`, `helm template`, `helm install`, and `helm test` commands against the chart.
- The checks are configurable. For example, if a chart requires additional values to be compliant with the checks, configure the values using the available options. The options are similar to those used by the `helm lint` and `helm template` commands.
- When there are no error messages, the `helm-lint` check passes the verification and is successful. Messages such as `Warning` and `info` do not cause the check to fail.
- Profiles define the checks needed based on the chart type: partner, redhat or community.
    - Profiles are versioned. Each new version may include updated checks, new checks, new annotations or chnaged annotations.
- The generated report is written to stdout but can optionally be written to a file.
- An error log is created for all verify commands but can be optionally suppressed.

## Types of Helm chart checks
Helm chart checks are categorized into the following types:

#### Table 1: Helm chart check types

| Check type | Description
|---|---
| Mandatory | Checks are required to pass and be successful for certification.
| Recommended | Checks are about to become mandatory; we recommend fixing any check failures.
| Optional | Checks are ready for customer testing. Checks can fail and still pass the verification for certification.
| Experimental | New checks introduced for testing purposes or beta versions.
> **_NOTE:_**  The current release of the chart-verifier includes only the mandatory and optional type of checks.

## Default set of checks for a Helm chart
The following table lists the set of checks for each profile version with details including the name and version of the check, and a description of the check.

#### Table 2: Helm chart default checks

| Profile v1.1 | Profile v1.0 | Description |
|:-------------------------------:|:-------------------------------:|---------------
| [is-helm-v3 v1.0](helm-chart-troubleshooting.md#is-helm-v3-v10) | [is-helm-v3 v1.0](helm-chart-troubleshooting.md#is-helm-v3-v10) | Checks that the given `uri` points to a Helm v3 chart.
| [has-readme v1.0](helm-chart-troubleshooting.md#has-readme-v10) | [has-readme v1.0](helm-chart-troubleshooting.md#has-readme-v10) | Checks that the Helm chart contains the `README.md` file.
| [contains-test V1.0](helm-chart-troubleshooting.md#contains-test-v10) | [contains-test v1.0](helm-chart-troubleshooting.md#contains-test-v10) | Checks that the Helm chart contains at least one test file.
| [has-kubeversion v1.1](helm-chart-troubleshooting.md#has-kubeversion-v11) | [has-kubeversion v1.0](helm-chart-troubleshooting.md#has-kubeversion-v10) | Checks that the `Chart.yaml` file of the Helm chart includes the `kubeVersion` field (v1.0) and is a valid semantic version (v1.1).
| [contains-values-schema v1.0](helm-chart-troubleshooting.md#contains-values-schema-v10) | [contains-values-schema v1.0](helm-chart-troubleshooting.md#contains-values-schema-v10) | Checks that the Helm chart contains a JSON schema file (`values.schema.json`) to validate the `values.yaml` file in the chart.
| [not-contains-crds v1.0](helm-chart-troubleshooting.md#not-contains-crds-v10) | [not-contains-crds v1.0](helm-chart-troubleshooting.md#not-contains-crds-v10) | Checks that the Helm chart does not include custom resource definitions (CRDs).
| [not-contain-csi-objects v1.0](helm-chart-troubleshooting.md#not-contain-csi-objects-v10) | [not-contain-csi-objects v1.0](helm-chart-troubleshooting.md#not-contain-csi-objects-v10) | Checks that the Helm chart does not include Container Storage Interface (CSI) objects.
| [images-are-certified v1.0](helm-chart-troubleshooting.md#images-are-certified-v10) | [images-are-certified v1.0](helm-chart-troubleshooting.md#images-are-certified-v10) | Checks that the images referenced by the Helm chart are Red Hat-certified.
| [helm-lint v1.0](helm-chart-troubleshooting.md#helm-lint-v10) | [helm-lint v1.0](helm-chart-troubleshooting.md#helm-lint-v10) | Checks that the chart is well formed by running the `helm lint` command.
| [chart-testing v1.0](helm-chart-troubleshooting.md#chart-testing-v10) | [chart-testing v1.0](helm-chart-troubleshooting.md#chart-testing-v10)  | Installs the chart and verifies it on a Red Hat OpenShift Container Platform cluster.
| [contains-values v1.0](helm-chart-troubleshooting.md#contains-values-v10)  | [contains-values  v1.0](helm-chart-troubleshooting.md#contains-values-v10) | Checks that the Helm chart contains the `values`[¹](https://github.com/redhat-certification/chart-verifier/blob/main/docs/helm-chart-checks.md#-for-more-information-on-the-values-file-see-values-and-best-practices-for-using-values) file.
| [required-annotations-present v1.0](helm-chart-troubleshooting.md#required-annotations-present-v10) | - | Checks that the Helm chart contains the annotation: ```charts.openshift.io/name```.

#
###### ¹ For more information on the `values` file, see [`values`](https://helm.sh/docs/chart_template_guide/values_files/) and [Best Practices for using values](https://helm.sh/docs/chart_best_practices/values/).

## Run Helm chart checks

There are two ways to run Helm chart checks, either through [containers with `podman`/`docker` command](#using-the-podman-or-docker-command-for-helm-chart-checks), or [run the binary directly (Linux only)](#using-the-chart-verifier-binary-for-helm-chart-checks-linux-only).

### Using the podman or docker command for Helm chart checks
This section provides help on the basic usage of Helm chart checks with the podman or docker command.

#### Prerequisites
- A container engine and the Podman or Docker CLI installed.
- Internet connection to check that the images are Red Hat certified.
- GitHub profile to submit the chart to the [OpenShift Helm Charts Repository](https://github.com/openshift-helm-charts).
- Red Hat OpenShift Container Platform cluster.

#### Procedure

- Run all the available checks for a remotely available chart using a `uri`, assuming the kube config file is available in ${HOME}/.kube:

  ```
  $ podman run --rm -i                                  \
          -e KUBECONFIG=/.kube/config                   \
          -v "${HOME}/.kube":/.kube                     \
          "quay.io/redhat-certification/chart-verifier" \
          verify                                        \
          <chart-uri>
  ```
- Run all the available checks for a chart local to your file system, assuming the chart is in the current directory and the kube config file is available in ${HOME}/.kube:

  ```
  $ podman run --rm                                     \
          -e KUBECONFIG=/.kube/config                   \
          -v "${HOME}/.kube":/.kube                     \
          -v $(pwd):/charts                             \
          "quay.io/redhat-certification/chart-verifier" \
          verify                                        \
          /charts/<chart>
  ```
- Get the list of options for the `verify` command:

  ```
  $ podman run -it --rm quay.io/redhat-certification/chart-verifier verify help
  ```
  The output is similar to the following example:
  ```
  Verifies a Helm chart by checking some of its characteristics

  Usage:
    chart-verifier verify <chart-uri> [flags]

  Flags:
    -S, --chart-set strings           set values for the chart (can specify multiple or separate values with commas: key1=val1,key2=val2)
    -G, --chart-set-file strings      set values from respective files specified via the command line (can specify multiple or separate values with commas: key1=path1,key2=path2)
    -X, --chart-set-string strings    set STRING values for the chart (can specify multiple or separate values with commas: key1=val1,key2=val2)
    -F, --chart-values strings        specify values in a YAML file or a URL (can specify multiple)
        --debug                       enable verbose output
    -x, --disable strings             all checks will be enabled except the informed ones
    -e, --enable strings              only the informed checks will be enabled
    -h, --help                        help for verify
        --kube-apiserver string       the address and the port for the Kubernetes API server
        --kube-as-group stringArray   group to impersonate for the operation, this flag can be repeated to specify multiple groups.
        --kube-as-user string         username to impersonate for the operation
        --kube-ca-file string         the certificate authority file for the Kubernetes API server connection
        --kube-context string         name of the kubeconfig context to use
        --kube-token string           bearer token used for authentication
        --kubeconfig string           path to the kubeconfig file
    -n, --namespace string            namespace scope for this request
    -V, --openshift-version string    set the value of certifiedOpenShiftVersions in the report
    -o, --output string               the output format: default, json or yaml
        --registry-config string      path to the registry config file (default "/home/baiju/.config/helm/registry.json")
        --repository-cache string     path to the file containing cached repository indexes (default "/home/baiju/.cache/helm/repository")
        --repository-config string    path to the file containing repository names and URLs (default "/home/baiju/.config/helm/repositories.yaml")
    -s, --set strings                 overrides a configuration, e.g: dummy.ok=false
    -f, --set-values strings          specify application and check configuration values in a YAML file or a URL (can specify multiple)
    -E, --suppress-error-log          suppress the error log (default: written to ./chartverifier/verifier-<timestamp>.log)
    -w, --write-to-file               write report to ./chartverifier/report.yaml (default: stdout)
  Global Flags:
        --config string   config file (default is $HOME/.chart-verifier.yaml)
  ```
- Run a subset of the checks:

  ```
  $ podman run --rm -i                                  \
          -e KUBECONFIG=/.kube/config                   \
          -v "${HOME}/.kube":/.kube                     \
          "quay.io/redhat-certification/chart-verifier" \
          verify -e images-are-certified,helm-lint      \
          <chart-uri>

  ```
- Run all the checks except a subset:

  ```
  $ podman run --rm -i                                  \
          -e KUBECONFIG=/.kube/config                   \
          -v "${HOME}/.kube":/.kube                     \
          "quay.io/redhat-certification/chart-verifier" \
          verify -x images-are-certified,helm-lint      \
          <chart-uri>
    ```
- Provide chart-override values:

  ```
  $ podman run --rm -i                                  \
          -e KUBECONFIG=/.kube/config                   \
          -v "${HOME}/.kube":/.kube                     \
          "quay.io/redhat-certification/chart-verifier" \
          verify -S default.port=8080                   \
          <chart-uri>
  ```
- Provide chart-override values from a file in the current directory:

  ```
  $ podman run --rm -i                                  \
          -e KUBECONFIG=/.kube/config                   \
          -v "${HOME}/.kube":/.kube                     \
          -v $(pwd):/values                             \
          "quay.io/redhat-certification/chart-verifier" \
          verify -F /values/overrides.yaml              \
          <chart-uri>
  ```

### Saving the report

By default the report is written to stdout which can be redirected to a file. For example:

```
  $ podman run --rm -i                                  \
          -e KUBECONFIG=/.kube/config                   \
          -v "${HOME}/.kube":/.kube                     \
          "quay.io/redhat-certification/chart-verifier" \
          verify -e images-are-certified,helm-lint      \
          <chart-uri> > report.yaml

  ```

Alternatively, use the ```-w```  option to write the report directly to the file ```./chartverifier/report.yaml```. To get this file a volume mount is required to ```/app/chartverifer```. For example:

```
  $ podman run --rm -i                                  \
          -e KUBECONFIG=/.kube/config                   \
          -v "${HOME}/.kube":/.kube                     \
          -v $(pwd)/chartverifier:/app/chartverifier    \
          -w                                            \
          "quay.io/redhat-certification/chart-verifier" \
          verify -e images-are-certified,helm-lint      \
          <chart-uri> > report.yaml

  ```
If the file already exists it is overwritten.

### The error log

By default an error log is written to  file ```./chartverifier/verify-<timestamp>.yaml```. It includes any error messages, the results of each check and additional information around chart testing. To get a copy of the error log a volume mount is required to ```/app/chartverifer```. For example: 

```
  $ podman run --rm -i                                  \
          -e KUBECONFIG=/.kube/config                   \
          -v "${HOME}/.kube":/.kube                     \
          -v $(pwd)/chartverifier:/app/chartverifier    \
          "quay.io/redhat-certification/chart-verifier" \
          verify -e images-are-certified,helm-lint      \
          <chart-uri> > report.yaml

  ```

If multiple logs are added to the same directory a maximum of 10 will be kept. The files will be deleted oldest first.

Use the ```-E``` flag to suppress error log output. 

Note: Error and warning messages are also output to stderr and are not suppressed by the ```-E``` option.


### Using the `chart-verifier` binary for Helm chart checks (Linux only)

Alternatively, download `chart-verifier` binary from the [release page](https://github.com/redhat-certification/chart-verifier/releases), unzip the tarball with `tar zxvf <tarball>`, and run `./chart-verifier verify` under the unzipped directory to perform Helm chart checks. Refer to the [procedures](#procedure) in the podman/docker section, for example,

```
$ podman run --rm -i                                  \
        -e KUBECONFIG=/.kube/config                   \
        -v "${HOME}/.kube":/.kube                     \
        "quay.io/redhat-certification/chart-verifier" \
        verify                                        \
        <chart-uri>
```

will become

```
$ ./chart-verifier verify <chart-uri>
```

By default, `chart-verifier` will assume kubeconfig is under $HOME/.kube, set environment variable KUBECONFIG for different kubeconfig file.

## Profiles

A profile defines a set of checks to run and an indication of whether each check is mandatory or optional. Four profiles are currently available:
- partner
  - Defines the requirements for a partner chart to pass helm chart certfication.
  - All checks are mandatory, that is they must all pass for a partner helm chart to be certified.
- redhat
  - Defines the requirements for a red hat internal chart to pass helm chart certfication.
  - All checks are mandatory, that is they must all pass for a Red Hat helm chart to be certified.
- community
  - Defines the requirements for a community chart to pass helm chart certfication.
  - The ```helm-lint``` check is the only mandatory check with all other checks optional.
- default
  - The default is the same as the partner profile and is used if a specific one is not specified.
  - All checks are mandatory.

Each profile also has a version and currently there are two profile versions: v1.0 and v1.1.

### Profile v1.1

#### Annotations

Annotations added to a v1.1 profile report are common to all profile types: partner, RedHat, community and default

| annotation        | description      |
|-------------------|------------------|
| [digests.chart](helm-chart-annotations.md#digests) | The sha value of the chart as calculated from the copy loaded into memory. |
| [digests.package](helm-chart-annotations.md#digests) | The sha value of the chart tarball if used to create the report. |
| [testedOpenShiftVersion](helm-chart-annotations.md#testedOpenShiftVersion) | The Open Shift version that was used by the chart-testing check. |
| [lastCertifiedTimestamp](helm-chart-annotations.md#lastCertifiedTimestamp) | The time that the report was created by the chart verifier |
| [supportedOpenShiftVersions](helm-chart-annotations.md#supportedOpenShiftVersions) | The Open Shift versions supported by the chart based on the kuberVersion attrinute in chart.yaml |

#### Checks

This table shows which checks are preformed and whether or not they ar mnandatory or optional for each profile type.

| check | partner | RedHat | community | default |
|-------|---------|--------|-----------|---------
| [is-helm-v3 v.1.0](helm-chart-troubleshooting.md#is-helm-v3-v10)  | mandatory | mandatory | optional | mandatory
| [has-readme v1.0](helm-chart-troubleshooting.md#has-readme-v10) | mandatory | mandatory | optional | mandatory
| [contains-test v1.0](helm-chart-troubleshooting.md#contains-test-v10) | mandatory | mandatory | optional | mandatory
| [has-kubeversion v1.1](helm-chart-troubleshooting.md#has-kubeversion-v11)| mandatory | mandatory | optional | mandatory
| [contains-values-schema v1.0](helm-chart-troubleshooting.md#contains-values-schema-v10) | mandatory | mandatory | optional | mandatory
| [not-contains-crds v1.0](helm-chart-troubleshooting.md#not-contains-crds-v10) |  mandatory | mandatory | optional | mandatory
| [not-contain-csi-objects v1.0](helm-chart-troubleshooting.md#not-contain-csi-objects-v10) |  mandatory | mandatory | optional | mandatory
| [images-are-certified v1.0](helm-chart-troubleshooting.md#images-are-certified-v10) | mandatory | mandatory | optional | mandatory
| [helm-lint v1.0](helm-chart-troubleshooting.md#helm-lint-v10) | mandatory | mandatory | optional | mandatory
| [chart-testing v1.0](helm-chart-troubleshooting.md#chart-testing-v10) | mandatory | mandatory | optional | mandatory
| [contains-values v1.0](helm-chart-troubleshooting.md#contains-values-v10)  | mandatory | mandatory | optional | mandatory
| [required-annotations-present v1.0](helm-chart-troubleshooting.md#required-annotations-present-v10) | mandatory | mandatory | optional | mandatory

### Profile 1.0

#### Annotations

Annotations added to a v1.0 profile report are common to all profile types: partner, RedHat, community and default

| annotation        | description      |
|-------------------|------------------|
| [digests.chart](helm-chart-annotations.md#digests) | The sha value of the chart as calculated from the copy loaded into memory. |
| [digests.package](helm-chart-annotations.md#digests) | The sha value of the chart tarball if used to create the report. |
| [certifiedOpenShiftVersion](helm-chart-annotations.md#certifiedOpenShiftVersion) | The Open Shift version that was used by the chart-testing check. |
| [lastCertifiedTimestamp](helm-chart-annotations.md#lastCertifiedTimestamp) | The time that the report was created by the chart verifier |

#### Checks

This table shows which checks are preformed and whether or not they ar mnandatory or potion for each profile type.

| check | partner | RedHat | community | default |
|-------|---------|--------|-----------|---------
| [is-helm-v3 v.1.0](helm-chart-troubleshooting.md#is-helm-v3-v10)  | mandatory | mandatory | optional | mandatory
| [has-readme v1.0](helm-chart-troubleshooting.md#has-readme-v10) | mandatory | mandatory | optional | mandatory
| [contains-test v1.0](helm-chart-troubleshooting.md#contains-test-v10) | mandatory | mandatory | optional | mandatory
| [has-kubeversion v1.0](helm-chart-troubleshooting.md#has-kubeversion-v10)| mandatory | mandatory | optional | mandatory
| [contains-values-schema v1.0](helm-chart-troubleshooting.md#contains-values-schema-v10) | mandatory | mandatory | optional | mandatory
| [not-contains-crds v1.0](helm-chart-troubleshooting.md#not-contains-crds-v10) |  mandatory | mandatory | optional | mandatory
| [not-contain-csi-objects v1.0](helm-chart-troubleshooting.md#not-contain-csi-objects-v10) |  mandatory | mandatory | optional | mandatory
| [images-are-certified v1.0](helm-chart-troubleshooting.md#images-are-certified-v10) | mandatory | mandatory | optional | mandatory
| [helm-lint v1.0](helm-chart-troubleshooting.md#helm-lint-v10) | mandatory | mandatory | optional | mandatory
| [chart-testing v1.0](helm-chart-troubleshooting.md#chart-testing-v10) | mandatory | mandatory | optional | mandatory
| [contains-values v1.0](helm-chart-troubleshooting.md#contains-values-v10)  | mandatory | mandatory | optional | mandatory

#### Running the chart verifier with a specific profile.

To specify which profile to use the --set flag:
```
    --set profile.vendorType=partner
        valid values based on current profiles: partner, community, redhat, default
        default is same as partner.
        If value specified is not specified or not recognized, default will be assumed.
        The flag name is case insensitive.
    --set profile.version=v1.1
        Valid values based on current profiles: v1.0, v1.0
        If value specified is not specified or not recognized, v1.1 will be assumed.
        The flag name is case insensitive.
```
For example:
```
$ podman run --rm -i                                                    \
          -e KUBECONFIG=/.kube/config                                   \
          -v "${HOME}/.kube":/.kube                                     \
          "quay.io/redhat-certification/chart-verifier"                 \
          verify --set profile.vendorType=partner, profile.version=v1.1 \
          <chart-uri>
```

## Chart Testing

### Cluster Config

You can configure the chart-testing check by performing one of the following steps:

* Option 1: Through the `--set` command line option:
    ```text
    $ chart-verifier                                                  \
        verify                                                        \
        --enable chart-testing                                        \
        --set chart-testing.buildId=${BUILD_ID}                       \
        --set chart-testing.upgrade=true                              \
        --set chart-testing.skipMissingValues=true                    \
        --set chart-testing.namespace=${NAMESPACE}                    \
        --set chart-testing.releaseLabel="app.kubernetes.io/instance" \
        --set chart-testing.release=${RELEASE}                        \
        some-chart.tgz
    ```
* Option 2: Create a YAML file (config.yaml) similar to the following example:
   ```text
    chart-testing:
        buildId: <BUILD_ID>
        upgrade: true
        skipMissingValues: true
        namespace: <NAMESPACE>
        releaseLabel: "app.kubernetes.io/instance"
        release: <RELEASE>
    ```

    Specify the file using the `--set-values` command line option:
    ```text
    $ chart-verifier verify --enable chart-testing --set-values config.yaml some-chart.tgz
    ```

    All settings are optional, if not set default values will be used.

### Override values

If the chart requires overrides values, these can be set using through the `--chart-set` command line options:

```
    -S, --chart-set strings           set values for the chart (can specify multiple or separate values with commas: key1=val1,key2=val2)
    -G, --chart-set-file strings      set values from respective files specified via the command line (can specify multiple or separate values with commas: key1=path1,key2=path2)
    -X, --chart-set-string strings    set STRING values for the chart (can specify multiple or separate values with commas: key1=val1,key2=val2)
    -F, --chart-values strings        specify values in a YAML file or a URL (can specify multiple)
```

### Check processing

The `chart-testing` check performs the following actions, keeping the semantics provided by [github.com/helm/chart-testing](https://github.com/helm/chart-testing):
1. Install: the chart being verified will be installed in the available OpenShift cluster utilizing the same semantics client-go uses to find the current context:
    1. `--kubeconfig flag`
    1. `KUBECONFIG` environment variable
    1. `$HOME/.kube/config`.
1. Test: once a release is installed for the chart being verified, performs the same actions as helm test would, which installing all chart resources containing the "helm.sh/hook": test annotation.

The check will be considered successful when the chart's installation and tests are all successful.
