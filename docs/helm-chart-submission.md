
# Submission of Helm charts for Red Hat OpenShift certification
 - [Submission options](#submission-options)
 - [Provider controlled delivery](#provider-controlled-delivery)

## Submission options

OpenShift Helm charts repository hosts Helm charts that are available by default with OpenShift. You can use this repository to submit the charts that need to be certified through a pull request. Once your pull request is merged, a CI/CD pipeline is created which publishes the chart in the gitHub release which is further reflected on the Helm charts [repository index](http://charts.openshift.io/).

The submission process of a Helm chart for OpenShift Helm Repository and Certification has been documented on the [OpenShift Helm Repository](https://github.com/openshift-helm-charts/charts). Note the instructions mentioned on the repository before submitting a chart.

The following options are available for submitting a chart for Red Hat OpenShift certification: 

| Option                                       | Description                                                                                                             |
|----------------------------------------------|-------------------------------------------------------------------------------------------------------------------------|
| Helm chart Tarball or the extracted Tarball  | Submit your chart with the specific or the extracted tarball. Here the chart-verifier report is optional.             |
| Verification report only                     | Submit your chart-verifier report without the chart.                                                                |
| Both verification report and the chart       | Submit both the chart-verifier report and the chart by placing the source or tarball under the versioned directory. |

> **_NOTE:_**  A chart-verifier report is an integral part of the submission process. With the options that do not require a report, a report will be generated as part of the submission process.

> **_NOTE:_**  It is recommended when submitting a chart to submit chart source over a tarball. 

> **_NOTE:_**  When submitting a Verification report only the report must be generated using the public url for the chart. 

For more information on the submission process, see: [OpenShift Helm Charts Repository documentation](https://github.com/openshift-helm-charts/charts/blob/main/docs/README.md).

For troubleshooting report related submission failures see: [Troubleshooting](./helm-chart-troubleshooting.md)

## Provider controlled delivery

By default, a submitted chart will be made available in the OpenShift Helm Chart Catalog on successful certification. In some cases this is undesirable and can be prevented using provider controlled delivery. With provider controlled delivery the provider of the chart controls access to the chart and this impacts report generation:

- The report must be generated using a tarball so that a package digest can be determined and included in the report.
  - if a tarball is not used the report will fail to generate.
- The chart URL may be considered private to the provider so the chart URL is not included in the report.

Provider controlled delivery is then based on the following conditions: 

1. When generating the Verification report the ```--provider-delivery``` flag is used.
   Example:
    ```
    $ podman run --rm -i                                  \
          -e KUBECONFIG=/.kube/config                   \
          -v "${HOME}/.kube":/.kube                     \
          "quay.io/redhat-certification/chart-verifier" \
          verify --provider-delivery                    \
          <chart-uri>
    ```
    This ensures that the [providerControlledDelivery annotation](helm-chart-annotations.md#providerControlledDelivery) is set to the value True in the verification report.

1. The OWNERS file for the submitted chart in the [openshift helm charts github repository](https://github.com/openshift-helm-charts/charts) includes a ```providerDelivery``` attribute which is set to the value True. 
   Example:
```
chart:
  name: mychart
  shortDescription: Test chart for testing chart submission workflows.
publicPgpKey: null
providerDelivery: True
users:
- githubUsername: myusername
vendor:
  label: redhat
  name: Redhat
```

If these preceding conditions are met when the chart is submitted for certification, successful certification will not result in the chart being published in the OpenShift catalog.