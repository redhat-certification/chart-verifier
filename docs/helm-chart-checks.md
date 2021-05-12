# Helm chart checks for Red Hat certification

Helm chart checks are a set of checks against which the Red Hat Helm chart-verifier tool verifies and validates whether a Helm chart is qualified for a certification from Red Hat. These checks contain metadata that have certain parameters and values with which a Helm chart must comply to be certified by Red Hat. A Red Hat-certified Helm chart qualifies in terms of readiness for distribution in the [Red Hat OpenShift Helm Chart Repository](https://github.com/openshift-helm-charts).

## Key features
- You can execute the checks individually. Individual checks can be included or excluded during the verification process using the command line options.
- Each check is independent, and there is no specific execution order.
- Executing the checks includes running the `helm lint` and `helm template` commands against the chart for Red Hat image certification.
- The checks are configurable. For example, if a chart requires additional values to be compliant with the checks, configure the values using the options available. The options are similar to those used by the `helm lint` and `helm template` commands.
- When there are no error messages, the `helm-lint` check passes the verification and is successful. Messages such as `Warning` and `info` do not cause the check to fail.

## Types of Helm chart checks
Helm chart checks are categorized into the following types:

#### Table 1: Helm chart check types

| Check type | Description
|---|---
| Mandatory | Checks are required to pass and be successful for certification.
| Recommended | Checks are about to become mandatory; we recommend fixing the check failures, if any.
| Optional | Checks are ready for customer testing. Checks can fail and still pass the verification for certification.
| Experimental | New checks introduced for testing purposes or beta versions.
> **_NOTE:_**  The current release of Helm charts verify only the mandatory type of checks.

## Default set of checks for a Helm chart
The following table lists the default set of checks with details including the names of the checks, their type, and description.These checks must be successful for a Helm chart to be certified by Red Hat.

#### Table 2: Helm chart default checks

| Name | Check type | Description
|------|------------|------------
| `is-helm-v3` | Mandatory | Checks that the given `uri` points to a Helm v3 chart.
| `has-readme` | Mandatory | Checks that the Helm chart contains the `README.md` file.
| `contains-test` | Mandatory | Checks that the Helm chart contains at least one test file.
| `has-kubeversion` | Mandatory | Checks that the `Chart.yaml` file of the Helm chart includes the `kubeVersion` field.
| `contains-values-schema` | Mandatory | Checks that the Helm chart contains a JSON schema file (`values.schema.json`) to validate the `values.yaml` file in the chart.
| `not-contains-crds` | Mandatory | Checks that the Helm chart does not include custom resource definitions (CRDs).
| `not-contain-csi-objects` | Mandatory | Checks that the Helm chart does not include Container Storage Interface (CSI)[¹](https://gist.github.com/Srivaralakshmi/eb3c9bf1d65ec297035f4a8d26057620#-for-more-information-on-csi-see-container-storage-interface) objects.
| `images-are-certified` | Mandatory | Checks that the images referenced by the Helm chart are Red Hat-certified.
| `helm-lint` | Mandatory | Checks that the chart is well formed by running the `helm lint` command.
| `contains-values` | Mandatory | Checks that the Helm chart contains the `values`[²](https://gist.github.com/Srivaralakshmi/eb3c9bf1d65ec297035f4a8d26057620#-for-more-information-on-the-values-file-see-values-and-best-practices-for-using-values) file.

###### ¹ For more information on CSI, see [Container Storage Interface](https://github.com/container-storage-interface/spec/blob/master/spec.md).

###### ² For more information on the `values` file, see [`values`](https://helm.sh/docs/chart_template_guide/values_files/) and [Best Practices for using values](https://helm.sh/docs/chart_best_practices/values/).

## Using the docker command for Helm chart checks
This section provides help on the basic usage of Helm chart checks with the docker command.

### Prerequisites
- A container engine and the `docker` command installed.
- Internet connection to check that the images are Red Hat certified.
- itHub profile to submit the chart to the [Red Hat OpenShift Helm Charts Repository](https://github.com/openshift-helm-charts).

### Procedure
To perform the tasks related to Helm chart checks, run the following commands:

- Runs all the available checks for the chart using a `uri`.

 ```ruby
 $ docker run -it --rm quay.io/redhat-certification/chart-verifier verify <chart-uri>
 ```
- Runs all the checks available locally on your system for the chart, from the same directory as the chart.

 ```ruby
 $ docker run -v $(pwd):/charts --rm quay.io/redhat-certification/chart-verifier verify /charts/<chart>
 ```
- Gets the list of options for the `verify` command.

 ```ruby
 $ docker run -it --rm quay.io/redhat-certification/chart-verifier verify help
 ```
 The output is similar to the following example:
 ```ruby
 $ docker run -it --rm quay.io/redhat-certification/chart-verifier verify help`

 Verifies a Helm chart by checking some of its characteristics

 Usage:
 chart-verifier verify <chart-uri> [flags]

 Flags:
   -S, --chart-set strings          set values for the chart (can specify multiple or separate values with commas:     key1=val1,key2=val2)
   -F, --chart-set-file strings     set values from respective files specified via the command line (can specify multiple or  separate values with commas: key1=path1,key2=path2)
   -X, --chart-set-string strings   set STRING values for the chart (can specify multiple or separate values with commas: key1=val1,key2=val2)
   -f, --chart-values strings       specify values in a YAML file or a URL (can specify multiple)
   -x, --disable strings            all checks will be enabled except the informed ones
   -e, --enable strings             only the informed checks will be enabled
   -h, --help                       help for verify
   -o, --output string              the output format: default, json or yaml
   -s, --set strings                overrides a configuration, e.g: dummy.ok=false

 Global Flags:
       --config string   config file (default is $HOME/.chart-verifier.yaml)
- Runs a subset of the checks.

 ```ruby
 $ docker run -it --rm quay.io/redhat-certification/chart-verifier verify -e images-are-certified,helm-lint
 ```
- Runs all the checks except a subset.

 ```ruby
 $ docker run -it --rm quay.io/redhat-certification/chart-verifier verify -x images-are-certified,helm-lint
 ```
- Provides chart-override values.

 ```ruby
 $ docker run -it --rm quay.io/redhat-certification/chart-verifier verify -S default.port=8080 images-are-certified,helm-lint
 ```
- Provides chart-override values in a file.

 ```ruby
 $ docker run -it --rm quay.io/redhat-certification/chart-verifier verify -F overrides.yaml images-are-certified,helm-lint
 ```
