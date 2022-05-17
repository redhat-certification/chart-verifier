# Certification Annotations

- [Verifier added annotations](#verifier-added-annotations)
- [Annotations by profile](#annotations-by-profile)
  - [verifier-version](#verifier-version)
  - [profile](#profile)  
  - [chart-uri](#chart-uri)
  - [digests](#digests) 
  - [lastCertifiedTimestamp](#lastCertifiedTimestamp)  
  - [certifiedOpenShiftVersions](#certifiedOpenShiftVersions)
  - [testedOpenShiftVersion](#testedOpenShiftVersion)
  - [supportedOpenShiftVersions](#supportedOpenShiftVersions)
  - [providerControlledDelivery](#providerControlledDelivery)  
- [Provider annotations](#provider-annotations)
  - [charts.openshift.io/provider](#chartsopenshiftioprovider)
  - [charts.openshift.io/name](#chartsopenshiftioname)
  - [charts.openshift.io/supportURL](#chartsopenshiftiosupportURL)
  - [charts.openshift.io/archs](#chartsopenshiftioarchs)  


## Verifier added annotations


The chart-verifier tool adds annotations to a generated report, for example:

```
metadata:
    tool:
        verifier-version: 1.4.0
        profile:
            VendorType: partner
            version: v1.1
        chart-uri: https://github.com/mmulholla/development/blob/main/charts/partners/test-org/psql-service/0.1.9/psql-service-0.1.9.tgz?raw=true
        digests:
            chart: sha256:94cbcb63531bc4457e7b0314f781070bbfe4affbdca98f67acadc381bf0f0b4f
            package: 4e9592ea31c0509bec308905289491b7056b78bdde2ab71a85c72be0901759b8
        lastCertifiedTimestamp: "2021-11-01T17:12:37.148895-04:00"
        testedOpenShiftVersion: 4.8
        supportedOpenShiftVersions: 4.5 - 4.8
        providerControlledDelivery: false
 
```

The annotations added differ based on the profiles version used:

## Annotations by profile

| Annotation                 | Profile Versions |
| -------------------------- |:-----------------
| [verifier-version](#verifier-version)                     | v1.0, v1.1
| [profile](#profile)                                       | v1.0, v1.1
| [chart-uri](#chart-uri)                                   | v1.0, v1.1
| [digests](#digests)                                       | v1.0, v1.1
| [lastCertifiedTimestamp](#lastCertifiedTimestamp)         | v1.0, v1.1
| [certifiedOpenShiftVersions](#certifiedOpenShiftVersions) | v1.0 
| [testedOpenShiftVersion](#testedOpenShiftVersion)         | v1.1
| [supportedOpenShiftVersions](#supportedOpenShiftVersions) | v1.1
| [providerControlledDelivery](#providerControlledDelivery) | v1.0,v1.1 

### verifier-version

The version of the chart-verifier which generated the report. 

### profile

profile incluse information aboy the profile that was used to generate the report.

### chart-uri

The location of the chart specified to the chart-verifier. For report-only submissions, this must be the public url of the chart.                                                                      |

### digests

digests may includes two digests:
- digests.chart:
    - sha:256 value of the chart as calculated from the copy of the chart loaded into memory by the chart-verifier.  
    - When submitting a report, this value must match the value generated as part of the submission process.
- digest.package:
    - The sha value of the chart tarball if used to create the report.
    - Not included if chart source was used.
    
### lastCertifiedTimestamp

The time when the report was generated.

### certifiedOpenShiftVersions

- The version of OCP that `chart-testing` check was performed on. If the role of the logged-in user prevents this from being accessed, the value must be specified using the `--openshift-version` flag.
- If the certifiedOpenShiftVersions is not set to a valid OpenShift version, the submission will fail.
- Renamed to testedOpenShiftVersions in profile version v1.1

### testedOpenShiftVersion

- The version of OCP that `chart-testing` check was performed on. If the role of the logged-in user prevents this from being accessed, the value must be specified using the `--openshift-version` flag.
- If the certifiedOpenShiftVersions is not set to a valid OpenShift version, the submission will fail.
- Renamed from certifiedOpenShiftVersions in profile version v1.1

### supportedOpenShiftVersions 

The Open Shift versions supported by the chart based on the kubeVersion attribute in chart.yaml.

### providerControlledDelivery

Used to control publication of a certified chart:
- True: provider will control publication of the chart
- False (default): The chart will be published in the OpenShift Helm chart catalogue when certified.

see: [Provider controlled delivery.](helm-chart-submission.md#provider-controlled-delivery)

## Provider annotations

The chart provider can also include annotations in `Chart.yaml`, which may be used when displaying the chart in the catalog, for example:

```
annotations:
   charts.openshift.io/archs: x86_64
   charts.openshift.io/name: PSQL RedHat Demo Chart
   charts.openshift.io/provider: RedHat
   charts.openshift.io/supportURL: https://github.com/dperaza4dustbit/helm-chart
```

### charts.openshift.io/provider

Name of chart provider (e.g., Red Hat), ready to be displayed in UI.

### charts.openshift.io/name

Human readable chart name, ready to be displayed in UI.
- This is mandatory with profile v1.1 (see: [has-kubeversion v1.1](helm-chart-troubleshooting.md#has-kubeversion-v11)) 

### charts.openshift.io/supportURL

Where users can find information about the chart provider's support.

### charts.openshift.io/archs

Comma separated list of supported architectures (e.g., x86_64, s390x, ...)

