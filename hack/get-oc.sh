#!/usr/bin/env bash

curl -O https://mirror.openshift.com/pub/openshift-v4/x86_64/clients/ocp/4.7.9/openshift-client-linux-4.7.9.tar.gz
mkdir -p /usr/local/bin
tar xfz  openshift-client-linux-4.7.9.tar.gz -C /usr/local/bin
