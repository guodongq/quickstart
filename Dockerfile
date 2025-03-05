FROM golang:1.24.1-alpine3.21 AS builder

WORKDIR /workpath

ENV GOPATH=/go
ENV GOCACHE=/go/caches/go-build

ADD ./libs /workpath/libs
ADD ./modules /workpath/modules
ADD ./go.mod /workpath/go.mod
ADD ./go.sum /workpath/go.sum

RUN CGO_ENABLED=0 go build ./modules/go-grpc-quickstart/main.go


FROM alpine:3.21

WORKDIR /app

COPY --from=builder /workpath/main ./quickstart

ENTRYPOINT ["/app/quickstart"]
