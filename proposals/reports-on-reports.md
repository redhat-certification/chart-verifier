# Adding commands to get information on reports and profiles

## Background

The CI workflow in the chart repository currently contains report specific information which leads to several problems.
* The worklflow may need to be updated when the report changes.
* The workflow needs separate paths for different versions of the report.
* The workflow contains logic which duplicates logic in the verifier code.
* The workflow needs access to the profiles

The new command is therefore design primarily for use from the workflow to isolate the workflow as much as possible from future updates to the report format and content.However, it may also be useful for users and will documented.

# Command specs

```chart-verifier report <subcommand> <options> <report-uri>```

```chart-verifier profile <subcommand> <options>```

## Sub-Commands

### Report Sub-commands
* metatdata : return information for metadata in the report
* digests : return information for digests in the report
* annotations : return information for annotations in the report
* results : return information on the checks
* all : all of the above (default) 

    
### Profile Sub-commands 
* show : output the content of a profile in yaml (default is all)
* list  : list the profile and version available
    
## Options

### Report Options 

* Set the prefix to be used for annotations. 
    * Default is ```charts.openshift.io```
    * for example: ````--set annotation.prefix=charts.openshift.io````

### Common Options

* Set the profile vendor type.
    * report command default is vendor type in the specified report.
    * profile command default is all vendor types.
    * for example: ```--set profile.vendortype=redhat```
    

* Set the profile version
    * report command default is version in the specified report.
    * profile command default is all versions.
    * for example: ```--set profile.version=v1.1```

## common options from verify:

*  ```-h, --help```                        
   * help for the specified command
*  ```-o, --output string```               
   * the output format: default, json or yaml. Default is json for all commands and subcommands except display.    
* ```-s, --set strings```                 
  * set a configuration, e.g: ```-s profile.vendortype=redhat```
* ```-f, --set-values strings```          
  * specify analyze configuration values in a YAML file or a URL (can specify multiple)
    
## Command and Output examples

### report metadata

command: ```chart-verifier report metadata report.yaml```

```{ "metadata": { "vendorType": "redhat", "profileVersion": "v1.1" } } ```

### report digests

command: ```chart-verifier report digests report.yaml```

```{ "digests": { "chart": "88888", "package": "8f9b3c" } } ```

### report annotations

command: ```chart-verifier report annotations report.yaml```

```{ "annotations": { "OCPVersion": "4.7.8", "digest": "88888", "LastCertifiedTimestamp": "2021-06-29T08:57:35.0023-04:00" } }```

### report results

command: ```chart-verifier report results report.yaml```

```{ "results": { "success": 11, "fail": 1, Messages [ "Mandatory check chart-testing not found"] }```

### report Full

command: ```chart-verifier report report.yaml```

```
{ "metadata": { "vendorType": "redhat", "profileVersion": "v1.1" },
  "digests": { "chart": "8f8f8f", "package": "8f9b3c" },
  "annotations": { "charts.openshift.io/OCPVersion": "4.7.8", "charts.openshift.io/digest": "8f8f8f", "charts.openshift.io/LastCertifiedTimestamp": "2021-06-29T08:57:35.0023-04:00" },
  "results": { "passed": 11, "failed": 1, "Messages": [ "Mandatory check chart-testing not found"] } }
```

## profile list

command: ```chart-verifier profile list```
```
{ "profiles": [{"type": "partner", "version": "v1.1"},{"type": "redhat", "version": "v1.1"},{"type": "community", "version": "v1.1"}]} 
```

## profile show -s profile.vendortype=partner, profile.version=v1.1

command: ```chart-verifier profile show -s profile.vendortype=partner, profile.version=v1.1```

```
apiversion: v1
kind: verifier-profile
vendorType: partner
version: 1.1
annotations:
  - "Digest"
  - "OCPVersion"
  - "LastCertifiedTimestamp"
checks:
    - name: v1.0/has-readme
      type: Mandatory
    - name: v1.0/is-helm-v3
      type: Mandatory
    - name: v1.0/contains-test
      type: Mandatory
    - name: v1.0/contains-values
      type: Mandatory
    - name: v1.0/contains-values-schema
      type: Mandatory
    - name: v1.0/has-kubeversion
      type: Mandatory
    - name: v1.0/not-contains-crds
      type: Mandatory
    - name: v1.0/helm-lint
      type: Mandatory
    - name: v1.0/not-contain-csi-objects
      type: Mandatory
    - name: v1.0/images-are-certified
      type: Mandatory
    - name: v1.0/chart-testing
      type: Mandatory
```

# design notes

* Command is added to cmd package
  * report.go - logic to process report command options only.
  * profile.go - logic to process profile command options only/  
* Logic is added to 'pkg/chart-verifer/commands/report' directory
    * output.go : schema definitions for command output
    * reporter.go : logic to process subcommands
    * commandBuilder.go : register command functions with the command
* Logic is added to 'pkg/chart-verifer/commands/profile' directory
    * output.go : schema definitions for command output
    * profiler.go : logic to process subcommands
    * commandBuilder.go : register command functions with the command    
* The output must remain consistent.
    * goal: worklow is not affected by profile/report updates unless absolutely necessary.
        * Add but don't take away.
    