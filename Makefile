include .env
export

.PHONY: db-up db-down db-reset db-logs db-shell migrate-up migrate-down migrate-status migrate-create migrate-force

db-up: ## MySQL コンテナ起動
	docker compose up -d mysql
	@echo "waiting for mysql to be healthy..."
	@until [ "$$(docker inspect -f '{{.State.Health.Status}}' musclead-mysql 2>/dev/null)" = "healthy" ]; do sleep 1; done
	@echo "mysql is ready."

db-down: ## MySQL コンテナ停止(データは残す)
	docker compose down

db-reset: ## MySQL コンテナ停止 + ボリューム削除
	docker compose down -v

db-logs:
	docker compose logs -f mysql

db-shell: ## MySQL CLI に接続
	docker exec -it musclead-mysql mysql -u$(DB_USER) -p$(DB_PASSWORD) $(DB_NAME)

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
