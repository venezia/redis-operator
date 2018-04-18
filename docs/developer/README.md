# Developer Information

## Build `redis-operator` Instructions

### As a Developer

Once this project has been checked out, one will notice that there is no
vendor folder.  Please use to fetch the required libraries.  In case you're
not familiar with `dep ensure`, here is a quick example:

```shell
compy5000:redis-operator venezia$ rm -rf vendor/
compy5000:redis-operator venezia$ go get -u github.com/golang/dep/cmd/dep
compy5000:redis-operator venezia$ dep ensure
compy5000:redis-operator venezia$ du -s -m vendor/
27	vendor/
```

### As a CI/CD process

There is an included `Dockerfile` that will build the project for you located
in [build/docker](../../build/docker/redis-operator)

To execute it, simply do something like

```shell
$ docker build -f ./build/docker/redis-operator/Dockerfile -t quay.io/venezia/redis-operator:v0.0.1 . 
```

## Required tools

### dep

[golang's dep](https://golang.github.io/dep/) is used to maintain the vendor folder.  If you intend to add
additional libraries to this project, please update the `Gopkg.toml` file
accordingly (or have dep do it for you)

This document doesn't intend to replicate dep's documentation - please see
dep's documentation for proper uses of dep.

### Kubernetes' code-generator

[Kubernetes' code-generator](https://github.com/kubernetes/code-generator) is
used by this project to parse the [api objects](../../pkg/apis/redis) and
convert them to a functional [client library](../../pkg/client).

Because the `code-generator` code is not directly used by the program, the
library is not available in the `vendor` folder.  Instead one needs to fetch
this library manually.

To do so, please do the following:

```shell
/go # go get k8s.io/code-generator
package k8s.io/code-generator: no Go files in /go/src/k8s.io/code-generator
/go # ls -l /go/src/k8s.io/code-generator/
```

Note that the error _no Go files in ..._ is expected, as the code generator
is not actually a go library, but rather just shell scripts.

Once that is working, you can then use [hack/update-codegen.sh](../../hack/update-codegen.sh)
which will look at the api folder and update the client library if need be.

As time permits, this usage will be cleaned up with the use of `client-gen`
and `openapi-gen` instead of the shell scripts.

