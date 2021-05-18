#!/usr/bin/env bash

# KUBECONFIG is either the value the user has configured in the
# environment or the default location `kubectl` uses.
KUBECONFIG="${KUBECONFIG:-${HOME}/.kube/config}"

# CHART_VERIFIER_IMAGE is the location of the image to use to verify
# the chart.
CHART_VERIFIER_IMAGE="${CHART_VERIFIER_IMAGE:-quay.io/redhat-certification/chart-verifier:main}"

# PODMAN contains the path for the `podman` program.
PODMAN="${PODMAN:-/usr/bin/podman}"

"${PODMAN}" run --rm -i                      \
            --log-level="error"              \
            --systemd="false"                \
            --cgroup-manager="cgroupfs"      \
            --cgroups="disabled"             \
            --events-backend="none"          \
            --log-driver="none"              \
            -e KUBECONFIG=/.kube/config      \
            -v "${KUBECONFIG}":/.kube/config \
            "${CHART_VERIFIER_IMAGE}"        \
            $*
