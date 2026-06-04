# Day 11 — インフラ大構築 Day(ACM / ALB / DNS / CORS)

## 🌐 ネットワーク基礎

- **ARN**: Amazon Resource Name の略。 AWS リソースを一意に識別する文字列
- **VPC**: 自分専用の仮想ネットワーク
- **CIDR**: `10.0.0.0/24` の `/24` がネットワーク部のビット数を示す記法。
  数字が大きいほど範囲が狭い(`/24` = 256 IP、 `/16` = 65,536 IP)
- **HTTPS のポート**: 443。 `curl -v`(verbose) でリクエスト/レスポンスのヘッダーを確認できる
- **alpine**: Linux distro の一つ。 軽量で Docker base image として頻用される。
  最終的にアプリは Linux OS の上で動くため、 alpine ベースでビルドする

## 🛠 Terraform 基本

- module ごとに設定を分けるのが基本
- **環境ごとに変わる / 機密性の高い値**: `variable` 経由でエントリポイントから受け取る
- **環境ごとに変わらない値**: `module/main.tf` に直書き
- root の `main.tf` で `module "xx" {}` を書き、 module に必要な変数を注入
- **`terraform.tfvars`**: `terraform apply/plan` 時に自動で読み込まれる変数値ファイル
- **module 間で動的に決まる値**(SG ID, Subnet ID 等):
  - 出力側 module の `outputs.tf` で公開
  - 呼び出し側で `module.<名前>.<output 名>` として参照
- **新 module を `main.tf` に追加した時** → `terraform init` が必要
  (module の中身を変えただけなら init 不要)

## 🐳 ECS の構造

> **クラスタ**(箱)の中で、 **サービス**(管理人)が **タスク定義**(レシピ)を元に
> **タスク**(コンテナ実体)を常時 N 個動かす。

| 概念 | 役割 |
|---|---|
| Cluster | リソースの入れ物 |
| Task Definition | コンテナ起動の仕様書(image, env, port 等) |
| Task | Task Definition から生成された実体(Docker container 相当) |
| Service | Task を desired_count 個維持する管理者 |

## 🔄 ALB のヘルスチェック

ALB は死んでいる Task にはリクエストを送らないため、 定期的にヘルスチェックを実施。
musclead では `/health` を 30 秒間隔で叩き、 2 回連続成功で healthy、 3 回連続失敗で unhealthy。

## 🔌 ENI(Elastic Network Interface)

- VPC 内で IP アドレス通信を可能にする仕組み
- クラウド版の NIC(物理 PC の LAN ポート / Wi-Fi アダプタに相当)
- 「**ドアのようなもの**」 で、 AWS サービス間の接続を支える
- SG(Security Group)で適切に管理しないと接続できない

## 🌍 DNS レコード

| レコード | 用途 | musclead での例 |
|---|---|---|
| **A (Alias)** | AWS リソースに紐付け | `api.musclead.com` → ALB |
| **CNAME** | 外部ドメインに紐付け(後に A に変更) | `app.musclead.com` → Vercel |
| **A** | IP 直指定 | `app.musclead.com` → Vercel anycast IP `76.76.21.21` |

### CNAME の制約
- **root ドメインに CNAME は使えない**(RFC 仕様、 SOA/NS と共存禁止)
- サブドメインなら CNAME OK
- → root を AWS リソースに向けたい時は **Alias(独自仕様)** を使う

## 🎯 今日完成した本番ルート

```
ユーザー
  ↓ HTTPS
app.musclead.com (Route 53 A → Vercel 76.76.21.21)
  ↓ NEXT_PUBLIC_API_BASE_URL
api.musclead.com (Route 53 A Alias → ALB)
  ↓ Listener 443
ALB (ACM 証明書で TLS 終端)
  ↓ Target Group (HTTP 8080)
ECS Service → Task (Fargate Spot ARM64)
  ↓
BE Go (ALLOWED_ORIGIN=https://app.musclead.com で CORS 許可)
  ↓
RDS MySQL
```
