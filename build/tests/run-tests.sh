#!/usr/bin/env bash

THIS_DIRECTORY="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PACKAGE_HOME=${THIS_DIRECTORY}/../../

cd $PACKAGE_HOME
go test -race \
  $(go list ./... | grep -v /pkg/client/ ) \
  -v -coverprofile .testCoverage.txt


