#!/usr/bin/env bash

CHART_VERIFIER_IMAGE="quay.io/redhat-certification/chart-verifier:main"
PODMAN="/usr/bin/podman"

"${PODMAN}" run --rm -i			\
       -e KUBECONFIG=/.kube/config	\
       -v "${HOME}/.kube":/.kube	\
       "${CHART_VERIFIER_IMAGE}"	\
       $*
