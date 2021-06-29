# adding a command to get information on a report

# command spec

chart-verifier analyze <subcommand> <options> <report-uri>

##sub-commands
* metatdata : return json information for metadata in the report
* digests : return json information for digests in the report
* annotations : return json information for annotations in the report
* results : return json information on the checks
* all : all of the above (default)

##Options

### set prefix for annotations

set the prefix to be used for annotations. Default is ```charts.openshift.io```

* ````--set annotation.prefix=charts.openshift.io````
* ```-a charts.openshift.io```

### set profile vendor type

Set the profile vendor type to use for analysis. Default is vendor type in the report. 
(change verify command set option to ```profile.vendorType``` - currently ```verifier.vendortype```?)

* ```--set profile.vendortype=redhat```
* ```-t redhat```

### set profile vendor type

Set the profile version to use for analysis. Default is vendor type in the report.

* ```--set profile.version=v1.1```
* ```-v v1.1```

## others from verify:

*  ```-h, --help```                        
   * help for analyze
*  ```-o, --output string```               
   * the output format: default, json or yaml    
* ```-s, --set strings```                 
  * overrides a configuration, e.g: profile.vendortype=redhat
* ```-f, --set-values strings```          
  * specify analyze configuration values in a YAML file or a URL (can specify multiple)
    
##Output examples

### metadata

```{ "metadata": { "vendorType": "redhat", "profileVersion": "v1.1" } } ```

### digests

```{ "digests": { "chart": "88888", "package": "8f9b3c" } } ```

### annotations

```{ "annotations": { "OCPVersion": "4.7.8", "digest": "88888", "LastCertifiedTimestamp": "2021-06-29T08:57:35.0023-04:00" } }```

### results

```{ "results": { "success": 11, "fail": 1, Messages [ "Mandatory check chart-testing not found"] }```

### All

```
{ "metadata": { "vendorType": "redhat", "profileVersion": "v1.1" },
  "digests": { "chart": "8f8f8f", "package": "8f9b3c" },
  "annotations": { "charts.openshift.io/OCPVersion": "4.7.8", "charts.openshift.io/digest": "8f8f8f", "charts.openshift.io/LastCertifiedTimestamp": "2021-06-29T08:57:35.0023-04:00" },
  "results": { "passed": 11, "failed": 1, "Messages": [ "Mandatory check chart-testing not found"] } }
```


# design notes

* Command is added to cmd package
  * analyze.go - logic to process command options only.
* Logic is added to 'pkg/chart-verifer/analyze' directory
    * output.go : schema definitions for command output
    * analyzers.go : logic to process subcommands 
* The output must remain consistent.
    * goal: worklow is not affected by profile/report updates unless absolutely necessary.
        * Add but don't take away.
    
Questions:
* do we want to register the analyzers with the command? 
    * similar to how checks are registered