# Lessons 0001: Terraform SG リファクタと ElastiCache 初回 apply で学んだこと

> 2026-06-08 Phase 2 Cache Foundation の作業中に遭遇した問題と解決策。
> 同じハマり方を将来繰り返さないための記録。

---

## 1. Security Group の inline と standalone は混在禁止

### 何が起きたか

```hcl
# network/main.tf
resource "aws_security_group" "rds" {
  ingress { ... }   # ← inline ingress
}

# bastion/main.tf
resource "aws_security_group_rule" "bastion_to_rds" { ... }   # ← standalone
```

同じ `rds` SG に対して **inline ingress と standalone rule の両方**を使っていた。
Terraform は `inline ingress` を「この SG の rule の完全な真実」 と扱うため、
**plan のたびに standalone で追加された bastion rule を destroy しようとする**。

### Terraform 公式の警告

> Terraform currently provides both a standalone `aws_security_group_rule`, and a
> Security Group resource with `ingress` and `egress` rules defined in-line.
> **At this time you cannot use a Security Group with in-line rules in conjunction
> with any Security Group Rule resources.**

### 解決策

**1 SG = 片方の方式に統一**する。 musclead は cross-module で rule 追加するため
**全て standalone (`aws_security_group_rule`) に統一**した。

```hcl
# SG 本体は ingress/egress を持たない
resource "aws_security_group" "rds" {
  name   = "musclead-rds-sg"
  vpc_id = aws_vpc.main.id
  tags   = { Name = "musclead-rds-sg" }
}

# Rule は別 resource
resource "aws_security_group_rule" "rds_mysql_in" {
  type                     = "ingress"
  from_port                = 3306
  to_port                  = 3306
  protocol                 = "tcp"
  security_group_id        = aws_security_group.rds.id
  source_security_group_id = aws_security_group.server_fargate.id
}
```

### 教訓

- module 跨ぎで rule を追加するなら **standalone 一択**
- musclead 内では SG を必ず network module で定義、 rule も standalone で書く

---

## 2. inline → standalone 移行で apply 中に Duplicate エラー

### 何が起きたか

inline rule を削除して standalone rule を追加するコード変更を 1 つの apply で実行したら、
AWS API で `InvalidPermission.Duplicate` エラー:

```
A duplicate Security Group rule was found on (sg-xxxxx).
operation error EC2: AuthorizeSecurityGroupIngress, ... already exists
```

### 原因

Terraform は 1 つの apply 内で:
1. SG resource の inline rule を「削除」 (SG resource の update)
2. standalone rule resource を「作成」 (新規 resource の create)

を **並列実行** する。 AWS は「同じ内容の rule は重複できない」 ため、
standalone rule の create 時にまだ inline rule が AWS 側に残っていて duplicate になる。

これは Terraform 公式に「This may be a side effect of a now-fixed Terraform issue」
と書かれているが、 実際にはまだ起きる。

### 解決策

**AWS CLI で先に AWS 側の inline rule を手動 revoke** → `terraform apply` で復元。

```bash
# 重複している rule を AWS から削除
aws --profile musclead-admin --region ap-northeast-1 \
  ec2 revoke-security-group-ingress \
  --group-id sg-xxxxx \
  --ip-permissions 'IpProtocol=tcp,FromPort=3306,ToPort=3306,UserIdGroupPairs=[{GroupId=sg-yyyyy}]'

# 即 terraform apply で standalone rule を作成 (数秒の通信断)
terraform apply
```

revoke → apply の間は通信が切れるので、 サービス稼働中なら **短時間で連続実行**する。

### 教訓

- inline → standalone のような「state vs AWS の食い違いが起きる移行」 は、
  事前に「duplicate になりそうか」 を読み取って AWS CLI 介入を計画する
- 学習用 musclead では数秒の本番断は許容、 商用なら blue/green 等で工夫が必要

---

## 3. AWS CLI 実行時は profile と region を明示

### 何が起きたか

```bash
aws ec2 revoke-security-group-ingress ...
# → "You must specify a region"
# → "Unable to locate credentials"
```

### 原因

Terraform は `provider "aws" { profile = "musclead-admin"; region = "ap-northeast-1" }`
で明示しているので動くが、 AWS CLI には伝わらない (別プロセス、 別設定)。

### 解決策

AWS CLI のコマンド全てに明示する:

```bash
aws --profile musclead-admin --region ap-northeast-1 ec2 ...
```

または環境変数で:
```bash
export AWS_PROFILE=musclead-admin
export AWS_REGION=ap-northeast-1
```

### 教訓

- Terraform と AWS CLI は **別々に認証情報を読む**
- 緊急時の手動オペレーションでは `--profile` `--region` を必ず明示

---

## 4. ElastiCache: Cluster vs Replication Group

### 選択肢

| Resource | 特徴 | 推奨度 |
|---|---|---|
| `aws_elasticache_cluster` | 古い API、 単一ノード前提 | △ |
| `aws_elasticache_replication_group` | 新規プロジェクト推奨、 単一ノードも OK | **✅** |

### 採用 (Replication Group)

```hcl
resource "aws_elasticache_replication_group" "main" {
  replication_group_id = "musclead-cache"
  engine               = "redis"
  engine_version       = "7.1"
  node_type            = "cache.t4g.micro"
  num_cache_clusters   = 1           # 単一ノード
  port                 = 6379

  subnet_group_name  = aws_elasticache_subnet_group.main.name
  security_group_ids = [var.cache_sg_id]
}
```

### output の差

| | Cluster | Replication Group |
|---|---|---|
| Endpoint | `cache_nodes[0].address` | `primary_endpoint_address` |

### 教訓

- AWS 公式は新規プロジェクトでも replication_group を推奨
- 単一ノードでも replication_group、 将来レプリカ追加 (`num_cache_clusters = 2`) も容易

---

## 5. `enable_cache` フラグ + `count` で 1 コマンド on/off

### パターン

```hcl
variable "enable_cache" {
  type    = bool
  default = false   # 普段は無効
}

module "cache" {
  source = "./modules/cache"
  count  = var.enable_cache ? 1 : 0
  ...
}

module "ecs" {
  cache_endpoint = var.enable_cache ? module.cache[0].endpoint : ""
}
```

### 使い方

```bash
terraform apply -var enable_cache=true    # 起動 ($13/月)
terraform apply -var enable_cache=false   # 停止 ($0/月)
```

### count 0 の時の参照

`var.enable_cache = false` → `module.cache` は **`module.cache[0]` 形式**になり、
参照すると「インデックス外」 エラー。 三項演算子で空文字に切り替えて回避:

```hcl
cache_endpoint = var.enable_cache ? module.cache[0].endpoint : ""
```

### BE 側

```go
host := os.Getenv("REDIS_HOST")
if host == "" {
    return NoOpCache{}   // キャッシュ無効モード
}
return NewRedisCache(host)
```

これで Terraform 側で停止しても BE は無停止で動く (ADR 0010 の Fail-open 設計)。

### 教訓

- 有償リソースは **default false** で安全側に倒す
- `count` + 三項演算子で「1 行 on/off」 が IaC のベストプラクティス
- BE 側も「リソース無し」 を許容する設計 (env 空 = NoOp) にすると IaC 切替が無停止になる

---

## 関連

- [ADR 0010: ElastiCache (Redis) 採用と Fail-open キャッシュ戦略](../adr/0010-elasticache-and-fail-open-cache.md)
