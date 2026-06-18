# musclead infrastructure (Terraform)

musclead を AWS にデプロイするための Terraform プロジェクト。 設計判断の
背景は [ADR 0007](../docs/adr/0007-infra-mvp-and-monorepo.md) を参照

## 構成サマリ

| 層 | サービス | モジュール |
|---|---|---|
| ネットワーク | VPC / Subnet / SG | `modules/network/` |
| データベース | RDS MySQL Single-AZ | `modules/rds/` |
| コンテナレジストリ | ECR Private | `modules/ecr/` |
| 実行環境 | ECS Fargate Spot (Backend / FE) | `modules/ecs/` |
| ロードバランサー | ALB + ACM | `modules/alb/` |
| DNS | Route 53 record | `modules/dns/` |
| シークレット | SSM Parameter Store | `modules/secrets/` |

## 前提

- Terraform 1.9.8 (`.tool-versions` で固定)
- AWS CLI 2.x
- AWS プロファイル `musclead-admin` 設定済
- tfstate 用 S3 バケット + DynamoDB ロックテーブル作成済

## tfstate backend

S3 + DynamoDB lock パターン。 名前は `backend.tf` を参照。

## ワークフロー

```bash
# 初期化
terraform init

# 差分確認
terraform plan

# 適用
terraform apply
```

## 注意

- リソース削除は `terraform destroy` で確実に終了させること
  (放置で課金が積み上がる)
- 月額アラーム($30) は別途設定済
