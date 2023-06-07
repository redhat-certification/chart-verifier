# Chart Verifier API

IMPORTANT:
- The API is first published in chart verifier release 1.8.0
- The API in release 1.8.0 is not finalized and is subject to change.
- The API is also available in development release 0.1.0.

## Overview

The chart-verifier API is written using the Go language and is created to enable services to run the chart verifier and obtain a report.
Use the chart-verifier API the same way as any Go package such as import and invoke. 

The chart-verifier API consists of the following Go language packages:

| Package | Description 
| --------| -------------
| [verifier](#verifier) | Provides an API to set the verify flags for the chart-verifier and run the verifier to generate a report. 
| [report](#report) | Provides an API to get and set report content as a string in the JSON or YAML format. 
| [reportSummary](#reportsummary) | Provides an API to set the report flags for the chart-verifier and generate a report summary. 
| [checks](#checks) | Provides an API to get a set containing all available checks. 

Each of these packages are now described in more detail. These are followed by an [example](#example) of use.

## verifier

### Go definition of the APIVerifier interface
```
func NewVerifier() ApiVerifier
type ApiVerifier interface {
	SetBoolean(key BooleanKey, value bool) ApiVerifier
	SetDuration(key DurationKey, duration time.Duration) ApiVerifier
	SetString(key StringKey, value []string) ApiVerifier
	SetValues(key ValuesKey, values map[string]interface{}) ApiVerifier
	EnableChecks(names []apichecks.CheckName) ApiVerifier
	UnEnableChecks(names []apichecks.CheckName) ApiVerifier
	Run(chart_uri string) (ApiVerifier, error)
	GetReport() *report.Report
}
```
### Description of each function of the APIVerifier interface

- NewVerifier - Creates a new ```Verifier```.
  
- SetBoolean: Sets a boolean flag. ```BooleanKey``` values are defined in the verifier package and include:
  - ```WebCatalogOnly```
  - ```Provider Delivery``` (deprecated - replaced by ```WebCatalogOnly```)  
  - ```SuppressErrorLog```
  - ```SkipCleanup```
    
- SetDuration: Sets a duration flag. ```DurationKey``` values are defined in the verifier package and include:
  - ```Timeout```
    
- SetString: Sets a string or string array flag. ```StringKey``` values are defined in the verifier package and include:
  - ```KubeApiServer```
  - ```KubeAsUser```
  - ```KubeCaFile```
  - ```KubeConfig```
  - ```KubeContext```
  - ```KubeToken```
  - ```Namespace```
  - ```OpenshiftVersion```
  -  ```RegistryConfig```
  -  ```RepositoryConfig```
  -  ```RepositoryCache```
  -  ```Config```
  -  ```ChartValues```
  -  ```KubeAsGroups```

- SetValues: Sets a map of string,value pairs. ```ValuesKey``` values are defined in the verifier package and include:
  - ```CommandSet```
  - ```ChartSet```
  - ```ChartSetFile```
  - ```ChartSetString```

- EnableChecks: Specifies a subset of checks to run. 
  - Any checks not specified will not be run. 
  - If no checks are specified all checks will be enabled. 
  - A  list of ```CheckName``` values that you can be enable are defined in the checks package, see [checks](#checks).

- UnEnableChecks: Specifies a subset of checks not to run.
  - Any checks not listed will be run. 
  - If no checks are specified all checks will be enabled. 
  - A list of ```CheckName``` values that you can un-enable are defined in the checks package, see [checks](#checks).

- Run: Runs the verifier based on the flags set and uri provided.

- GetReport: Use after ```Run``` to get the verifier report see [Report](#report).

## Report

### Go definition of the APIReport interface
```
func NewReport() APIReport
type APIReport interface {
    GetContent(ReportFormat) (string, error)
    SetContent(string) APIReport
    SetURL(url *url.URL) APIReport
    Load() (*Report, error)
}
```

### Description of each function of the APIVerifier interface

- NewReport: Creates a new ```Report```.
  
- GetContent: Gets the report as a string in either the JSON or YAML format. ReportFormat values are defined and available in the report package:
  - ```JsonReport``` - for the JSON format.
  - ```YamlReport``` - for the YAML format.
    
- SetContent: Sets the report content from a string, for example a string as returned by ```GetContent```. The format of the report YAML/JSON will be determined based on the report content.
  
- SetUrl: Sets the URL of a report, the report being in string format, for example a string as returned by ```GetContent```. The format of the report YAML/JSON will be determined based on the report content. 
  
- Load: Loads a report based on content set using ```SetContent``` or ```SetUrl```. This will be called internally when the report is needed but can be used to check if a report will load without error.

## ReportSummary

### Go definition of the APIReportSummary interface
```
func NewReportSummary() APIReportSummary
type APIReportSummary interface {
	SetReport(report *apireport.Report) APIReportSummary
	GetContent(SummaryType, SummaryFormat) (string, error)
	SetValues(values map[string]interface{}) APIReportSummary
	SetBoolean(key BooleanKey, value bool) APIReportSummary
}

```

### Description of each function of the APIVerifier interface

- NewReportSummary: Creates a ```ReportSummary```.
  
- SetReport: Sets the report from which the summary should be generated. For example a report as returned by ```report.NewReport```.
  
- GetContent: Gets the report summary as a string in either the JSON or YAML format. ReportFormat values are defined and availalble in the reportsummary package:
    - ```JsonReport``` - for the JSON format.
    - ```YamlReport``` - for the YAML format.
  
- SetValues: Sets value flags to customize content of the report summary. 
  - For example, to customize the result summary to be for a different profile.vendortype than is in the report:
    - set value ```profile.vendortype``` to the required profile (partner/redhat/community).
    - see also: [profiles](helm-chart-checks.md#profiles).

- SetBoolean: Used to set a boolean flag. ```BooleanKey``` values are defined in the reportsummary package and include:
    - ```SkipDigestCheck``` - Intended for testing purpoises only.

## Checks

### Go definition of the GetChecks function
```
func GetChecks() []CheckName
```
### Description of the GetChecks function

- GetChecks: Get an array of ```CheckName``` types. The array content provides the following values that can be used for the verifier.EnableChecks and verifier.UnEnableChecks attributes.
    - ```ChartTesting```
    - ```ContainsTest```
    - ```ContainsValuesSchema```
    - ```ContainsValues```
    - ```HasKubeVersion```
    - ```HasReadme```
    - ```HelmLint```
    - ```ImagesAreCertified```
    - ```IsHelmV3```
    - ```NotContainCsiObjects```
    - ```NotContainsCRDs```
    - ```NotContainValuesSchemaRemoteRef```
    - ```RequiredAnnotationsPresent``` 


# Example:

This example shows a basic invocation of the chart-verifier API, getting and printing the resulting report and the report summary of the report.

Note: 
- The example does not include error checking code for clarity purposes.
  - The example with error checking is available [here](https://github.com/redhat-certification/chart-verifier/blob/main/pkg/chartverifier/samples/sample.go).
- For full use of the chart-verifier see:
    - [https://github.com/redhat-certification/chart-verifier/blob/main/cmd/verify.go](https://github.com/redhat-certification/chart-verifier/blob/main/cmd/verify.go)
    - [https://github.com/redhat-certification/chart-verifier/blob/main/cmd/report.go](https://github.com/redhat-certification/chart-verifier/blob/main/cmd/report.go)

1. Import the packages of the chart-verifier API:
```
import (
	"fmt"
	"github.com/redhat-certification/chart-verifier/pkg/chartverifier/checks"
	"github.com/redhat-certification/chart-verifier/pkg/chartverifier/report"
	"github.com/redhat-certification/chart-verifier/pkg/chartverifier/reportsummary"
	"github.com/redhat-certification/chart-verifier/pkg/chartverifier/verifier"
)
```

2. Set a profile.vendortype of redhat, un-enable the chart testing check, and run verify:
```
	// Run verification for a chart, but omit the chart testing check and run checks based on the redhat profile
	commandSet := make(map[string]interface{})
	commandSet["profile.vendortype"] = "redhat"

	verifier, verifierErr := verifier.NewVerifier().
		SetValues(verifier.CommandSet, commandSet).
		UnEnableChecks([]checks.CheckName{checks.ChartTesting}).
		Run("https://github.com/redhat-certification/chart-verifier/blob/main/tests/charts/psql-service/0.1.9/psql-service-0.1.9.tgz?raw=true")

```
3. Get and print the report, created from the previous step, in the YAML format:
```
	// Get and print the report from the verify command
	report, reportErr := verifier.GetReport().
		GetContent(report.YamlReport)
	fmt.Println("report content:\n", report)
```
4. Set the `profile.vendortype` value to `partner`, get and print a report summary of the previous report in the JSON format. 
```
	// Get and print the report summary  of the report, but using the partner profile.
	values := make(map[string]interface{})
	values["profile.vendortype"] = "partner"

	reportSummary, summmaryErr := reportsummary.NewReportSummary().
		SetReport(verifier.GetReport()).
		SetValues(values).
		GetContent(reportsummary.AllSummary, reportsummary.JsonReport)

	fmt.Println("report summary content:\n", reportSummary)
```
