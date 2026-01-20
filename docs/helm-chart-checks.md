# Helm chart checks for Red Hat OpenShift certification

Helm chart checks are a set of checks against which the Red Hat Helm chart-verifier tool verifies and validates whether a Helm chart is qualified for a certification from Red Hat. These checks contain metadata that have certain parameters and values with which a Helm chart must comply to be certified by Red Hat. A Red Hat-certified Helm chart qualifies in terms of readiness for distribution in the [OpenShift Helm Chart Repository](https://github.com/openshift-helm-charts).

- [Helm chart checks for Red Hat OpenShift certification](#helm-chart-checks-for-red-hat-openshift-certification)
  - [Key features](#key-features)
  - [Types of Helm chart checks](#types-of-helm-chart-checks)
      - [Table 1: Helm chart check types](#table-1-helm-chart-check-types)
  - [Default set of checks for a Helm chart](#default-set-of-checks-for-a-helm-chart)
      - [Table 2: Helm chart default checks](#table-2-helm-chart-default-checks)
- [](#)
          - [¹ For more information on the `values` file, see `values` and Best Practices for using values.](#-for-more-information-on-the-values-file-see-values-and-best-practices-for-using-values)
  - [Run Helm chart checks](#run-helm-chart-checks)
    - [Using the podman or docker command for Helm chart checks](#using-the-podman-or-docker-command-for-helm-chart-checks)
      - [Prerequisites](#prerequisites)
      - [Procedure](#procedure)
  - [Signed charts](#signed-charts)

## Key features
- You can execute all of the mandatory checks.
  - Include or exclude individual checks during the verification process by using the command line options.
- Each check is independent, and there is no specific execution order.
- Executing the checks includes running the `helm lint`, `helm template`, `helm install`, and `helm test` commands against the chart.
- The checks are configurable. For example, if a chart requires additional values to be compliant with the checks, configure the values using the available options. The options are similar to those used by the `helm lint` and `helm template` commands.
- When there are no error messages, the `helm-lint` check passes the verification and is successful. Messages such as `Warning` and `info` do not cause the check to fail.
- Profiles define the checks needed based on the chart type: partner, redhat or community.
  - Profiles are versioned. Each new version can include updated checks, new checks, new annotations, or changed annotations.
- The generated report, by default, is written to stdout.
  - Alternatively the ```--write-to-file``` flag can be used to write to a ```report.yaml``` file.
- An error log is created for all verify commands but can be optionally suppressed.
- You can indicate that a chart is not to be published in the OpenShift catalog.
- From chart verifier version 1.9.0 the generated report includes a sha value based on the report content. This is used during the submission process to verify the integrity of the report.
- You can verify a signed chart. See: [Signed Charts](#signed-charts).

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

| Profile v1.4| Profile v1.3 | Profile v1.2 | Profile v1.1 | Profile v1.0 | Description |
|---|---|---|---|---|---|
| [is-helm-v3 v1.0](helm-chart-troubleshooting.md#is-helm-v3-v10) | [is-helm-v3 v1.0](helm-chart-troubleshooting.md#is-helm-v3-v10) | [is-helm-v3 v1.0](helm-chart-troubleshooting.md#is-helm-v3-v10) | [is-helm-v3 v1.0](helm-chart-troubleshooting.md#is-helm-v3-v10) | [is-helm-v3 v1.0](helm-chart-troubleshooting.md#is-helm-v3-v10) | Checks that the given `uri` points to a Helm v3 chart. |
| [has-readme v1.0](helm-chart-troubleshooting.md#has-readme-v10) | [has-readme v1.0](helm-chart-troubleshooting.md#has-readme-v10) | [has-readme v1.0](helm-chart-troubleshooting.md#has-readme-v10) | [has-readme v1.0](helm-chart-troubleshooting.md#has-readme-v10) | [has-readme v1.0](helm-chart-troubleshooting.md#has-readme-v10) | Checks that the Helm chart contains the `README.md` file. |
| [contains-test V1.0](helm-chart-troubleshooting.md#contains-test-v10) | [contains-test V1.0](helm-chart-troubleshooting.md#contains-test-v10) | [contains-test V1.0](helm-chart-troubleshooting.md#contains-test-v10) | [contains-test V1.0](helm-chart-troubleshooting.md#contains-test-v10) | [contains-test v1.0](helm-chart-troubleshooting.md#contains-test-v10) | Checks that the Helm chart contains at least one test file. |
| [has-kubeversion v1.1](helm-chart-troubleshooting.md#has-kubeversion-v11) | [has-kubeversion v1.1](helm-chart-troubleshooting.md#has-kubeversion-v11) | [has-kubeversion v1.1](helm-chart-troubleshooting.md#has-kubeversion-v11) | [has-kubeversion v1.1](helm-chart-troubleshooting.md#has-kubeversion-v11) | [has-kubeversion v1.0](helm-chart-troubleshooting.md#has-kubeversion-v10) | Checks that the `Chart.yaml` file of the Helm chart includes the `kubeVersion` field (v1.0) and is a valid semantic version (v1.1). |
| [contains-values-schema v1.0](helm-chart-troubleshooting.md#contains-values-schema-v10) | [contains-values-schema v1.0](helm-chart-troubleshooting.md#contains-values-schema-v10) | [contains-values-schema v1.0](helm-chart-troubleshooting.md#contains-values-schema-v10) | [contains-values-schema v1.0](helm-chart-troubleshooting.md#contains-values-schema-v10) | [contains-values-schema v1.0](helm-chart-troubleshooting.md#contains-values-schema-v10) | Checks that the Helm chart contains a JSON schema file (`values.schema.json`) to validate the `values.yaml` file in the chart. |
| [not-contains-crds v1.0](helm-chart-troubleshooting.md#not-contains-crds-v10) | [not-contains-crds v1.0](helm-chart-troubleshooting.md#not-contains-crds-v10) | [not-contains-crds v1.0](helm-chart-troubleshooting.md#not-contains-crds-v10) | [not-contains-crds v1.0](helm-chart-troubleshooting.md#not-contains-crds-v10) | [not-contains-crds v1.0](helm-chart-troubleshooting.md#not-contains-crds-v10) | Checks that the Helm chart does not include custom resource definitions (CRDs). |
| [not-contain-csi-objects v1.0](helm-chart-troubleshooting.md#not-contain-csi-objects-v10) | [not-contain-csi-objects v1.0](helm-chart-troubleshooting.md#not-contain-csi-objects-v10) | [not-contain-csi-objects v1.0](helm-chart-troubleshooting.md#not-contain-csi-objects-v10) | [not-contain-csi-objects v1.0](helm-chart-troubleshooting.md#not-contain-csi-objects-v10) | [not-contain-csi-objects v1.0](helm-chart-troubleshooting.md#not-contain-csi-objects-v10) | Checks that the Helm chart does not include Container Storage Interface (CSI) objects. |
| [images-are-certified v1.1](helm-chart-troubleshooting.md#images-are-certified-v10) | [images-are-certified v1.1](helm-chart-troubleshooting.md#images-are-certified-v10) | [images-are-certified v1.1](helm-chart-troubleshooting.md#images-are-certified-v10) | [images-are-certified v1.0](helm-chart-troubleshooting.md#images-are-certified-v10) | [images-are-certified v1.0](helm-chart-troubleshooting.md#images-are-certified-v10) | Checks that the images referenced by the Helm chart are Red Hat-certified. |
| [helm-lint v1.0](helm-chart-troubleshooting.md#helm-lint-v10) | [helm-lint v1.0](helm-chart-troubleshooting.md#helm-lint-v10) | [helm-lint v1.0](helm-chart-troubleshooting.md#helm-lint-v10) | [helm-lint v1.0](helm-chart-troubleshooting.md#helm-lint-v10) | [helm-lint v1.0](helm-chart-troubleshooting.md#helm-lint-v10) | Checks that the chart is well formed by running the `helm lint` command. |
| [chart-testing v1.0](helm-chart-troubleshooting.md#chart-testing-v10) | [chart-testing v1.0](helm-chart-troubleshooting.md#chart-testing-v10) | [chart-testing v1.0](helm-chart-troubleshooting.md#chart-testing-v10) | [chart-testing v1.0](helm-chart-troubleshooting.md#chart-testing-v10) | [chart-testing v1.0](helm-chart-troubleshooting.md#chart-testing-v10) | Installs the chart and verifies it on a Red Hat OpenShift Container Platform cluster. |
| [contains-values v1.0](helm-chart-troubleshooting.md#contains-values-v10) | [contains-values v1.0](helm-chart-troubleshooting.md#contains-values-v10) | [contains-values v1.0](helm-chart-troubleshooting.md#contains-values-v10) | [contains-values v1.0](helm-chart-troubleshooting.md#contains-values-v10) | [contains-values  v1.0](helm-chart-troubleshooting.md#contains-values-v10) | Checks that the Helm chart contains the `values`[¹](https://github.com/redhat-certification/chart-verifier/blob/main/docs/helm-chart-checks.md#-for-more-information-on-the-values-file-see-values-and-best-practices-for-using-values) file. |
| [required-annotations-present v1.0](helm-chart-troubleshooting.md#required-annotations-present-v10) | [required-annotations-present v1.0](helm-chart-troubleshooting.md#required-annotations-present-v10) | [required-annotations-present v1.0](helm-chart-troubleshooting.md#required-annotations-present-v10) | [required-annotations-present v1.0](helm-chart-troubleshooting.md#required-annotations-present-v10) | - | Checks that the Helm chart contains the annotation: ```charts.openshift.io/name```. |
| [signature-is-valid v1.0](helm-chart-troubleshooting.md#signature-is-valid-v10) | [signature-is-valid v1.0](helm-chart-troubleshooting.md#signature-is-valid-v10) | [signature-is-valid v1.0](helm-chart-troubleshooting.md#signature-is-valid-v10) | - | - | Verifies a signed chart based on a provided public key. |
| [has-notes v1.0](helm-chart-troubleshooting.md#has-notes-v10) | [has-notes v1.0](helm-chart-troubleshooting.md#has-notes-v10) | - | - | - | Checks that the Helm chart contains the `NOTES.txt` file in the templates directory. |
| [cluster-is-not-eol v1.0](helm-chart-troubleshooting.md#cluster-is-not-eol-v10) | - | - | - | - | Checks that Helm chart was tested on a non EOL cluster. |
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
          -e KUBECONFIG=/.kube/config:z                 \
          -v "${HOME}/.kube":/.kube:z                   \
          "quay.io/redhat-certification/chart-verifier" \
          verify                                        \
          <chart-uri>
  ```
- Run all the available checks for a chart local to your file system, assuming the chart is in the current directory and the kube config file is available in ${HOME}/.kube:

  ```
  $ podman run --rm                                     \
          -e KUBECONFIG=/.kube/config                   \
          -v "${HOME}/.kube":/.kube:z                   \
          -v $(pwd):/charts:z                           \
          "quay.io/redhat-certification/chart-verifier" \
          verify                                        \
          /charts/<chart>
  ```
- Get the list of options for the `verify` command:

  ```
  $ podman run -it --rm quay.io/redhat-certification/chart-verifier verify --help
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
          --helm-install-timeout duration   helm install timeout (default 5m0s)
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
      -k, --pgp-public-key string       file containing gpg public key of the key used to sign the chart
      -W, --web-catalog-only            set this to indicate that the distribution method is web catalog only (default: false)
          --registry-config string      path to the registry config file (default "/home/baiju/.config/helm/registry.json")
          --repository-cache string     path to the file containing cached repository indexes (default "/home/baiju/.cache/helm/repository")
          --repository-config string    path to the file containing repository names and URLs (default "/home/baiju/.config/helm/repositories.yaml")
      -s, --set strings                 overrides a configuration, e.g: dummy.ok=false
      -f, --set-values strings          specify application and check configuration values in a YAML file or a URL (can specify multiple)
      -E, --suppress-error-log          suppress the error log (default: written to ./chartverifier/verifier-<timestamp>.log)
          --timeout duration            time to wait for completion of chart install and test (default 30m0s)
          --write-junitxml-to string    If set, will write a junitXML representation of the result to the specified path in addition to the configured output format
      -w, --write-to-file               write report to ./chartverifier/report.yaml (default: stdout)
    Global Flags:
          --config string   config file (default is $HOME/.chart-verifier.yaml)
  ```
- Run a subset of the checks:

  ```
  $ podman run --rm -i                                  \
          -e KUBECONFIG=/.kube/config                   \
          -v "${HOME}/.kube":/.kube:z                   \
          "quay.io/redhat-certification/chart-verifier" \
          verify -e images-are-certified,helm-lint      \
          <chart-uri>

  ```
- Run all the checks except a subset:

  ```
  $ podman run --rm -i                                  \
          -e KUBECONFIG=/.kube/config                   \
          -v "${HOME}/.kube":/.kube:z                   \
          "quay.io/redhat-certification/chart-verifier" \
          verify -x images-are-certified,helm-lint      \
          <chart-uri>
    ```
- Provide chart-override values:

  ```
  $ podman run --rm -i                                  \
          -e KUBECONFIG=/.kube/config                   \
          -v "${HOME}/.kube":/.kube:z                   \
          "quay.io/redhat-certification/chart-verifier" \
          verify -S default.port=8080                   \
          <chart-uri>
  ```
- Provide chart-override values from a file in the current directory:

  ```
  $ podman run --rm -i                                  \
          -e KUBECONFIG=/.kube/config                   \
          -v "${HOME}/.kube":/.kube:z                   \
          -v $(pwd):/values:z                           \
          "quay.io/redhat-certification/chart-verifier" \
          verify -F /values/overrides.yaml              \
          <chart-uri>
  ```

### Timeout Option

Increase the timeout value if chart-testing is going to take more time, default value is 30m.

  ```
  $ podman run --rm -i                                  \
          -e KUBECONFIG=/.kube/config                   \
          -v "${HOME}/.kube":/.kube:z                   \
          -v $(pwd):/values:z                           \
          "quay.io/redhat-certification/chart-verifier" \
          verify --timeout 40m                          \
          <chart-uri>
  ```
Note: In case chart-testing takes more time, it is advised to submit the report for certification since the certification process will use the default value of 30m.

### Saving the report

By default the report is written to stdout which can be redirected to a file. For example:

```
  $ podman run --rm -i                                  \
          -e KUBECONFIG=/.kube/config                   \
          -v "${HOME}/.kube":/.kube:z                   \
          "quay.io/redhat-certification/chart-verifier" \
          verify -e images-are-certified,helm-lint      \
          <chart-uri> > report.yaml

  ```

Alternatively, use the ```-w```  option to write the report directly to the file ```./chartverifier/report.yaml```. To get this file a volume mount is required to ```/app/chartverifer```. For example:

```
  $ podman run --rm -i                                  \
          -e KUBECONFIG=/.kube/config                   \
          -v "${HOME}/.kube":/.kube:z                   \
          -v $(pwd)/chartverifier:/app/chartverifier:z  \
          -w                                            \
          "quay.io/redhat-certification/chart-verifier" \
          verify -e images-are-certified,helm-lint      \
          <chart-uri>

  ```
If the file already exists it is overwritten.

An additional report can be written in JUnit XML format if requested with the
`--write-junitxml-to` flag, passing in the desired output filename. 

```
  $ podman run --rm -i                                         \
          -e KUBECONFIG=/.kube/config                          \
          -v "${HOME}/.kube":/.kube:z                          \
          -v $(pwd)/chartverifier:/app/chartverifier:z         \
          -w                                                   \
          "quay.io/redhat-certification/chart-verifier"        \
          verify                                               \
          --write-junitxml-to /app/chartverifier/report-junit.xml \
          <chart-uri>
```

JUnitXML is not an additional report format that can be used for certification
or validation using chart-verifier, and is only intended to be consumed by user
tooling. The YAML or JSON report is always written as specified.

### The error log

By default an error log is written to  file ```./chartverifier/verify-<timestamp>.yaml```. It includes any error messages, the results of each check and additional information around chart testing. To get a copy of the error log a volume mount is required to ```/app/chartverifer```. For example:

```
  $ podman run --rm -i                                  \
          -e KUBECONFIG=/.kube/config                   \
          -v "${HOME}/.kube":/.kube:z                   \
          -v $(pwd)/chartverifier:/app/chartverifier:z  \
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
        -v "${HOME}/.kube":/.kube:z                   \
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

A profile defines a set of checks to run and an indication of whether each check is mandatory or optional. Following profiles are currently available:
- partner
  - Defines the requirements for a partner chart to pass helm chart certfication.
  - All checks are mandatory, that is they must all pass for a partner helm chart to be certified.
- redhat
  - Defines the requirements for a red hat internal chart to pass helm chart certfication.
  - All checks are mandatory, that is they must all pass for a Red Hat helm chart to be certified.
- community
  - Defines the requirements for a community chart to pass helm chart certfication.
  - The ```helm-lint``` check is the only mandatory check with all other checks optional.
- developer-console
  - Defines the requirements for a developer-console chart to be validated.
  - The checks which are enabled for this profile are ```helm-lint```,```is-helm-v3```,```contains-values```,```contains-values-schema```,```has-kubeversion``` and ```has-readme```. All these checks are mandatory.
- default
  - The default is the same as the partner profile and is used if a specific one is not specified.
  - All checks are mandatory.

Each profile also has a version and currently there are five profile versions: v1.0, v1.1, v1.2, v1.3 and v1.4. The `developer-console` just has one profile version v1.0.

### Profile v1.4

Compared to profile v1.3, adds a new check:

| check                                                                           | partner | RedHat | community | default |
|---------------------------------------------------------------------------------|---------|--------|-----------|---------
| [cluster-is-not-eol v1.0](helm-chart-troubleshooting.md#cluster-is-not-eol-v10) | optional | optional | optional | optional

### Profile v1.3

Compared to profile v1.2, adds a new check:

| check | partner | RedHat | community | default |
|-------|---------|--------|-----------|---------
| [has-notes v1.0](helm-chart-troubleshooting.md#has-notes-v10) | optional | optional | optional | optional

### Profile v1.2

Compared to profile v1.1, adds a new check:

| check | partner | RedHat | community | default |
|-------|---------|--------|-----------|---------
| [signature-is-valid v1.0](helm-chart-troubleshooting.md#signature-is-valid-v10) | mandatory | mandatory | optional | mandatory
| [images-are-certified v1.1](helm-chart-troubleshooting.md#images-are-certified-v11) | mandatory | mandatory | optional | mandatory

### Profile v1.1

#### Annotations

Annotations added to a v1.1 profile report are common to all profile types: partner, RedHat, community and default

| annotation        | description      |
|-------------------|------------------|
| [digests.chart](helm-chart-annotations.md#digests) | The sha value of the chart as calculated from the copy loaded into memory. |
| [digests.package](helm-chart-annotations.md#digests) | The sha value of the chart tarball if used to create the report. |
| [testedOpenShiftVersion](helm-chart-annotations.md#testedOpenShiftVersion) | The OpenShift version that was used by the chart-testing check. |
| [lastCertifiedTimestamp](helm-chart-annotations.md#lastCertifiedTimestamp) | The time that the report was created by the chart verifier |
| [supportedOpenShiftVersions](helm-chart-annotations.md#supportedOpenShiftVersions) | The OpenShift versions supported by the chart based on the kuberVersion attrinute in chart.yaml |

#### Checks

This table shows which checks are preformed and whether or not they are mandatory or optional for each profile type.

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
| [certifiedOpenShiftVersion](helm-chart-annotations.md#certifiedOpenShiftVersion) | The OpenShift version that was used by the chart-testing check. |
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

#### Running the chart verifier with a specific profile

To specify which profile to use the --set flag:
```
    --set profile.vendorType=partner
        valid values based on current profiles: partner, community, redhat, default
        default is same as partner.
        If value is not specified or not recognized, default will be assumed.
        The flag name is case insensitive.
    --set profile.version=v1.1
        Valid values based on current profiles: v1.0, v1.0
        If value is not specified or not recognized, v1.1 will be assumed.
        The flag name is case insensitive.
```
For example:
```
$ podman run --rm -i                                                    \
          -e KUBECONFIG=/.kube/config                                   \
          -v "${HOME}/.kube":/.kube:z                                   \
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
1. Install: the chart being verified will be installed in the available OpenShift cluster using the same semantics client-go uses to find the current context:
    1. `--kubeconfig flag`
    1. `KUBECONFIG` environment variable
    1. `$HOME/.kube/config`.
1. Test: once a release is installed for the chart being verified, performs the same actions as helm test would, which installing all chart resources containing the "helm.sh/hook": test annotation.

The check will be considered successful when the chart's installation and tests are all successful.

### Chart testing timeouts

For the chart install and test check there are two configurable timeout options:

- ```--helm-install-timeout```
    - limits how long the `chart-testing` check waits for chart install to complete.
    - default is 5 minutes
    - set for example:
      ```--helm-install-timeout  10m0s```
- ```--timeout```
    - limits how long the check waits for the `chart-testing` check to complete: 
        1. chart install
        2. wait for deployments to be available
        3. run the test
    - default is 30 minutes
    - set for example:
        ```--timeout  60m0s```
      
Notes: 
- The timeouts are independent. 
  - Changing one timeout does not impact the other.
  - If ```helm-install-timeout``` is increased, consider also increasing ```timeout```
- The [helm chart certification process](./helm-chart-submission.md#submission-of-helm-charts-for-red-hat-openShift-certification) uses default timeout values.
  - If a helm chart can only pass the chart testing check with modified timeouts a verifier report must be included in the chart submission.  

## Images are Certified

The **images-are-certified** check must be able to render your templates in
order to extract image references. In that process, your chart's kubeVersion
constraint will be evaluated against a "server" version. No server is required
for this check, and Chart Verifier will mock out the server version information
for you (as is also done via `helm template`).

However, you may find that Chart Verifier mocks out server version that does not
fall within your chart's constraints. In most cases, this will happen if you
have an upper bound on your supported kubeVersion (e.g. `<= v1.30.0`).

If this is the case, you can override the mocked version
string to another valid semantic version that falls within your constraints.
This is done using the `--set images-are-certified.kube-version=<yourversion>`
flag.

For example, to set the kube-version value for this check to `v1.29`:

```shell
    $ chart-verifier                                                  \
        verify                                                        \
        --enable images-are-certified                                 \
        --set images-are-certified.kube-version=v1.29                 \
        some-chart.tgz
```

## Signed charts

In profile v1.2 a new mandatory check is added for signed charts. For information on signed charts see [helm provenance and integrity](https://helm.sh/docs/topics/provenance/).
- For a signed chart:
  - The check requires a pgp public key file to run.
    - Ensures a signed chart is validly signed for the public key which will be provided to users to verify the chart.
    - Specify the public key file using the flag: ```--pgp-public-key <public-key-file>```
    - The check runs ```helm verify``` using the public key
       - If ```helm verify``` fails the check will fail.
       - For information on ```helm verify``` see [helm verify](https://helm.sh/docs/helm/helm_verify) 
    - To create the pgp public key file:
      - run: ```gpg --export -a <User-Name> > <public-key-file>```
        - User-Name is the user name of the secret key used to sign the chart.
  - If a pgp public key is not provided the check result will be "SKIPPED" which is considered a PASS for chart certification purposes.
- For a non-signed chart:
  - the check result will be "SKIPPED" which is considered a PASS for chart certification purposes.
    
For troubleshooting this check see: [signature-is-valid v1.0](helm-chart-troubleshooting.md#signature-is-valid-v10).
    