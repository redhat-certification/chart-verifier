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

An exception process starts if the Pull request fails, if the report contains one or more failures, or has missing tests. If you submit the report without a chart, run the report against the chart in the corresponding final location. Then, the verifier records the `chart-uri` specified during the report run, and in the absence of a submitted chart, you can use the `chart-uri` for publication.

If the report is to be submitted with a chart, it must be run against the same. This is because the submission process does not have access to the values and the report generated would include failures.

For more information on the submission process, please refer to [OpenShift Helm Charts Repository documentation](https://github.com/openshift-helm-charts/charts/blob/main/docs/README.md).
