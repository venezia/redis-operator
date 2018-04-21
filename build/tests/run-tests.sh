#!/usr/bin/env bash

PACKAGE_NAME=gitlab.com/mvenezia/redis-operator

go test -race \
  ${PACKAGE_NAME}/cmd/redis-operator \
  ${PACKAGE_NAME}/pkg/controller \
  ${PACKAGE_NAME}/pkg/redis \
  ${PACKAGE_NAME}/pkg/util/k8sutil \
  ${PACKAGE_NAME}/pkg/util/retryutil \
  -v -coverprofile .testCoverage.txt


