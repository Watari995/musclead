# ADR 0011: 体重時系列 API の設計 (全フィールド入り + 期間プリセット + 即時 invalidation)

## ステータス
採用 (2026-06-08)

## コンテキスト

体重グラフ機能のため時系列データを返す API が必要。 既存の `GET /weights` (一覧、 offset/limit) はページネーション用途で、 グラフには適さない。

要件:
- グラフは画像参考 UI (筋肉量タブ + 期間プリセット切替 + 過去スクロール)
- 体重 / 体脂肪率 / 骨格筋量 を **タブ切替で表示、 切替時に待ちたくない**
- record / update / delete した瞬間にグラフ反映 (stale 禁止)
- 期間: 1週間 / 1ヶ月 / 3ヶ月 / 半年 / 1年 のプリセット
- 過去スクロールで継続取得
- 1日複数回保存対応 (raw record をそのまま返す、 集計しない)

## 判断

### ① エンドポイント: 専用 `GET /weights/timeseries`

既存の `GET /weights` (offset/limit ページネーション) とは別エンドポイント。 用途が違うので分離。

```
GET /weights/timeseries?period=1year[&before=2025-06-08T00:00:00Z]
```

| パラメータ | 必須 | 値 |
|---|---|---|
| `period` | ✅ | `1week` / `1month` / `3months` / `halfyear` / `1year` |
| `before` | - | ISO 8601、 未指定なら now。 「これ以前 〜 period 分」 を返す。 過去スクロール用 |

from/to の自由日付指定は**しない** (UI に該当機能なし、 当面は period プリセットで十分)。

### ② レスポンス: 全フィールド入り (C 案)

```json
{
  "period": "1year",
  "weights": [
    {
      "id": "...",
      "weight_kg": "70.5",
      "body_fat_percentage": "15.2",
      "skeletal_muscle_kg": "30.1",
      "measured_at": "2026-06-01T08:00:00Z"
    }
  ]
}
```

- record の **全フィールドを 1 回で返す**
- FE はタブ切替で「どのフィールドをプロットするか」 を切替えるだけ (再 fetch なし、 瞬時)

### ③ キャッシュ戦略: Redis ZSET + HASH で record 単位の部分更新

データ構造:

```
ZSET  weights:timeseries:{user_id}:idx
    score=measured_at(unix), member=weight_id

HASH  weights:timeseries:{user_id}:data
    field=weight_id, value=JSON (record 全フィールド)
```

操作:

| Use case | Redis 操作 |
|---|---|
| 期間取得 | `ZRANGEBYSCORE idx <from> <to>` → id 一覧 → `HMGET data <ids...>` (2 RTT) |
| record 追加 | `MULTI` → `ZADD idx` + `HSET data` → `EXEC` (atomic) |
| record 更新 (measured_at 不変) | `HSET data <id> <newJson>` (ZSET 触らない) |
| record 更新 (measured_at 変更) | `MULTI` → `ZADD idx <newScore>` + `HSET data` → `EXEC` |
| record 削除 | `MULTI` → `ZREM idx <id>` + `HDEL data <id>` → `EXEC` |

TTL は **長め (24 時間)** に設定可。 部分更新で常に最新を保つので、 全体 TTL の依存度が低い。

### ④ Cache interface 拡張: weight 専用 `WeightTimeseriesCache`

shared の `Cache` interface (Get/Set/Delete) は **触らない**。 weight 固有のデータ構造操作を weight module 内に閉じる:

```go
// weight/internal/domain/weight_timeseries_cache.go
type WeightTimeseriesCache interface {
    FindByPeriod(ctx context.Context, userID UserID, from, to time.Time) (weights []*Weight, hit bool, err error)
    Add(ctx context.Context, weight *Weight) error
    Update(ctx context.Context, weight *Weight) error
    Delete(ctx context.Context, userID UserID, weightID WeightID, measuredAt time.Time) error
}
```

infra 実装は go-redis を直接使い、 JSON encode/decode と value object 復元を内部で行う (= Repository と同じ責務分担)。 これにより usecase はキャッシュからそのまま entity を受け取れて、 caller 側に decode boilerplate を持ち込まない。

NoOp 実装も同 interface を満たすよう用意 (常に miss / 無視)。

## なぜ部分更新 (ZSET + HASH) を選んだか

代替案:

| 案 | データ構造 | invalidate 戦略 |
|---|---|---|
| **A: 全削除 + 次回再構築** | string (JSON 全体) | 記録/更新/削除のたびに SCAN + DEL で `weights:timeseries:{user_id}:*` を全削除 |
| **B: 全体を読み書き直し** | string (JSON 全体) | GET → decode → 差分追加 → encode → SET |
| **C: ZSET + HASH の record 単位部分更新** ✅ | Sorted Set + Hash | 該当 record のみ ZADD/HSET/HDEL で incremental に同期 |

### A (全削除) を却下した理由

- record 1 件追加 / 更新 / 削除で **全期間のキャッシュ (1week〜1year × 全 before 過去) が消える**
- 次の GET 時に DB 全再読 → 体感に出ない程度ではあるが、 cache の意味が弱い
- record 編集が頻発する musclead では「全削除→再構築」 を毎回繰り返すのは設計的にもったいない

### B (全体読み書き) を却下した理由

- read-modify-write の競合リスク (同時に 2 つの記録があった場合 lost update)
- JSON 全体を毎回 encode/decode = CPU と帯域の無駄
- record 数が増えるほど不利

### C を採用した理由

- record 単位で **過去のデータをピンポイントで更新・削除できる** ← user 要件
- ZSET によりインデックス化されていて、 期間 (from/to) 絞り込みが Redis 側で完結 (`ZRANGEBYSCORE`)
- `MULTI/EXEC` で ZSET と HASH の操作を atomic に揃えられる
- record 追加: O(log N)、 期間取得: O(log N + M)、 操作はほぼ全てミリ秒
- stale が原理的に発生しない (常に最新)
- **学習価値**: Redis の代表的なデータ構造 (Sorted Set + Hash) と atomic 操作 (MULTI/EXEC) を実装で習得できる

### shared `Cache` interface に持ち込まなかった理由

ZSET / HASH / `MULTI/EXEC` は Redis 固有のデータ構造操作。 汎用 `Cache` interface (`Get/Set/Delete`) に持ち込むと:

- NoOpCache 等の代替実装が複雑化
- 他 module (meal / training 等) が将来 Cache を使う際に **不要な ops が露出**
- interface の抽象レベルが壊れる

そのため shared の `Cache` はそのまま (汎用 string KV) で残し、 weight 固有の cache は weight module 内に専用 interface (`WeightTimeseriesCache`) として閉じる。 infra 実装は shared の Cache 経由ではなく go-redis を直接使う (Repository が gorp を直接使うのと同じ責務分担)。

## なぜ全フィールド入り (C 案) にしたか

代替案:
- A: 別 API (`/weights/timeseries`、 `/body-fat/timeseries` 等)
- B: 統合 API + type で BE 切替 (`?type=weight_kg`)

却下理由:
- A: タブ切替で別 endpoint 呼び出し → ローディング待ち、 user 要件「待ちたくない」 に反する。 重複コードも増える
- B: タブ切替で同じ endpoint だが type 違い → 再 fetch、 やはりローディング待ち発生

採用 (C):
- 1 回の取得で全フィールド入り → タブ切替は FE 内のメモ化された表示切替だけ → 瞬時
- 当面 5 フィールド (体重 / 体脂肪率 / 骨格筋量 / BMI / 内臓脂肪 等) 想定で 1 年分 ~700KB、 mobile でも実用範囲
- キャッシュ key 数も少なく済む (type 軸が無い)

## なぜ即時 invalidation を採用したか (ADR 0010 から方針変更)

ADR 0010 では「TTL 任せの stale 許容」 と決めていたが、 user 要件確認の結果:
> 「stale したらダメ、 record した瞬間にグラフ反映してほしい」

を明示された。 musclead のユーザー体験 (記録 → 即グラフ確認) では stale 許容できない。

|  | ADR 0010 当初 | 本 ADR で変更 |
|---|---|---|
| invalidation | TTL 任せ | record 単位の部分更新 |
| TTL | 10 分 | 24 時間 (部分更新で常に最新を保つので長めで OK) |
| 整合性 | 最大 10 分遅延 | record 直後に最新 |

実装コスト:
- weight 専用 `WeightTimeseriesCache` interface 新設
- weight/infra に Redis ZSET + HASH 実装
- record / update / delete usecase に cache 同期呼び出し追加 (best effort、 失敗してもアプリ続行)

## 将来の肥大化対応 (B 案へのスムーズ移行パス)

5 フィールドを超えて肥大化 (10 フィールド = 1MB 超) したら、 **optional `type` パラメータを追加**して B 案に進化可能:

```
GET /weights/timeseries?period=1year                # 既存、 全部入り (後方互換)
GET /weights/timeseries?period=1year&type=weight_kg # 新規、 1 フィールドのみ
```

旧 client は引き続き動き、 新 client は type 指定で軽量取得。 API パスは同じなので **段階的移行可能**。

## やらないこと

- 自由 from/to 日付指定 (UI 不要)
- 期間プリセット以外の値 (e.g. `2weeks`) は受け付けない
- 集計値 (avg / min / max) の BE 計算: raw record を返して FE が必要なら計算
- type / fields 絞り込み (将来追加可能)
- 別 type (体脂肪率) 専用 endpoint
- 体重以外の entity (meal, training) でも同じ pattern を適用するかは別 ADR で決める

## 影響

### BE
- `dto`: TimeseriesWeightsResponse 追加 (既存 `WeightDTO` を再利用)
- `weight/internal/domain/weight_timeseries_cache.go`: `WeightTimeseriesCache` interface 新規
- `weight/internal/domain/weight_repository.go`: `FindAllByUserIDAndPeriod` (新規 method、 cache miss 時の DB fallback で使う)
- `weight/internal/infra/weight_repository.go`: 期間絞り SQL 実装
- `weight/internal/infra/weight_timeseries_cache.go`: Redis (ZSET + HASH + MULTI/EXEC) 実装 + NoOp 実装
- `weight/internal/usecase/get_weight_timeseries.go`: Cache → 失敗時は Repo → cache populate
- `weight/internal/usecase/{record,update,delete}_weight.go`: cache 同期呼び出し追加 (best effort)
- `weight/internal/handler/weight_handler.go`: `GET /weights/timeseries` 追加 + swag
- `weight/weight.go` (facade): NewModule が WeightTimeseriesCache 依存を受け取る
- `cmd/server/main.go`: shared cache client を元に weight 専用 cache を組み立てて weight.NewModule に渡す
- 既存 `shared/domain/cache.go` (Get/Set/Delete) は **触らない**

### FE
- `features/weight/api/weights.ts`: useTimeseriesQuery (TanStack Query)
- `features/weight/ui/WeightGraph.tsx` (新規): recharts で line chart、 タブ切替で type 変更
- `app/weights/page.tsx`: グラフセクション追加

### Infra
- ElastiCache 既に起動済 (ADR 0010 で apply 済み)

## 関連 ADR
- [ADR 0010](0010-elasticache-and-fail-open-cache.md): ElastiCache 採用 (本 ADR で「TTL 任せ」 → 「明示 invalidation」 に方針更新)
