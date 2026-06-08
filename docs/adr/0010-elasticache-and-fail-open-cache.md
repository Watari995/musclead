# ADR 0010: ElastiCache (Redis) 採用と Fail-open キャッシュ戦略

## ステータス
採用 (2026-06-08)

## コンテキスト

体重機能 Phase 2 で 日/週/月集計 API のキャッシュ層が必要になった。 さらに将来の集計系エンドポイントでも再利用したい。

要件:
- 学習目的: SODA 入社準備のため、 業界標準なキャッシュ運用を経験する
- コスト最適化: 学習期間中も自由に on/off できる
- 信頼性: キャッシュ層が落ちてもアプリ全体は動き続ける
- 再現性: コードを残し、 いつでも復活できる

## 判断

### ① キャッシュ層: **AWS ElastiCache (Redis 7.x)**

- Node type: `cache.t4g.micro` 単一ノード、 Multi-AZ なし
- Engine: Redis 7.x
- Subnet: RDS と同じ public_subnet
- SG: ECS Fargate からの 6379 のみ許可

セルフホスト (Fargate / EC2) も検討したが却下 (後述)。

### ② 切り替え機構: **`enable_cache` 変数 + count 制御**

```hcl
variable "enable_cache" { type = bool, default = false }

module "cache" {
  count  = var.enable_cache ? 1 : 0
  source = "./modules/cache"
  ...
}

module "ecs" {
  cache_endpoint = var.enable_cache ? module.cache[0].endpoint : ""
}
```

`terraform apply -var enable_cache=true/false` の 1 コマンドで起動 / 停止。

### ③ キャッシュ戦略: **Cache-aside + Fail-open**

```
Request → cache.Get
            ├ ヒット → 返す
            ├ ミス  → DB 取得 → cache.Set (best effort) → 返す
            └ ERROR → DB 取得 → 返す
```

cache の Get/Set エラーは **全て無視して DB にフォールバック**。 cache 障害でアプリ全体は落ちない。

### ④ Cache interface 抽象化 + 環境変数で実装切替

```go
type Cache interface {
    Get(ctx, key) ([]byte, error)
    Set(ctx, key, val, ttl) error
}

// 本番 (REDIS_HOST 設定済)
type RedisCache struct { client *redis.Client }

// destroy 中 / ローカル無し (REDIS_HOST="")
type NoOpCache struct{}  // 常に miss を返す
```

```go
func newCache() Cache {
    if os.Getenv("REDIS_HOST") == "" { return NoOpCache{} }
    return NewRedisCache(...)
}
```

### ⑤ TTL: **集計系は短めに (5〜15 分)**

- 体重日次集計: 10 分
- 整合性より単純さ優先
- 明示的な invalidation はしない (TTL 任せ)

## なぜ ElastiCache を採用したか

| 観点 | ElastiCache | Fargate Redis セルフ | ローカル Docker のみ |
|---|---|---|---|
| 月額 | $13 | $5-10 (+EFS) | $0 |
| 運用負荷 | ゼロ | 永続化・パッチ・監視を自分で | - |
| SODA 想定 | ✅ 業界標準 | 少数派 | 本番なし |
| 学習価値 | 本番マネージドサービスの運用感 | コンテナ運用に時間を取られ散漫 | 本番経験ゼロ |

コスト差は誤差レベル ($1〜2 / 月)、 学習価値の差が大きいため ElastiCache を採用。

## なぜ Fail-open を採用したか

代替案: cache 障害で 5xx を返す。

却下理由:
- musclead の cache は **性能のための補助層**、 真実は DB にある
- cache 障害でユーザー体験を壊すのは過剰反応
- 業界標準 (Twitter / Slack / GitHub) は cache 障害時の grace degradation
- Fail-open ならキャッシュ destroy 中も無停止で動く = 学習用途と相性が良い

## なぜ enable_cache フラグ方式を採用したか

代替案: `terraform destroy -target=module.cache` で消す。

却下理由:
- ECS module が `module.cache.endpoint` を参照しているため、 cache だけ destroy すると **参照エラー**で plan が壊れる
- enable_cache フラグなら ECS の env も `""` で同期更新される
- 1 コマンド (`apply -var enable_cache=...`) で切替可、 オペレーションがシンプル

## 影響

### インフラ (Terraform)

- 新規: `modules/cache/{main,variables,outputs}.tf`
- 変更: `modules/network/{main,outputs}.tf` に `aws_security_group.cache` 追加
- 変更: `modules/ecs/{main,variables}.tf` に `cache_endpoint` + `REDIS_HOST` / `REDIS_PORT` 追加
- 変更: `main.tf` に `enable_cache` 変数 + `module "cache"` 配線

### BE

- 新規: `internal/shared/infra/cache/redis_cache.go` (RedisCache 実装)
- 新規: `internal/shared/infra/cache/noop_cache.go` (NoOpCache)
- 新規: `internal/shared/domain/cache.go` (Cache interface)
- 修正: `cmd/server/main.go` で env 見て切替
- 修正: 集計 usecase は Cache 依存を受け取る (cache-aside パターン適用)

### 運用

- 通常: `enable_cache=false` (キャッシュなし、 $0)
- 学習 / 動作確認時: `enable_cache=true` (キャッシュあり、 月 $13 ベース)
- 入社後: コードは残しつつ `enable_cache=false` で停止

### コスト見積

| シナリオ | 月額 |
|---|---|
| 通常 (停止中) | $0 |
| 学習継続中 (常時起動) | $13 |
| 学習期間 1 ヶ月のみ起動 | $13 |

## やらないこと

- **複数ノード / Multi-AZ / Replication Group**: musclead 規模では不要
- **TLS / 認証 (AUTH トークン)**: VPC 内 SG で絞れば十分
- **キャッシュの明示的 invalidation**: TTL 任せで単純さ優先
- **memcached の選択肢**: Redis の方が SODA 想定との一致度が高い
- **CloudWatch Alarm**: 学習用途では過剰、 必要になれば別 ADR

## 関連 ADR

- [ADR 0007](0007-infra-mvp-and-monorepo.md): インフラ / monorepo 構成
