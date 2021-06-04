#!/usr/bin/env bash

HELM_VERSION=${HELM_VERSION-3.5.0}
DOWNLOAD_DIR=${DOWNLOAD_DIR-/tmp/chart-verifier}
INSTALL_DIR=${INSTALL_DIR-/usr/local/bin}

mkdir -p "${DOWNLOAD_DIR}"
pushd "${DOWNLOAD_DIR}"
curl -O "https://mirror.openshift.com/pub/openshift-v4/clients/helm/${HELM_VERSION}/helm-linux-amd64.tar.gz
tar xfz helm-linux-amd64.tar.gz
mkdir -p "${INSTALL_DIR}"
mv "${DOWNLOAD_DIR}/helm-linux-amd64" "${INSTALL_DIR}/helm"
popd
