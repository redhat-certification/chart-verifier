#!/usr/bin/env bash

curl -O https://mirror.openshift.com/pub/openshift-v4/clients/helm/3.5.0/helm-linux-amd64.tar.gz
mkdir -p /usr/local/bin
tar xfz helm-linux-amd64.tar.gz -C /usr/local/bin
mv /usr/local/bin/helm-linux-amd64 /usr/local/bin/helm
