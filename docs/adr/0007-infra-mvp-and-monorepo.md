# ADR 0007: インフラ MVP 構成と Terraform を monorepo に置く判断

## ステータス
採用 (2026-06-03) / 更新 (2026-06-06): FE を Vercel に移行、 Bastion EC2 と GitHub Actions OIDC を追加

## コンテキスト

musclead は個人開発の SaaS 練習プロジェクト。 ローカル開発
(Backend Go + FE Next.js + MySQL on Docker) で機能実装は概ね完了したので、
AWS に **本番に近い構成**でデプロイし、 「設計から運用まで一気通貫で
触れる」 経験を積むのが本 ADR の目的。

決定すべきは大きく 3 つ:

1. **どの AWS サービスで動かすか**(コストと学習価値のバランス)
2. **Terraform コードをどこに置くか**(monorepo か別リポか)
3. **MVP の範囲をどこで切るか**(本番化はどこまでやらない)

## 決定

### 1. AWS / Vercel スタック (MVP 構成)

| 層 | サービス | 月額 (初年度) | 月額 (翌年) |
|---|---|---|---|
| FE hosting | Vercel Hobby (Next.js SSR) | $0 | $0 |
| Backend container | ECS Fargate Spot 0.25 vCPU / 0.5 GB (ARM64) | $3 | $3 |
| LB | ALB (target group は Backend のみ) | $18 | $18 |
| DB | RDS db.t4g.micro Single-AZ (free tier 利用) | $0 | $13 |
| Bastion | EC2 t4g.nano (普段は stopped、 必要時のみ起動) | ~$0 | ~$0 |
| Image registry | ECR Private (server リポジトリのみ) | $0.1 | $0.1 |
| Domain / Cert | Route 53 + ACM | $0.5 | $0.5 |
| Secrets | SSM Parameter Store SecureString | $0 | $0 |
| Logs | CloudWatch Logs 7 日保持 | $0.5 | $0.5 |
| CI/CD | GitHub Actions + OIDC Role (Actions 側無料枠内) | $0 | $0 |
| **合計** | | **~$22 (¥3,500)** | **~$35 (¥5,500)** |

ドメイン: `musclead.com` を Route 53 で取得 ($15/年)。

サブドメイン分離:
- `app.musclead.com` → **Vercel** (FE / Next.js SSR)
- `api.musclead.com` → ALB → ECS Fargate (Backend)

### 2. Terraform は monorepo (musclead リポジトリ内 `terraform/`)

```
musclead/
├── server/
├── web/
├── sql/
├── docs/
└── terraform/
    ├── bootstrap/        # tfstate 用 S3 + DynamoDB
    └── modules/
        ├── network/      # VPC / Subnet / IGW / SG
        ├── rds/          # RDS MySQL Single-AZ
        ├── ecr/          # ECR Private (server のみ)
        ├── ecs/          # ECS Cluster / Task / Service (ARM64)
        ├── alb/          # ALB / Listener / Target Group
        ├── acm/          # ACM 証明書
        ├── dns/          # Route 53 record
        ├── secrets/      # SSM Parameter Store
        ├── bastion/      # 踏み台 EC2 (SSM port forward 用)
        └── github_oidc/  # GitHub Actions OIDC Role
```

### 3. MVP で「やらないこと」 を明示

- Multi-AZ HA(Single-AZ)
- Auto Scaling(常時 1 タスク)
- WAF(ALB 直接公開)
- Sentry / 監視アラート(CloudWatch Logs のみ)
- Secrets Manager(SSM Parameter Store で代替)
- IAM 最小権限の細分化(Terraform 用に Administrator アタッチ、 後で絞る)
- VPC Endpoint(NAT 無し前提なので不要)
- Private Subnet + NAT Gateway(コスト優先で Public Subnet のみ)

## 理由

### a. なぜ Fargate Spot か

- On-Demand $8/月 → Spot $3/月 (約 70% 割引)
- 2 分前通知で interruption する可能性あるが、 個人デモ規模では実害ほぼ無し
- 「コスト最適化の引き出し」 として Spot を一度実運用するのが目的

### b. なぜ ALB + ECS Fargate (App Runner ではなく)

- ALB / Target Group / Listener Rule / Health Check 周りをマネージドサービスで一通り触っておきたい
- App Runner だとこれらが隠蔽されて学べない
- +$15/月の追加コストは学習用途として許容

### c. なぜ FE を Vercel に変更したか (Fargate ではなく)

当初は FE も ECS Fargate Spot で動かす計画だったが、 以下の理由で Vercel に変更:

- **Next.js は Vercel で動かすのが純正運用** で、 SSR / ISR / Edge Functions / Image Optimization まで Zero-config で揃う
- Hobby tier で個人 SaaS 規模は完全に無料枠内 → 月 $3 削減
- Fargate task の deploy / rollback / preview 環境を自前で組むよりも、 PR ごとの Preview URL がそのまま使える
- 「Backend は AWS、 FE は Vercel」 の hybrid 構成は実務でもよくある形なので学習価値も保てる
- `output: "export"` で S3 静的化する案は、 既存コードに動的セグメント
  (`/routines/[id]/edit` 等) が多数あるため不可

詳細は ADR 0008 (RSC 移行検討) も参照。

### d. なぜ Bastion EC2 (t4g.nano) を入れたか

- RDS は Single-AZ + Public Subnet 配置だが、 SG inbound は ECS task と
  bastion のみ → 直接外から接続できない
- ローカルから TablePlus 等で接続するために、 **SSM Session Manager の
  port forward** で bastion を中継して RDS に届ける
- t4g.nano $0.0042/hr で、 普段は stopped、 必要時のみ alias で start/stop
- bastion SG は **inbound なし** (SSM は outbound のみで成立)、 SSH key も不要
- 公開鍵 + 22 番ポート開放の伝統的 bastion より圧倒的に安全

### e. なぜ Public Subnet のみ (NAT 無し)

- NAT Gateway $32/月は MVP には重い
- Public Subnet + Security Group で締めれば実用上は問題なし
- ECS task は SG で「ALB からのみ inbound」 に絞る → 公開 IP 持っても外から直接叩けない
- 将来「本番化 ADR」 で Private Subnet + NAT に変更する予定

### f. なぜ monorepo (musclead リポ内 `terraform/`)

別リポを推す典型理由(複数チーム / 別 cadence / 権限分離 / 共有 infra)が
musclead では全く当てはまらない。 一方 monorepo のメリットは以下:

- **アプリと infra の同時変更が 1 PR で完結**
  例: 新 API endpoint 追加 → Backend handler + ALB listener rule を同じ PR で。
  別リポだと 2 PR 間の調整必要、 マージ順序事故のリスク。
- **コミット履歴で「なぜこの infra 変更か」 が見える**
  半年後に `alb.tf` の git blame で `feat(billing): add /api/billing` まで
  辿れる。 別リポだとクロスリポ検索が必要。
- **CI/CD パイプライン 1 個で済む**
  認証設定・secrets 管理が 1 箇所。 メンテ負荷半減。
- **個人プロジェクトでは管理対象を減らすのが正義**
  リポ切り替え・issue tracker・README 2 重メンテのコストは無視できない。

### g. なぜ初年度は free tier 前提か

- RDS db.t4g.micro: 12 ヶ月無料(月 750h までだが 1 インスタンス常時 = 744h、 範囲内)
- 月 $13 → $0 に削減
- 学習用途で `db.t3.micro` / `db.t4g.micro` は十分

### h. なぜ Secrets Manager ではなく SSM Parameter Store か

- Secrets Manager: $0.40/secret/月
- SSM Parameter Store Standard: 無料
- 個人開発で rotation 必須の secret はない
- 階層名 (`/musclead/prod/db-url`) で組織化しておけば、 本番では
  Secrets Manager に昇格しやすい設計

### i. なぜ GitHub Actions OIDC を採用したか(IAM ユーザーのアクセスキーではなく)

- 当初は「初回検証優先で CI/CD は後回し」 としていたが、 初回 apply が落ち着いた段階で OIDC を導入
- アクセスキー方式は **長期 credential が GitHub Secrets に残り続ける** ため、 漏洩時のリスクが大きい
- OIDC + IAM Role の AssumeRole は短期 credential のみで動き、 ローテーション不要
- `allowed_branch = "main"` で main ブランチからの実行のみに絞る
- ECR push と ECS Task Definition の更新権限のみ付与(最小権限)

## 不採用案

A. **AWS App Runner で Backend をデプロイ**
   - 棄却: ALB / ECS Service / Task Definition の経験が積めない。
     マネージドコンテナ運用の学習として ECS Fargate を選択。

B. **Next.js を S3 + CloudFront で静的配信**
   - 棄却: `output: "export"` × 動的セグメントでビルド失敗。
     既存コードの大幅書き換えが必要。

C. **FE も ECS Fargate Spot で運用 (当初案)**
   - 棄却: 上記 c. の通り Vercel に変更。 月 $3 削減 + deploy / preview の運用コスト削減。

D. **EC2 + docker compose**
   - 棄却: 「コンテナをマネージドサービスで運用する」 経験が積めない。
     コスト最安だが学習価値が薄い。

E. **Multi-AZ + WAF + Secrets Manager 最初から入れる**
   - 棄却: コスト $80/月超、 個人 SaaS 練習用途には過剰。
     段階的に強化する方針。

F. **Terraform コードを別リポ `musclead-infra` に分離**
   - 棄却: 別リポを正当化する要件(複数チーム / 別 cadence / 権限分離 /
     共有 infra)が個人 SaaS では存在しない。
     管理コスト > メリット。

G. **GitHub Actions に IAM ユーザーのアクセスキーを置く**
   - 棄却: 長期 credential 漏洩リスク。 OIDC + AssumeRole で短期 credential 化。

## 結果

- `terraform/` ディレクトリを musclead リポジトリ直下に作成済
- VPC / RDS / ECR / ECS / ALB / ACM / DNS / Secrets / Bastion / GitHub OIDC の 10 module 構成
- FE は Vercel (Hobby tier)、 Backend は ECS Fargate Spot (ARM64)
- 月額 ~¥3,500 (初年度) で動かす
- ドメイン `musclead.com`、 `app.musclead.com` → Vercel、 `api.musclead.com` → ALB
- 将来 ADR で本番化 (Multi-AZ / WAF / Secrets Manager / 観測性ツール) を計画

## 実施状況

### Phase 1: 手動準備 ✅
1. AWS アカウント + IAM ユーザー + MFA + 予算アラーム
2. Route 53 でドメイン取得
3. Terraform 用 IAM ユーザー + アクセスキー(OIDC 移行前の bootstrap 用)
4. tfstate 保存先 S3 バケット + DynamoDB ロックテーブル (`terraform/bootstrap/`)

### Phase 2: Terraform でインフラ構築 ✅
1. `modules/network/`: VPC / Subnet / IGW / Route Table / SG (ALB / Backend Fargate / RDS / Bastion)
2. `modules/rds/`: RDS MySQL + DB Subnet Group
3. `modules/ecr/`: ECR Private (server リポジトリのみ、 FE は Vercel)
4. `modules/secrets/`: SSM Parameter Store (jwt_secret / db_user / db_password / db_host)
5. `modules/acm/`: ACM 証明書 (`api.musclead.com`)
6. `modules/alb/`: ALB + Listener + Target Group
7. `modules/ecs/`: ECS Cluster + Task Definition (ARM64) + Service
8. `modules/dns/`: Route 53 record (`api` → ALB Alias、 `app` → Vercel A record)
9. `modules/bastion/`: 踏み台 EC2 + SSM port forward
10. `modules/github_oidc/`: GitHub Actions OIDC Provider + IAM Role

### Phase 3: アプリのデプロイ ✅
1. Dockerfile 作成 (Backend のみ、 FE は Vercel が build)
2. Docker build → ECR push (初回は手動、 以降は GitHub Actions)
3. ECS task が新イメージを pull → 稼働
4. RDS に手動マイグレーション実行 (`goose up`)
5. Vercel に web/ を連携、 `app.musclead.com` 割り当て
6. 動作確認 (`https://app.musclead.com` → `https://api.musclead.com`)

### Phase 4: 将来の本番化 ADR で順次強化
- Multi-AZ 化
- WAF
- Secrets Manager 昇格
- 観測性ツール (Sentry / OpenTelemetry / CloudWatch アラーム)
- IAM 最小権限の細分化
- Private Subnet + NAT (or VPC Endpoint)
