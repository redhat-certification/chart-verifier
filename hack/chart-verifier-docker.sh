#!/usr/bin/env bash

# KUBECONFIG is either the value the user has configured in the
# environment or the default location `kubectl` uses.
KUBECONFIG="${KUBECONFIG:-${HOME}/.kube/config}"

# CHART_VERIFIER_IMAGE is the location of the image to use to verify
# the chart.
CHART_VERIFIER_IMAGE="${CHART_VERIFIER_IMAGE:-quay.io/redhat-certification/chart-verifier:main}"

# DOCKER contains the path for the `podman` program.
DOCKER="${DOCKER:-/usr/bin/docker}"

"${DOCKER}" run --rm -i                      \
            -e KUBECONFIG=/.kube/config      \
            -v "${KUBECONFIG}":/.kube/config \
            "${CHART_VERIFIER_IMAGE}"        \
            $*

