# Submission of Helm Charts for verification

OpenShift Helm Charts repository hosts Helm Charts that are available by default with OpenShift. You can use this repository to submit the charts that need to be certified through a Pull Request. Once your Pull Request is merged, a CI/CD pipeline is created which publishes the chart in the GitHub releases which is further reflected on the Helm charts [repository index](http://charts.openshift.io/).

The submission process of a Helm Chart for Red Hat Helm Repository and Certification has been documented on the [Red Hat Helm Repository](https://github.com/openshift-helm-charts/charts). Please note the instructions mentioned on the repository before submitting a chart.

The following options are available for submitting a Chart for inclusion in Red Hat Helm Repository and Certification: 

| Option                                       | Description                                                                                                             |
|----------------------------------------------|-------------------------------------------------------------------------------------------------------------------------|
| Helm Chart Tarball or the extracted Tarball  | Submit your Chart with the specific or the extracted tarball. Here the chart-verifier's report is optional.             |
| Chart Verification Report only               | Submit your Chart Verification Report without the Chart.                                                                |
| Both Chart Verification Report and the Chart | Submit both the Chart Verification Report and the Chart by placing the source or tarball under the versioned directory. |

> **_NOTE:_**  A verifier report is an integral part of the submission process. With the options that do not require a report, a report will be generated as part of the submission process.

> **_NOTE:_**  It is preferred when submitting a chart to submit chart source over a tarball. 

For more information on the submission process, please refer to [OpenShift Helm Charts Repository documentation](https://github.com/openshift-helm-charts/charts/blob/main/docs/README.md).

## Submission failures from a submitted report

- One or more mandatory checks have failed or are missing from the report. 
  - Submission will fail if any [mandatory checks](./helm-chart-checks.md#default-set-of-checks-for-a-helm-chart) indicate failure or are absent from the report.
  - If you are unsure as to why a check failed see [Trouble shooting check failures](./helm-chart-troubleshooting.md#trouble-shooting-check-failures)  
- The digest in the report does not match the digest calculated for the submitted chart. Common issues:
  - Chart was updated after the report was generated.
  - Report was generated against a different form to the chart submitted. For example report was generated from the chart source but the chart tarball was used for submission.
  - For more information see [Verifier added annotations](./helm-chart-annotations.md#verifier-added-annotations)  
- The certifiedOpenShiftVersions does not contain a valid value.
  - This annotation must contain a current or recent OpenShift version. 
  - For more information see [Verifier added annotations](./helm-chart-annotations.md#verifier-added-annotations)
- The chart uri is not a valid url. 
    - For a report only submission the report must include a valid url for the chart.
    - For more information see [error-with-the-chart-url-when-submitting-report](https://github.com/openshift-helm-charts/charts/blob/main/docs/README.md#error-with-the-chart-url-when-submitting-report)
    - For more information see [Verifier added annotations](./helm-chart-annotations.md#verifier-added-annotations)
   