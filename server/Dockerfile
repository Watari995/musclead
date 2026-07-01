# Build context は repo root を想定:
#   docker buildx build -f server/Dockerfile <repo-root>
#
# 理由: sql/migrations を image に含める必要があるため(server/ より上の階層)

# builder: Go server
# --platform=$BUILDPLATFORM = ホスト arch で動かす(QEMU 不要)
# ARG TARGETARCH = buildx が自動で "arm64" 等を渡してくれる
# Go のクロスコンパイル機能で ARM 用バイナリを x86 上で生成 → 3-5 倍高速
FROM --platform=$BUILDPLATFORM golang:1.26.3-alpine AS builder
ARG TARGETOS
ARG TARGETARCH

WORKDIR /build

COPY server/go.mod server/go.sum ./
RUN go mod download

COPY server/ ./

# Swagger docs を build 前に生成(server/docs/ は .gitignore で git 管理外)
# main.go が `_ "github.com/Watari995/musclead/docs"` を import してるため必須
RUN go install github.com/swaggo/swag/cmd/swag@latest \
    && /go/bin/swag init -g cmd/server/main.go -o docs

# GOOS / GOARCH で「ターゲット用バイナリ吐いて」 と Go に指示
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -ldflags="-s -w" -o /app/server ./cmd/server
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -ldflags="-s -w" -o /app/goal-checker ./cmd/goal-checker

# -- Stage 2: migrate CLI ---------
# curl で arm64 バイナリを download するだけなので Go 不要
# --platform=$BUILDPLATFORM でホスト arch 実行 → ネットワーク I/O のみ、 高速
FROM --platform=$BUILDPLATFORM alpine:3.20 AS migrate
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
COPY --from=builder /app/goal-checker /app/goal-checker
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
