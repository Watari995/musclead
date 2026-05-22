# 💪 musclead

> 筋トレ・食事・体重を一元管理する個人向けSaaS

## 🏗️ 構成(Monorepo)

```
musclead/
├── server/      # Go BE (Connect-RPC + sqlc + MySQL)
├── web/         # React FE
├── mobile/      # Flutter(将来)
├── proto/       # Protobuf 定義(BE/FE/Mobile 共通)
├── sql/         # マイグレーション + sqlc クエリ
├── terraform/   # IaC(AWS)
├── docs/        # ドメインモデル / ADR
└── .github/workflows/  # CI/CD(path filter 分離)
```

## 🛠️ 技術スタック

| 領域 | 技術 |
|---|---|
| 言語 (BE) | Go 1.23+ |
| API | Connect-RPC |
| ORM | sqlc |
| DB | MySQL 8.0 / Aurora Serverless v2 |
| Architecture | DDD + Modular Monolith |
| Infra | AWS (ECS Fargate, ALB, S3, CloudFront) |
| IaC | Terraform |
| 監視 | Sentry + CloudWatch |
| FE | React + TypeScript + Connect-Web |
| CI/CD | GitHub Actions + AWS OIDC |

## 🚀 開発

```bash
# DB 起動
docker-compose up -d

# サーバー起動
cd server && go run ./cmd/server

# ヘルスチェック
curl http://localhost:8080/healthz
```

## 📚 ドキュメント

- [ドメインモデル](docs/domain-model.md)
- [ADR](docs/adr/)
