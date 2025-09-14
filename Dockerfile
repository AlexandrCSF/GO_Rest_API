FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache \
    go \
    gcc \
    g++ \
    musl-dev \
    pkgconfig \
    librdkafka-dev

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=1 go build -o apiserver ./cmd/apiserver

FROM alpine:latest

RUN apk add --no-cache ca-certificates librdkafka

WORKDIR /root/

COPY --from=builder /app/apiserver .

COPY --from=builder /app/configs ./configs

COPY --from=builder /app/static ./static

EXPOSE 8080

CMD ["./apiserver"]
