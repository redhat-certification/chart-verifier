# Helm chart checks for Red Hat OpenShift certification

Helm chart checks are a set of checks against which the Red Hat Helm chart-verifier tool verifies and validates whether a Helm chart is qualified for a certification from Red Hat. These checks contain metadata that have certain parameters and values with which a Helm chart must comply to be certified by Red Hat. A Red Hat-certified Helm chart qualifies in terms of readiness for distribution in the [OpenShift Helm Chart Repository](https://github.com/openshift-helm-charts).

## Key features
- You can execute all of the mandatory checks.
  - Include or exclude individual checks during the verification process by using the command line options.
- Each check is independent, and there is no specific execution order.
- Executing the checks includes running the `helm lint`, `helm template`, `helm install`, and `helm test` commands against the chart.
- The checks are configurable. For example, if a chart requires additional values to be compliant with the checks, configure the values using the available options. The options are similar to those used by the `helm lint` and `helm template` commands.
- When there are no error messages, the `helm-lint` check passes the verification and is successful. Messages such as `Warning` and `info` do not cause the check to fail.

## Types of Helm chart checks
Helm chart checks are categorized into the following types:

#### Table 1: Helm chart check types

| Check type | Description
|---|---
| Mandatory | Checks are required to pass and be successful for certification.
| Recommended | Checks are about to become mandatory; we recommend fixing any check failures.
| Optional | Checks are ready for customer testing. Checks can fail and still pass the verification for certification.
| Experimental | New checks introduced for testing purposes or beta versions.
> **_NOTE:_**  The current release of the chart-verifier includes only the mandatory type of checks.

## Default set of checks for a Helm chart
The following table lists the default set of checks with details including the names of the checks, their type, and description. These checks must be successful for a Helm chart to be certified by Red Hat.

#### Table 2: Helm chart default checks

| Name | Check type | Description
|------|------------|------------
| `is-helm-v3` | Mandatory | Checks that the given `uri` points to a Helm v3 chart.
| `has-readme` | Mandatory | Checks that the Helm chart contains the `README.md` file.
| `contains-test` | Mandatory | Checks that the Helm chart contains at least one test file.
| `has-kubeversion` | Mandatory | Checks that the `Chart.yaml` file of the Helm chart includes the `kubeVersion` field.
| `contains-values-schema` | Mandatory | Checks that the Helm chart contains a JSON schema file (`values.schema.json`) to validate the `values.yaml` file in the chart.
| `not-contains-crds` | Mandatory | Checks that the Helm chart does not include custom resource definitions (CRDs).
| `not-contain-csi-objects` | Mandatory | Checks that the Helm chart does not include Container Storage Interface (CSI) objects.
| `images-are-certified` | Mandatory | Checks that the images referenced by the Helm chart are Red Hat-certified.
| `helm-lint` | Mandatory | Checks that the chart is well formed by running the `helm lint` command.
| `chart-testing` | Mandatory | Installs the chart and verifies it on a Red Hat OpenShift Container Platform cluster.
| `contains-values` | Mandatory | Checks that the Helm chart contains the `values`[¹](https://github.com/redhat-certification/chart-verifier/blob/main/docs/helm-chart-checks.md#-for-more-information-on-the-values-file-see-values-and-best-practices-for-using-values) file.

#
###### ¹ For more information on the `values` file, see [`values`](https://helm.sh/docs/chart_template_guide/values_files/) and [Best Practices for using values](https://helm.sh/docs/chart_best_practices/values/).

## Using the podman or docker command for Helm chart checks
This section provides help on the basic usage of Helm chart checks with the podman or docker command.

### Prerequisites
- A container engine and the Podman or Docker CLI installed.
- Internet connection to check that the images are Red Hat certified.
- GitHub profile to submit the chart to the [OpenShift Helm Charts Repository](https://github.com/openshift-helm-charts).
- Red Hat OpenShift Container Platform cluster.

### Procedure

- Run all the available checks for the chart using a `uri`:

  ```
  $ podman run -it --rm quay.io/redhat-certification/chart-verifier verify <chart-uri>
  ```
- Run all the checks available locally on your system for the chart, from the same directory as the chart:

  ```
  $ podman run -v $(pwd):/charts --rm quay.io/redhat-certification/chart-verifier verify /charts/<chart>
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
    -l, --log-output                  output logs after report (default: false) 
    -n, --namespace string            namespace scope for this request
    -V, --openshift-version string    set the value of certifiedOpenShiftVersions in the report
    -o, --output string               the output format: default, json or yaml
        --registry-config string      path to the registry config file (default "/home/baiju/.config/helm/registry.json")
        --repository-cache string     path to the file containing cached repository indexes (default "/home/baiju/.cache/helm/repository")
        --repository-config string    path to the file containing repository names and URLs (default "/home/baiju/.config/helm/repositories.yaml")
    -s, --set strings                 overrides a configuration, e.g: dummy.ok=false
    -f, --set-values strings          specify application and check configuration values in a YAML file or a URL (can specify multiple)

  Global Flags:
        --config string   config file (default is $HOME/.chart-verifier.yaml)
  ```
- Run a subset of the checks:

  ```
  $ podman run -it --rm quay.io/redhat-certification/chart-verifier verify -e images-are-certified,helm-lint
  ```
- Run all the checks except a subset:

  ```
  $ podman run -it --rm quay.io/redhat-certification/chart-verifier verify -x images-are-certified,helm-lint
  ```
- Provide chart-override values:

  ```
  $ podman run -it --rm quay.io/redhat-certification/chart-verifier verify -S default.port=8080 images-are-certified,helm-lint
  ```
- Provide chart-override values in a file:

  ```
  $ podman run -it --rm quay.io/redhat-certification/chart-verifier verify -F overrides.yaml images-are-certified,helm-lint
  ```

## Chart Testing

### Cluster Config

You can configure the chart-testing check by performing one of the following steps:

* Option 1: Through the `--set` command line option:
    ```text
    $ chart-verifier                                                 \
        verify                                                       \
        --enable chart-testing                                       \
        --set chart-testing.buildId=${BUILD_ID}                       \
        --set chart-testing.upgrade=true                              \
        --set chart-testing.skipMissingValues=true                    \
        --set chart-testing.namespace=${NAMESPACE}                    \
        --set chart-testing.releaseLabel="app.kubernetes.io/instance" \
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
    ```

    Specify the file using the `--set-values` command line option:
    ```text
    $ chart-verifier verify --enable chart-testing --set-values config.yaml some-chart.tgz
    ```

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
