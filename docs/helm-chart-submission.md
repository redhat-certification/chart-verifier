# Submission of Helm charts for certification

OpenShift Helm charts repository hosts Helm charts that are available by default with OpenShift. You can use this repository to submit the charts that need to be certified through a pull request. Once your pull request is merged, a CI/CD pipeline is created which publishes the chart in the gitHub release which is further reflected on the Helm charts [repository index](http://charts.openshift.io/).

The submission process of a Helm chart for OpenShift Helm Repository and Certification has been documented on the [OpenShift Helm Repository](https://github.com/openshift-helm-charts/charts). Note the instructions mentioned on the repository before submitting a chart.

The following options are available for submitting a chart for inclusion in OpenShift Helm Repository and Certification: 

| Option                                       | Description                                                                                                             |
|----------------------------------------------|-------------------------------------------------------------------------------------------------------------------------|
| Helm chart Tarball or the extracted Tarball  | Submit your chart with the specific or the extracted tarball. Here the chart-verifier report is optional.             |
| Verification report only                     | Submit your chart verification report without the chart.                                                                |
| Both verification report and the chart       | Submit both the chart verification report and the chart by placing the source or tarball under the versioned directory. |

> **_NOTE:_**  A chart verification report is an integral part of the submission process. With the options that do not require a report, a report will be generated as part of the submission process.

> **_NOTE:_**  It is recommended when submitting a chart to submit chart source over a tarball. 

For more information on the submission process, see: [OpenShift Helm Charts Repository documentation](https://github.com/openshift-helm-charts/charts/blob/main/docs/README.md).

For troubleshooting report related submission failures see: [Troubleshooting](./helm-chart-troubleshooting.md)
 