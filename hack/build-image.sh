#!/usr/bin/env bash
#
# Copyright 2021 Red Hat
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

# CONTAINER_RUNTIME contains the program to use, either podman or
# docker. If not informed, will find either in the path.
for cr in podman docker; do
  [ -z "$CONTAINER_RUNTIME" ] && $cr --help >/dev/null 2>&1 && CONTAINER_RUNTIME=$cr
done;

COMMIT_ID=$(git rev-parse --short HEAD)

EXTRA_ARGS=""

case "$CONTAINER_RUNTIME" in
    *docker)
        EXTRA_ARGS="${EXTRA_ARGS} --progress=plain"
        ;;
esac

if [[ "$(uname -a)" =~ Linux.*WSL.*Linux && "$CONTAINER_RUNTIME" =~ podman ]]; then
    # systemd and journald are not available in WSL.
    EXTRA_ARGS="${EXTRA_ARGS} --log-level=error --cgroup-manager=cgroupfs --events-backend=none"
fi

"$CONTAINER_RUNTIME" build $EXTRA_ARGS -t quay.io/redhat-certification/chart-verifier:"$COMMIT_ID" .
