#!/usr/bin/env bash

CHART_VERIFIER_IMAGE="quay.io/redhat-certification/chart-verifier:main"
DOCKER="/usr/bin/docker"

"${DOCKER}" run --rm -i			\
       -e KUBECONFIG=/.kube/config	\
       -v "${HOME}/.kube":/.kube	\
       "${CHART_VERIFIER_IMAGE}"	\
       $*

