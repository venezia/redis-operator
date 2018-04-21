# Initial preparations
FROM golang:1.10.1-alpine
#RUN apk update && apk add git

# Let's add the code in
ADD . /go/src/gitlab.com/mvenezia/redis-operator

# Time to intall and rund ep
#RUN go get -u github.com/golang/dep/cmd/dep
WORKDIR /go/src/gitlab.com/mvenezia/redis-operator
#RUN dep ensure

# Build package in container
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o ./redis-operator cmd/redis-operator/main.go

# Now that we have build a static binary, we can stuff it in a bare-bones alpine image
FROM alpine
COPY --from=0 /go/src/gitlab.com/mvenezia/redis-operator/redis-operator /redis-operator

ENTRYPOINT /redis-operator