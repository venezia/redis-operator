# Initial preparations
FROM golang:1.10.1-alpine
RUN apk update && apk add git bash

RUN go get k8s.io/code-generator/cmd/client-gen && \
    go install k8s.io/code-generator/cmd/client-gen && \
    go install k8s.io/code-generator/cmd/conversion-gen && \
    go install k8s.io/code-generator/cmd/deepcopy-gen && \
    go install k8s.io/code-generator/cmd/defaulter-gen && \
    go install k8s.io/code-generator/cmd/go-to-protobuf && \
    go install k8s.io/code-generator/cmd/import-boss && \
    go install k8s.io/code-generator/cmd/informer-gen && \
    go install k8s.io/code-generator/cmd/lister-gen && \
    go install k8s.io/code-generator/cmd/openapi-gen && \
    go install k8s.io/code-generator/cmd/set-gen && \
    go get -u github.com/golang/dep/cmd/dep
