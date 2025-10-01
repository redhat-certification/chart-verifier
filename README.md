# **chart-verifier**: Rules based tool to certify Helm charts

The **chart-verifier** CLI tool allows you to validate the Helm chart against a configurable list of checks. The tool ensures that the Helm charts include the associated metadata and formatting, and are distribution ready.

The tool allows users to validate a Helm chart URL and provides a report where each check has a `positive` or `negative` result. A negative result from a check indicates a problem with the chart, which needs correction. It ensures that the Helm chart works seamlessly on Red Hat OpenShift and can be submitted as a certified Helm chart in the [OpenShift Helm Repository](https://github.com/openshift-helm-charts).

The input is provided through the command-line interface, with the only required input parameter being the `uri` option. The output is represented through a YAML format with descriptions added in a human-readable format. The report should be submitted with a full set of checks thus validating the Helm chart.

The tool provides the following features:

-   Helm chart verification: Verifies if a Helm chart is compliant with a certain set of independent checks with no particular execution order.
-   Red Hat OpenShift Certified chart validation: Verifies the Helm chart's readiness for being certified and submitted in the OpenShift Helm Repository.    
-   Report generation: Generates a verification report in a YAML format.    
-   Customizable checks: Defines the checks you wish to execute during the verification process.

For more information see:

- [The command line interface and checks performed](docs/helm-chart-checks.md)
    - [Key features](docs/helm-chart-checks.md#key-features)
    - [Types of Helm chart checks](docs/helm-chart-checks.md#types-of-helm-chart-checks)
    - [Default set of checks for a Helm chart](docs/helm-chart-checks.md#default-set-of-checks-for-a-helm-chart)
    - [Run Helm chart checks](docs/helm-chart-checks.md#run-helm-chart-checks)
    - [Profiles](docs/helm-chart-checks.md#profiles)
    - [Chart Testing](docs/helm-chart-checks.md#chart-testing)
- [Certification Annotations](docs/helm-chart-annotations.md)
    - [Verifier added annotations](docs/helm-chart-annotations.md#verifier-added-annotations)
    - [Annotations by profile](docs/helm-chart-annotations.md#annotations-by-profile)
    - [Provider annotations](docs/helm-chart-annotations.md#provider-annotations)
- [Introduction to the submission process](docs/helm-chart-submission.md)
    - [Submission options](docs/helm-chart-submission.md#submission-options)
    - [Provider controlled delivery](docs/helm-chart-submission.md#provider-controlled-delivery)
- [Troubleshooting](docs/helm-chart-troubleshooting.md)
    - [Check failures](docs/helm-chart-troubleshooting.md#troubleshooting-check-failures)
    - [Report related submission failures](docs/helm-chart-troubleshooting.md#report-related-submission-failures)
- [Chart Verifier API](docs/helm-chart-api.md)
    - [Overview](docs/helm-chart-api.md#overview)
    - [Verifier](docs/helm-chart-api.md#verifier)
    - [Report](docs/helm-chart-api.md#report)
    - [Report Summary](docs/helm-chart-api.md#reportsummary)
    - [Checks](docs/helm-chart-api.md#checks)  
    - [Example](docs/helm-chart-api.md#example)
    

For developer specific information, see:

- [Additional information for developers.](docs/helm-chart-developer.md)
- [Creating a chart-verifier release.](docs/helm-chart-release.md)

Sample Change More Context