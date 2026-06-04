# builder
FROM golang:1.26.3-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

# ソース全部コピー
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server ./cmd/server

# -- Stage2: runtime ----
FROM alpine:3.20

# CA証明書(HTTPS通信、TLS検証に必要)
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/server /app/server

EXPOSE 8080

# 非 root実行(セキュリティ)
RUN adduser -D -u 10001 appuser
USER appuser

CMD ["/app/server"]