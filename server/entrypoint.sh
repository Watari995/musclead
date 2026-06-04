#!/bin/sh
# BE container entrypoint:
#   1. DB マイグレーションを最新まで適用
#   2. server 起動
#
# 失敗時(DB 未起動、 migration エラー等)は server を起動せずに exit。
set -eu

DB_URL="mysql://${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}?parseTime=true&multiStatements=true"

echo "[entrypoint] running migrations..."
migrate -path /app/migrations -database "${DB_URL}" up
echo "[entrypoint] migrations applied."

echo "[entrypoint] starting server..."
exec /app/server
