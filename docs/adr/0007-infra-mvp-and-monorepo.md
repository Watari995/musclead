# ADR 0007: インフラ MVP 構成と Terraform を monorepo に置く判断

## ステータス
採用 (2026-06-03)

## コンテキスト

musclead は SODA 入社準備 (2026-07 入社) の練習プロジェクト。 ローカル開発
(BE Go + FE Next.js + MySQL on Docker) で機能実装は概ね完了したので、
AWS に **本番に近い構成**でデプロイし、 「設計から運用まで一気通貫で
触れる」 経験を積むのが本 ADR の目的。

決定すべきは大きく 3 つ:

1. **どの AWS サービスで動かすか**(コストと学習価値のバランス)
2. **Terraform コードをどこに置くか**(monorepo か別リポか)
3. **MVP の範囲をどこで切るか**(本番化はどこまでやらない)

## 決定

### 1. AWS スタック (MVP 構成)

| 層 | サービス | 月額 (初年度) | 月額 (翌年) |
|---|---|---|---|
| BE container | ECS Fargate Spot 0.25 vCPU / 0.5 GB | $3 | $3 |
| FE container | ECS Fargate Spot 0.25 vCPU / 0.5 GB (Next.js SSR) | $3 | $3 |
| LB | ALB (target group 2 つ、 host header 振り分け) | $18 | $18 |
| DB | RDS db.t4g.micro Single-AZ (free tier 利用) | $0 | $13 |
| Image registry | ECR Private | $0.1 | $0.1 |
| Domain / Cert | Route 53 + ACM | $0.5 | $0.5 |
| Secrets | SSM Parameter Store SecureString | $0 | $0 |
| Logs | CloudWatch Logs 7 日保持 | $0.5 | $0.5 |
| **合計** | | **~$25 (¥4,000)** | **~$38 (¥6,000)** |

ドメイン: `musclead.com` を Route 53 で取得 ($15/年)。

サブドメイン分離:
- `app.musclead.com` → FE Fargate task
- `api.musclead.com` → BE Fargate task

### 2. Terraform は monorepo (musclead リポジトリ内 `terraform/`)

```
musclead/
├── server/
├── web/
├── sql/
├── docs/
└── terraform/    ← 新規
```

### 3. MVP で「やらないこと」 を明示

- Multi-AZ HA(Single-AZ)
- Auto Scaling(常時 1 タスク)
- WAF(ALB 直接公開)
- Sentry / 監視アラート(CloudWatch Logs のみ)
- Secrets Manager(SSM Parameter Store で代替)
- CI/CD 自動デプロイ(初回は手動 `terraform apply` + `docker push`)
- IAM 最小権限の細分化(Terraform 用に Administrator アタッチ、 後で絞る)
- VPC Endpoint(NAT 無し前提なので不要)

## 理由

### a. なぜ Fargate Spot か

- On-Demand $8/月 → Spot $3/月 (約 70% 割引)
- 2 分前通知で interruption する可能性あるが、 個人デモ規模では実害ほぼ無し
- SODA 決済基盤で Spot は使わないだろうが、 「コスト最適化の引き出し」 を持つ意味は大きい

### b. なぜ ALB + ECS Fargate (App Runner ではなく)

- ALB / Target Group / Listener Rule / Health Check 周りは SODA 実務で必須スキル
- App Runner だとこれらが隠蔽されて学べない
- +$15/月の追加コストは 1 ヶ月限定で許容
- ADR で「なぜ ALB を選んだか」 を文書化することで面接アピール材料になる

### c. なぜ Next.js を Fargate (S3 静的化ではなく)

- 既存コードに動的セグメント (`/routines/[id]/edit` 等) が多数あり、
  `output: "export"` だとビルド失敗することを事前検証で確認
- 既存コード変更ゼロで動かすには SSR が必要 → Fargate task で動かす
- 詳細は [docs/notes/nextjs-static-export-investigation.md] (後日記載予定)
- AWS Amplify Hosting も候補だが、 ECS 統一で Terraform 管理を一元化したい

### d. なぜ Public Subnet のみ (NAT 無し)

- NAT Gateway $32/月は MVP には重い
- Public Subnet + Security Group で締めれば実用上は問題なし
- ECS task は SG で「ALB からのみ inbound」 に絞る → 公開 IP 持っても外から直接叩けない
- 将来 ADR 0008「本番化」 で Private Subnet + NAT に変更する予定

### e. なぜ monorepo (musclead リポ内 `terraform/`)

別リポを推す典型理由(複数チーム / 別 cadence / 権限分離 / 共有 infra)が
musclead では全く当てはまらない。 一方 monorepo のメリットは以下:

- **アプリと infra の同時変更が 1 PR で完結**
  例: 新 API endpoint 追加 → BE handler + ALB listener rule を同じ PR で。
  別リポだと 2 PR 間の調整必要、 マージ順序事故のリスク。
- **コミット履歴で「なぜこの infra 変更か」 が見える**
  半年後に `alb.tf` の git blame で `feat(billing): add /api/billing` まで
  辿れる。 別リポだとクロスリポ検索が必要。
- **CI/CD パイプライン 1 個で済む**
  認証設定・secrets 管理が 1 箇所。 メンテ負荷半減。
- **個人プロジェクトでは管理対象を減らすのが正義**
  リポ切り替え・issue tracker・README 2 重メンテのコストは無視できない。

### f. なぜ初年度は free tier 前提か

- RDS db.t4g.micro: 12 ヶ月無料(月 750h までだが 1 インスタンス常時 = 744h、 範囲内)
- 月 $13 → $0 に削減
- 学習用途で `db.t3.micro` / `db.t4g.micro` は十分

### g. なぜ Secrets Manager ではなく SSM Parameter Store か

- Secrets Manager: $0.40/secret/月
- SSM Parameter Store Standard: 無料
- 個人開発で rotation 必須の secret はない
- 階層名 (`/musclead/prod/db-url`) で組織化しておけば、 本番では
  Secrets Manager に昇格しやすい設計

### h. なぜ初回は GitHub Actions CD を入れないか

- 初回デプロイで「想定通り動くか」 の検証が先
- 本番並みの自動化は GitHub Actions OIDC + AWS IAM Role 連携が必要
  (準備自体が学習コスト)
- ADR 0008 で CI/CD 化を別途検討

## 不採用案

A. **AWS App Runner で BE をデプロイ**
   - 棄却: ALB / ECS Service / Task Definition の経験が積めない。
     SODA で必要なスキルが抜ける。

B. **Next.js を S3 + CloudFront で静的配信**
   - 棄却: `output: "export"` × 動的セグメントでビルド失敗。
     既存コードの大幅書き換えが必要。

C. **AWS Amplify Hosting で FE をデプロイ**
   - 棄却: Terraform で統一管理しづらく、 Amplify 独自の deploy 機構を学ぶ
     コストが musclead 文脈で割に合わない。

D. **EC2 + docker compose**
   - 棄却: 「コンテナをマネージドサービスで運用する」 経験が積めない。
     コスト最安だが学習価値が薄い。

E. **Multi-AZ + WAF + Secrets Manager 最初から入れる**
   - 棄却: コスト $80/月超、 1 ヶ月学習用途には過剰。
     ADR 0008 で段階的に強化する。

F. **Terraform コードを別リポ `musclead-infra` に分離**
   - 棄却: 別リポを正当化する要件(複数チーム / 別 cadence / 権限分離 /
     共有 infra)が個人 SaaS では存在しない。
     管理コスト > メリット。

## 結果

- `terraform/` ディレクトリを musclead リポジトリ直下に作成
- VPC / RDS / ECR / ECS / ALB / Route 53 を Terraform で 6 module 構成
- 初回デプロイは手動 `terraform apply` + `docker push`
- 月額 ¥4,000 (初年度) で動かす
- ドメイン `musclead.com` + サブドメイン `app` / `api`
- ADR 0008(将来) で本番化 (Multi-AZ / WAF / CI/CD / Secrets Manager) を計画

## 実施計画

### Phase 1: 手動準備 (Terraform 動かす前提)
1. AWS アカウント + IAM ユーザー + MFA + 予算アラーム(完了済)
2. Route 53 でドメイン取得(完了済)
3. Terraform 用 IAM ユーザーのアクセスキー作成(または OIDC 設定)
4. tfstate 保存先 S3 バケット + DynamoDB ロックテーブルを手動 or bootstrap module で作成

### Phase 2: Terraform でインフラ構築
1. `terraform/network/`: VPC / Subnet / IGW / Route Table / SG
2. `terraform/rds/`: RDS MySQL + DB Subnet Group
3. `terraform/ecr/`: ECR Private リポジトリ (BE / FE 用)
4. `terraform/ecs/`: ECS Cluster + Task Definition + Service
5. `terraform/alb/`: ALB + Listener + Target Group + ACM
6. `terraform/dns/`: Route 53 record (Alias for ALB)
7. `terraform/secrets/`: SSM Parameter Store

### Phase 3: アプリのデプロイ
1. Dockerfile 作成 (BE / FE)
2. Docker build → ECR push
3. ECS task が新イメージを pull → 稼働
4. RDS に手動マイグレーション実行
5. 動作確認 (`https://app.musclead.com` でアクセス)

### Phase 4: ADR 0008 を書いてから順次強化
- CI/CD 自動化
- Multi-AZ 化
- WAF / Secrets Manager / 観測性ツール
