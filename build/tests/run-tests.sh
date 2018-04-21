#!/usr/bin/env bash

PACKAGE_NAME=gitlab.com/mvenezia/redis-operator

go test -race \
  $(go list ${PACKAGE_NAME}/... | grep -v /pkg/client/ ) \
  -v -coverprofile .testCoverage.txt


