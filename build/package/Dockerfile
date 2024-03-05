FROM golang:1.21.5 AS builder
WORKDIR /go/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /go/bin/acadebot ./cmd/acadebot/

FROM alpine:3.14
COPY --chown=65534:65534 --from=builder /go/bin/acadebot .
USER 65534

ENTRYPOINT [ "./acadebot" ]