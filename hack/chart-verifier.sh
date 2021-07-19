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

VERIFY_ARGS=""

# http://abhipandey.com/2016/03/getopt-vs-getopts/
OPTSPEC="S:V:h-:"
while getopts "$OPTSPEC" optchar; do
    case "${optchar}" in
    h)
        echo "$0 [-V <openshiftVersion>] <chart>"
        exit 0
        ;;
    V)
        VALUE=${OPTARG#*=}
        OPENSHIFT_VERSION="--openshift-version=${VALUE}"
        ;;
    S)
        VALUE=${OPTARG#*=}
        FLAG=${OPTARG%=$VALUE}
        VERIFY_ARGS="${VERIFY_ARGS} -S $FLAG=$VALUE"
        ;;
    esac
done

# $OPTIND is the index of the first positional argument found
CHART_TO_VERIFY="${*:$OPTIND:1}"

[ -z "$CHART_TO_VERIFY" ] && {
    echo "Chart is missing, exiting."
    exit 1
}

case "${CHART_TO_VERIFY}" in
    http://*|https://*)
        ;;
    *)
        CHART_TO_VERIFY=${CHART_TO_VERIFY#file://}
        CHART_TO_VERIFY=$(realpath $CHART_TO_VERIFY)
        CHART_GUEST_BASEDIR="/charts"
        CHART_HOST_BASEDIR=$(dirname $CHART_TO_VERIFY)
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
