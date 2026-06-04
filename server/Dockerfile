# Build context は repo root を想定:
#   docker buildx build -f server/Dockerfile <repo-root>
#
# 理由: sql/migrations を image に含める必要があるため(server/ より上の階層)

# builder: Go server
FROM golang:1.26.3-alpine AS builder

WORKDIR /build

COPY server/go.mod server/go.sum ./
RUN go mod download

COPY server/ ./

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server ./cmd/server

# -- Stage 2: migrate CLI ---------
# golang-migrate を別 stage で取得(production image に含める)
FROM golang:1.26.3-alpine AS migrate
ARG TARGETARCH
RUN apk add --no-cache curl tar \
    && curl -L "https://github.com/golang-migrate/migrate/releases/download/v4.18.1/migrate.linux-${TARGETARCH}.tar.gz" \
       | tar xz -C /usr/local/bin

# -- Stage 3: runtime ---------
FROM alpine:3.20

# CA証明書(HTTPS通信、TLS検証に必要)
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/server /app/server
COPY --from=migrate /usr/local/bin/migrate /usr/local/bin/migrate

# Migration ファイル同梱(起動時に適用)
COPY sql/migrations /app/migrations

# Entrypoint: migrate → server 起動
COPY server/entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

EXPOSE 8080

# 非 root実行(セキュリティ)
RUN adduser -D -u 10001 appuser
USER appuser

ENTRYPOINT ["/app/entrypoint.sh"]
