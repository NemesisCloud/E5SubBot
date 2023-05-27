FROM golang:alpine as builder

WORKDIR /app

ENV GO111MODULE=on \
    CGO_ENABLED=0

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -ldflags '-w -s' -o E5SubBot .

RUN apk update && apk add --no-cache ca-certificates

COPY config.example.yml /app/build/config.yml

FROM alpine:latest

RUN apk add tzdata
COPY --from=builder /app/build/E5SubBot /

ENV DB_SLL_MODE true

ENTRYPOINT ["/E5SubBot"]
