# syntax = docker/dockerfile:1.0.0-experimental
# Build Container
FROM golang:1.16.1 as builder

ENV GO111MODULE on
ENV GOPRIVATE "bitbucket.org/latonaio"
WORKDIR /go/src/bitbucket.org/latonaio

COPY go.mod .

RUN git config --global url."git@bitbucket.org:".insteadOf "https://bitbucket.org/"
RUN mkdir /root/.ssh/ && touch /root/.ssh/known_hosts && ssh-keyscan -t rsa bitbucket.org >> /root/.ssh/known_hosts
RUN --mount=type=secret,id=ssh,target=/root/.ssh/id_rsa go mod download

COPY . .

RUN go build -o salesforce-api-kube .


# Runtime Container
FROM alpine:3.12

RUN apk add --no-cache libc6-compat tzdata

COPY --from=builder /go/src/bitbucket.org/latonaio/config.json .
COPY --from=builder /go/src/bitbucket.org/latonaio/config.test.json .
COPY --from=builder /go/src/bitbucket.org/latonaio/salesforce-api-kube .

CMD ["./salesforce-api-kube"]
