#!/usr/bin/env bash

# Copyright (c) 2019 StackRox Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# deploy.sh
#
# Sets up the environment for the admission controller webhook demo in the active cluster.

set -euo pipefail

basedir="$(dirname "$0")/deployment"
keydir="$(mktemp -d)"

# Generate keys into a temporary directory.
echo "Generating TLS keys ..."
"${basedir}/generate-keys.sh" "$keydir"

# Create the `webhook-demo` namespace. This cannot be part of the YAML file as we first need to create the TLS secret,
# which would fail otherwise.
echo "Creating Kubernetes objects ..."
kubectl get ns -oname | grep webhook-demo >/dev/null 2>&1 || kubectl create ns webhook-demo

# Create the TLS secret for the generated keys.
kubectl create secret tls webhook-server-tls \
    --cert "${keydir}/webhook-server-tls.crt" \
    --key "${keydir}/webhook-server-tls.key" \
    --dry-run=client -o yaml | kubectl apply -n webhook-demo -f -

# Read the PEM-encoded CA certificate, base64 encode it, and replace the `${CA_PEM_B64}` placeholder in the YAML
# template with it. Then, create the Kubernetes resources.
CA_PEM_B64=$(openssl base64 -A <"${keydir}/ca.crt") CHECKSUM=$(openssl md5 md5sum "${keydir}/webhook-server-tls.crt" | gawk '{print $NF}') \
    envsubst <"${basedir}/deployment.yaml.template" | kubectl apply -n webhook-demo -f -

# Delete the key directory to prevent abuse (DO NOT USE THESE KEYS ANYWHERE ELSE).
rm -rf "$keydir"

echo "The webhook server has been deployed and configured!"
