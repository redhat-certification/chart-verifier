# Red Hat Helm Chart Verifier Part 1

## Introducion

Red Hat provides a chart verifier tool for use by partners planning to submit a chart for 
certification and subsequent availability in the openshift catalogue. In this blog I cover
running the chart verifier and intepreting the results.

- add informartion on why to use it.

## Setting up to run the verifier

First you will need either [podman](https://podman.io/getting-started/installation) or [docker](https://docs.docker.com/engine/install/).

Then a cluster for the verifier to run tests on. There are lots of choices available, including:
- [codeready containers](https://developers.redhat.com/products/codeready-containers/getting-started).
- [minikube](https://minikube.sigs.k8s.io/docs/start/).
- remote OpenShift cluster.


## Running the chart verifier
- Basic podman/docker commands and the output to expect 

### Specify chart values
- using the varu=ious optin for setting chart values.

### Running a subset of tests
- how to run a subset and why you might do it.

## Dealing with failures
- basically point to troubleshooting on doc.
- what it you can't fix them all
  - exceptions
  - community  
    
## When you need a report
- submission with report v without.

# More information
- doc links


# Red Hat Helm Chart Verifier Part 1

## Running locally
- clone and build

## Using profiles
- intro to profiles and how to use them

## Creating your own profile.
- how to create you own profile and why you might do it.

## Adding checks
- not recommended, open an issue.
