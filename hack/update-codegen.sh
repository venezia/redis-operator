#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname ${BASH_SOURCE})/..
CODEGEN_PKG=./../../../../../../..${GOPATH}/src/k8s.io/code-generator

${CODEGEN_PKG}/generate-groups.sh "deepcopy,client,informer,lister" \
  gitlab.com/mvenezia/redis-operator/pkg/client \
  gitlab.com/mvenezia/redis-operator/pkg/apis \
"redis:v1alpha1" \
--output-base "$(dirname ${BASH_SOURCE})/../../../../" \
--go-header-file ${SCRIPT_ROOT}/hack/custom-boilerplate.go.txt

