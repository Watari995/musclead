include .env
export

SWAG := $(shell go env GOPATH)/bin/swag

.PHONY: db-up db-down db-reset db-logs db-shell \
        migrate-up migrate-down migrate-status migrate-create migrate-force \
        run build test tidy swag-init swag-install help

help: ## ターゲット一覧を表示
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# ---- Docker / MySQL ----

db-up: ## MySQL コンテナ起動
	docker compose up -d mysql
	@echo "waiting for mysql to be healthy..."
	@until [ "$$(docker inspect -f '{{.State.Health.Status}}' musclead-mysql 2>/dev/null)" = "healthy" ]; do sleep 1; done
	@echo "mysql is ready."

db-down: ## MySQL コンテナ停止(データは残す)
	docker compose down

db-reset: ## MySQL コンテナ停止 + ボリューム削除
	docker compose down -v

db-logs: ## MySQL ログを tail
	docker compose logs -f mysql

db-shell: ## MySQL CLI に接続
	docker exec -it musclead-mysql mysql -u$(DB_USER) -p$(DB_PASSWORD) $(DB_NAME)

# ---- Migration (golang-migrate) ----

migrate-up: ## 全 migration を適用
	migrate -path $(MIGRATION_DIR) -database "$(DB_URL)" up

migrate-down: ## 1つロールバック
	migrate -path $(MIGRATION_DIR) -database "$(DB_URL)" down 1

migrate-status: ## 現在の version を表示
	migrate -path $(MIGRATION_DIR) -database "$(DB_URL)" version

migrate-create: ## name=xxx で新規 migration を作成
	migrate create -ext sql -dir $(MIGRATION_DIR) -seq $(name)

migrate-force: ## version=N で強制的に version を設定(dirty 状態の復旧用)
	migrate -path $(MIGRATION_DIR) -database "$(DB_URL)" force $(version)

# ---- Go (build / run / test) ----

run: ## サーバー起動(本番風、 リロードなし)
	cd server && go run ./cmd/server

dev: ## サーバー起動(air でホットリロード、 開発用)
	cd server && $(shell go env GOPATH)/bin/air

air-install: ## air をインストール
	GOBIN=$(shell go env GOPATH)/bin go install github.com/air-verse/air@latest

build: ## サーバービルド
	cd server && go build -o ../bin/server ./cmd/server

test: ## ユニットテスト
	cd server && go test ./...

tidy: ## go mod tidy
	cd server && go mod tidy

# ---- swag (OpenAPI) ----

swag-install: ## swag CLI をインストール
	GOBIN=$(shell go env GOPATH)/bin go install github.com/swaggo/swag/cmd/swag@latest

swag-init: ## OpenAPI ドキュメントを生成(server/docs に出力)
	cd server && $(SWAG) init \
		-g cmd/server/main.go \
		--parseInternal \
		--parseDependency \
		--output ./docs
