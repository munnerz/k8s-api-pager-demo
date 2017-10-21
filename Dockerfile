FROM golang:1.9-alpine as builder


ENV CGO_ENABLED 0
ENV LDFLAGS ""
WORKDIR $GOPATH/src/github.com/srossross/k8s-test-controller
COPY *.go ./
COPY pkg ./pkg
COPY vendor ./vendor
RUN apk update && apk add git

# FIXME: not sure why these are not in vendor
RUN go get github.com/imdario/mergo golang.org/x/crypto/ssh golang.org/x/sys/unix
RUN go build -ldflags "${LDFLAGS}" -o /test-controller ./main.go

FROM alpine
COPY --from=builder /test-controller  /test-controller
CMD ["/test-controller"]
