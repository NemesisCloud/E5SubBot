FROM golang:alpine as builder

WORKDIR /app

ENV GO111MODULE=on \
    CGO_ENABLED=0

# cache
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . .
RUN go build -ldflags '-w -s' -o E5SubBot .

RUN apk update && apk add --no-cache ca-certificates

RUN mkdir build && cp E5SubBot build && mv config.yml build/config.yml

FROM alpine:latest

RUN apk add tzdata
COPY --from=builder /app/build /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENV DB_SLL_MODE true

ENTRYPOINT ["/E5SubBot"]