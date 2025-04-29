####################################################################################################
# Quickstart Build stage which performs the actual build of Quickstart binaries
####################################################################################################
FROM docker.io/library/golang:1.24.2 AS quickstart-build

WORKDIR /go/src/github.com/guodongq/quickstart

COPY go.* ./
RUN go mod download

# Perform the build
COPY . .
RUN make quickstart-all

####################################################################################################
# Final image
####################################################################################################
FROM ubuntu

RUN apt-get update && apt-get install -y tini

COPY --from=quickstart-build /go/src/github.com/guodongq/quickstart/dist/quickstart* /usr/local/bin/

RUN ln -s /usr/local/bin/quickstart /usr/local/bin/go-grpc && \
    ln -s /usr/local/bin/quickstart /usr/local/bin/cli

ENTRYPOINT ["/usr/bin/tini", "--"]
