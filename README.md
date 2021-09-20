# **chart-verifier**: Rules based tool to certify Helm charts

testing

[![Docker Repository on Quay](https://quay.io/repository/redhat-certification/chart-verifier/status "Docker Repository on Quay")](https://quay.io/repository/redhat-certification/chart-verifier)

The **chart-verifier** CLI tool allows you to validate the Helm chart against a configurable list of checks. The tool ensures that the Helm charts include the associated metadata and formatting, and are distribution ready.

The tool allows users to validate a Helm chart URL and provides a report where each check has a `positive` or `negative` result. A negative result from a check indicates a problem with the chart, which needs correction. It ensures that the Helm chart works seamlessly on Red Hat OpenShift and can be submitted as a certified Helm chart in the [OpenShift Helm Repository](https://github.com/openshift-helm-charts).

The input is provided through the command-line interface, with the only required input parameter being the `uri` option. The output is represented through a YAML format with descriptions added in a human-readable format. The report should be submitted with a full set of checks thus validating the Helm chart.

The tool provides the following features:

-   Helm chart verification: Verifies if a Helm chart is compliant with a certain set of independent checks with no particular execution order.
-   Red Hat OpenShift Certified chart validation: Verifies the Helm chart's readiness for being certified and submitted in the OpenShift Helm Repository.    
-   Report generation: Generates a verification report in a YAML format.    
-   Customizable checks: Defines the checks you wish to execute during the verification process.

For more information see:

- [The command line interface and checks performed.](docs/helm-chart-checks.md)
- [Annotations in the report,  chart-verifier CLI tool, and submitter provided.](docs/helm-chart-annotations.md)
- [Introduction to the submission process.](docs/helm-chart-submission.md)
- [Troubleshooting check failures.](docs/helm-chart-troubleshooting.md)

For developer specific information, see:

- [Additional information for developers.](docs/helm-chart-developer.md)
- [Creating a chart-verifier release.](docs/helm-chart-release.md)
