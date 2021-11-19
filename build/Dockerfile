# syntax = docker/dockerfile:1.0.0-experimental
# Build Container
FROM golang:1.17.3 as builder

ENV GO111MODULE on
ENV GOPRIVATE "github.com/latonaio"
WORKDIR /go/src/github.com/latonaio

COPY go.mod .

RUN go mod download

COPY . .

RUN go build -o salesforce-api-kube .


# Runtime Container
FROM alpine:3.14

RUN apk add --no-cache libc6-compat tzdata

COPY --from=builder /go/src/github.com/latonaio/config.json .
COPY --from=builder /go/src/github.com/latonaio/config.test.json .
COPY --from=builder /go/src/github.com/latonaio/salesforce-api-kube .

CMD ["./salesforce-api-kube"]
