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

The availability of this collection of dictionaries can provide support for existing OpenShift versions as well **
upcoming** releases, so in the case deprecation notices are present in the API server those could also be used to
provide guidance to the Helm Chart Developer, informing which resources are deprecated and how to use more appropriate
resources if possible.

