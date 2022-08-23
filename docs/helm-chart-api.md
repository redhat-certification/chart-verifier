# Chart Verifier go language API

IMPORTANT:
- The API was first published in chart verifier release 1.8.0
- The API in release 1.8.0 is not finalized and is subject to change.
- The API is also available in development release 0.1.0.

## Overview

The API is written using the go language and was created to enable services to run the chart verifier and obtain a report. 
It can be used in the same way as any go package, import and invoke. 

The API includes the following 4 go language packages: verifier, report, reportSummary and checks:

| Package | Description 
| --------| -------------
| [verifier](#verifier) | Provides an API to set the verify flags for the chart-verifier and run the verifier to generate a report. 
| [report](#report) | Provides an api to get and set report content as a string in json or yaml format. 
| [reportSummary](#reportsummary) | Provides an api to set chart verifier report flags and generate a report summary. 
| [checks](#checks) | Provides an api to get a set containing all available checks. 

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

- NewVerifier - Used to create a new ```Verifier```.
  
- SetBoolean: Used to set a boolean flag. ```BooleanKey``` values are defined in the verifier package and include:
  - ```ProviderDelivery```
  - ```SuppressErrorLog```
    
- SetDuration: Used to set a duration flag. ```DurationKey``` values are defined in the verifier package and include:
  - ```Timeout```
    
- SetString: Used to set a string or string array flag. ```StringKey``` values are defined in the verifier package and include:
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

- SetValues: Used to set a map of string,value pairs. ```ValuesKey``` values are defined in the verifier package and include:
  - ```CommandSet```
  - ```ChartSet```
  - ```ChartSetFile```
  - ```ChartSetString```

- EnableChecks: Used to specify a subset of checks to run, any checks not listed will not be run. If called with an empty list all checks will be enabled. For a list of ```CheckName``` values which can be enabled are defined in the checks package, see [checks](#checks).

- UnEnableChecks: Used to specify a subset of checks which should not run, any checks not listed will be run. If called with an empty list there will be no effect. For a list of ```CheckName``` values which can be un-enabled are defined in the checks package, see [checks](#checks).

- Run: Used to run the verifier verify command based on the flags set and uri provided.

- GetReport: used after run to get the verifier report see [Report](#report).

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

- NewReport: Used to create a new ```Report```.
  
- GetContent: Get the report as a string in either json or yaml format. ReportFormat values are defined and available in the report package:
  - ```JsonReport``` - for json format.
  - ```YamlReport``` - for yaml format.
    
- SetContent: Set the report content from a string, for example a string as returned by ```GetContent```. The format of the report yaml/json will be determined based on the report content.
  
- SetUrl: Set the URL of a report, the report being in string format, for example a string as returned by ```GetContent```. The format of the report yaml/json will be determined based on the report content. 
  
- Load: Can be used to load a report based on content set using ```SetContent``` or ```SetUrl```. This will be called internally when the report is needed but can be used to check if a report will load without error.

## ReportSummary

### Go definition of the APIReportSummary interface
```
func NewReportSummary() APIReportSummary
type APIReportSummary interface {
	SetReport(report *apireport.Report) APIReportSummary
	GetContent(SummaryType, SummaryFormat) (string, error)
	SetValues(values map[string]interface{}) APIReportSummary
}

```

### Description of each function of the APIVerifier interface

- NewReportSummary: Used to create a ```ReportSummary```.
  
- SetReport: Set the report from which the summary should be generated. For example a report as returned by ```report.NewReport```.
  
- GetContent: Get the report summary as a string in either json or yaml format. ReportFormat values are defined and availalble in the reportsummary package:
    - ```JsonReport``` - for json format.
    - ```YamlReport``` - for yaml format.
  
- SetValues: Used to set value flags to tailor content of the report sumary. Can be used to set a profile vendorType and/or version which by default are the values set in the report.

## Checks

### Go definition of the GetChecks function
```
func GetChecks() []CheckName
```
### Description of the GetChecks function

- GetChecks: Used to get an array of ```CheckName``` types. The array content provide value with can be used for ```verifier.EnableChecks``` and ```verifier.UnEnableChecks```:
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
    - ```RequiredAnnotationsPresent``` 


# Example:

Note: 
- The following example does not include error checking code for clarity purposes.
- For full use of the chart-verifier see:
    - [https://github.com/redhat-certification/chart-verifier/blob/main/cmd/verify.go](https://github.com/redhat-certification/chart-verifier/blob/main/cmd/verify.go)
    - [https://github.com/redhat-certification/chart-verifier/blob/main/cmd/report.go](https://github.com/redhat-certification/chart-verifier/blob/main/cmd/report.go)

1. import the packages of the chart-verifier API:
```
import (
	"fmt"
	"github.com/redhat-certification/chart-verifier/pkg/chartverifier/checks"
	"github.com/redhat-certification/chart-verifier/pkg/chartverifier/report"
	"github.com/redhat-certification/chart-verifier/pkg/chartverifier/reportsummary"
	"github.com/redhat-certification/chart-verifier/pkg/chartverifier/verifier"
)
```

2. Set a profile vendortype of redhat, un-enable the chart testing check, and run verify:
```
	// Run verify command for a chart, but omit the chart testing check and run checks based on the redhat profile
	commandSet := make(map[string]interface{})
	commandSet["profile.vendortype"] = "redhat"

	verifier, verifierErr := verifier.NewVerifier().
		SetValues(verifier.CommandSet, commandSet).
		UnEnableChecks([]checks.CheckName{checks.ChartTesting}).
		Run("https://github.com/redhat-certification/chart-verifier/blob/main/tests/charts/psql-service/0.1.9/psql-service-0.1.9.tgz?raw=true")

```
3. Get and print the report created from the previous step in YAML format:
```
	// Get and print the report from the verify command
	report, reportErr := verifier.GetReport().
		GetContent(report.YamlReport)
	fmt.Println("report content:\n", report)
```
4. Use the vendortype value of the profile as partner to get and print a report summary of the previous report in JSON format. 
```
	// Get and print the report summary  of the report, but using the partnet profile.
	values := make(map[string]interface{})
	values["profile.vendortype"] = "partner"

	reportSummary, summmaryErr := reportsummary.NewReportSummary().
		SetReport(verifier.GetReport()).
		SetValues(values).
		GetContent(reportsummary.AllSummary, reportsummary.JsonReport)

	fmt.Println("report summary content:\n", reportSummary)
```

The above example shows a basic invocation of the chart verifier, getting and printing the resulting report and the report summary of the report. 
- For full use of the chart-verifier API see:
    - [https://github.com/redhat-certification/chart-verifier/blob/main/cmd/verify.go](https://github.com/redhat-certification/chart-verifier/blob/main/cmd/verify.go)
    - [https://github.com/redhat-certification/chart-verifier/blob/main/cmd/report.go](https://github.com/redhat-certification/chart-verifier/blob/main/cmd/report.go)
