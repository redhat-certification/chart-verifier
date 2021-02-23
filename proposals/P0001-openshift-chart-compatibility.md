# OpenShift Chart Compatibility

* Proposal N.: 0001
* Authors: isuttonl@redhat.com
* Status: **Draft**

## Abstract

This feature proposal outlines a mechanism to verify whether a chart is compatible with a particular version of
OpenShift by comparing the Helm Chart's rendered resources resource information (namely the GVK) against a set of
supported resources by default at a given OpenShift version.

Please note this also could be extended to verify against other Kubernetes distributions as well, but is at the time of
this writing a non-goal.

## Motivation

There is currently no support for checking whether a Helm Chart references resources not available in an OpenShift
cluster, resulting in scenarios where the Helm Chart installation will fail since Helm doesn't know which resources (as
specified by their Group, Version and Kind information) are expected in the cluster until the moment the rendered
resource lands in the cluster.

Verifying this specific information can help a developer to verify whether a Helm Chart can be installed on a specific
OpenShift release. Additionally, having this information at hand for future releases could also help users to verify
upcoming deprecations on the API level and also ways to migrate those resources.

## Rationale

Different from other checks, this particular verification depend on external information in order to validate the Helm
Chart for OpenShift compatibility issues, and some assumptions should be made considering the volatility of the data
being verified-both chart and OpenShift versions.

One of these assumptions is that **resource definitions present by default in a specific OpenShift version** will be
available by the program to verify the Helm Chart against; in other words, a dictionary containing all API resources,
groups and versions **distributed by default** will either be bundled within the resulting binary or available to be
downloaded on demand for each OpenShift version.

The availability of this collection of dictionaries can provide support for existing OpenShift versions as well
**upcoming** releases, so in the case deprecation notices are present in the API server those could also be used to
provide guidance to the Helm Chart Developer, informing which resources are deprecated and how to use more appropriate
resources if possible.

There are at least two approaches regarding the layout of the dictionaries' values mentioned above:

1. A list of GVKs that are present in a default OpenShift installation, where rendered resources from a Helm chart can 
   be matched against values present in this list, either strictly (all rendered resources must exist in this list) or
   other heuristics enabled by the available data; in other words, a Helm chart must contain only rendered resources
   present in this list.
   
   This is important to bear in mind that lists could be appended, so a chart could be checked against two different 
   OpenShift versions.
   
   Another interesting point in this approach is those GVK lists also could represent subsystems, such as containing
   only Knative APIGroup; this means that at least in theory, the Knative set together with the OpenShift set could
   be enough to statically validate a Helm chart against a set of expected resources available in the cluster.
   
1. The OpenAPI V3 Schema from the default OpenShift installation, where rendered resources from a Helm chart can be 
   matched against the available schema.
   
   Kubernetes and OpenShift expose through a REST interface the resources available in the cluster (which is also used
   to collect the resources in this page: https://docs.openshift.com/container-platform/4.6/rest_api/objects/index.html)
   .

   This approach offers less ways to overlay default cluster API resources with specific API resources, unlike the 
   former, basically due the fact that a well formed OpenAPI V3 Schema requires its dependencies to be declared; in
   reality this means that a Knative OpenAPI V3 Schema would include, other than specific Knative resources, Kubernetes 
   resources it depends on.

   It should be observed that the validation of rendered resources using the OpenAPI V3 Schema available from OpenShift
   would be the pretty similar if not the same same as performed by `kubectl`.

Both approaches use the same authoritative data source, which makes the choice of verification a matter of the desired
resolution, where a collection of GVKs provide only a low resolution verification and using the OpenAPI V3 Schema 
offers a higher resolution verification.

Both approaches would also use external (or embedded) artifacts containing verification data, so the verification
could be performed for past, current and upcoming OpenShift or Kubernetes versions, increasing the tool's usefulness.

# Command Line Interface

It is expected the command line interface to support the `--set` option, used to propagate configurations to the checks
at runtime in the following format: `--set KEY=VALUE`, where `KEY` represents a field path, such as `compat.version`, 
in a data-structure representing the configuration, and `VALUE` represents the value the specified configuration
key should have, such as `openshift-4.6`.

It is also expected the configuration file to have the check internal name as keys, and dictionaries as values, 
representing the entire check configuration; those values will be interpreted by each individual check, since they're
opaque to the main process.

## Command Line Options

The version to be matched against should be informed via the command line interface using the `--set` flag:

```text
> chart-verifier verify --enable compat --set compat.version=openshift-4.6 chart.tgz
```

To validate against multiple versions, a comma separated list could be used:

```text
> chart-verifier verify --enable compat --set compat.version=openshift-4.6,openshift-4.7 chart.tgz
```

## Configuration

The same settings stated above could be materialized in a configuration file, such as below:

```yaml
compat:
  version: openshift-4.6,openshift-4.7
```

Then used as:

```text
> chart-verifier verify --config cv.yaml --enable compat chart.tgz
```

# References

1. `oc proxy & curl localhost:8001/openapi/v2 > openapi.json`
1. https://stackoverflow.com/a/48804996
1. https://github.com/go-openapi