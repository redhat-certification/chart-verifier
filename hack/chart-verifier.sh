#!/usr/bin/env bash

# KUBECONFIG is either the value the user has configured in the
# environment or the default location `kubectl` uses.
KUBECONFIG="${KUBECONFIG:-${HOME}/.kube/config}"

# CHART_VERIFIER_IMAGE is the location of the image to use to verify
# the chart.
CHART_VERIFIER_IMAGE="${CHART_VERIFIER_IMAGE:-quay.io/redhat-certification/chart-verifier:latest}"

# CONTAINER_RUNTIME contains the program to use, either podman or
# docker. If not informed, will find either in the path.
for cr in podman docker; do
  [ -z "$CONTAINER_RUNTIME" ] && $cr --help >/dev/null 2>&1 && CONTAINER_RUNTIME=$cr
done;

# EXTRA_ARGS contains extra arguments to the container runtime.
EXTRA_ARGS=""

if [[ "$(uname -a)" =~ Linux.*WSL.*Linux && "$CONTAINER_RUNTIME" =~ podman ]]; then
    # systemd and journald are not available in WSL.
    EXTRA_ARGS="${EXTRA_ARGS} --log-level=error --systemd=false --cgroup-manager=cgroupfs --cgroups=disabled --events-backend=none --log-driver=none"
fi

# "[-o <openshiftVersion] <chart>"

PARSED_ARGUMENTS=$(getopt -a -o "S:V:h" --long "chart-set:,openshift-version:,help" -- "$@" 2>/dev/null)
[ $? -eq 0 ] || {
    PARSED_ARGUMENTS="--help"
}

eval set -- "${PARSED_ARGUMENTS}"

VERIFY_ARGS=""

while true; do
    case "$1" in
        -h|--help)
            echo "$0 [-V|--openshift-version <openshiftVersion>] <chart>"
            exit 0
            ;;
        -V|--openshift-version)
            shift
            V=$1
            shift
            OPENSHIFT_VERSION="--openshift-version=${V}"
            ;;
        -S|--chart-set)
            FLAG=$1
            shift
            VALUE=$1
            shift
            VERIFY_ARGS="${VERIFY_ARGS} $FLAG=$VALUE"
            ;;
        --)
            shift
            break
            ;;
    esac
done

[ -z "$1" ] && {
    echo "Chart is missing, exiting."
    exit 1
}

case "$1" in
    http://*|https://*)
        CHART_TO_VERIFY=$1
        ;;
    *)
        CHART_TO_VERIFY=${1#file://}
        CHART_GUEST_BASEDIR="/charts"
        CHART_HOST_BASEDIR=$(dirname $(realpath "$CHART_TO_VERIFY"))
        CHART_NAME=$(basename "$CHART_TO_VERIFY")
        CHART_TO_VERIFY="${CHART_GUEST_BASEDIR}/${CHART_NAME}"

        EXTRA_ARGS="${EXTRA_ARGS} -v ${CHART_HOST_BASEDIR}/${CHART_NAME}:${CHART_GUEST_BASEDIR}/${CHART_NAME}:z"
        ;;
esac

# Execute the command.
$CONTAINER_RUNTIME                     \
    run --rm -i                        \
    -e KUBECONFIG=/.kube/config        \
    $EXTRA_ARGS                        \
    -v "${KUBECONFIG}":/.kube/config:z \
    "${CHART_VERIFIER_IMAGE}"          \
    verify ${OPENSHIFT_VERSION} ${VERIFY_ARGS} ${CHART_TO_VERIFY}
