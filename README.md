# 💪 musclead

> 筋トレ・食事・体重を一元管理する個人向けSaaS

## 🏗️ 構成(Monorepo)

```
musclead/
├── server/      # Go バックエンド(net/http + gorp + MySQL)
├── web/         # React フロントエンド
├── mobile/      # Flutter(将来)
├── sql/         # マイグレーション (golang-migrate)
├── terraform/   # IaC(AWS)
├── docs/        # ドメインモデル / ADR
└── .github/workflows/  # GitHub Actions(vet / build / test)
```

## 🛠️ 技術スタック

| 領域 | 技術 |
|---|---|
| 言語 (BE) | Go 1.26+ |
| HTTP | net/http(Go 1.22 ServeMux)|
| ORM | gorp |
| DB | MySQL 8.0 / Aurora Serverless v2 |
| マイグレーション | golang-migrate |
| API ドキュメント | swag(OpenAPI 自動生成)|
| Architecture | DDD + Modular Monolith(strict) |
| Infra | AWS(ECS Fargate / ALB / S3 / CloudFront) |
| IaC | Terraform |
| 監視 | Sentry + CloudWatch |
| FE | React + TypeScript |
| CI/CD | GitHub Actions |
| テスト | testify(assert / mock) |

## 🚀 開発

```bash
# DB 起動
make db-up

# マイグレーション
make migrate-up

# サーバー起動(ホットリロード、 air)
make dev

# テスト
cd server && go test ./...

# ヘルスチェック
curl http://localhost:8080/health
```

## 📚 ドキュメント

- [ドメインモデル](docs/domain-model.md)
- [ADR](docs/adr/)
