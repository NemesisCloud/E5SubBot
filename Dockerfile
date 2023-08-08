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

FROM alpine:latest

RUN apk add tzdata

COPY --from=builder /app/E5SubBot /E5SubBot
COPY . .
ENV DB_SLL_MODE true

ENTRYPOINT ["/E5SubBot"]
