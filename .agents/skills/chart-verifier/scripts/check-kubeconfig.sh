#!/bin/sh

if [ -z "${KUBECONFIG:-}" ]; then
  echo "KUBECONFIG is not set"
  exit 0
fi

if [ ! -e "${KUBECONFIG}" ]; then
  echo "KUBECONFIG path is not accessible"
  exit 0
fi

echo "KUBECONFIG is set"
