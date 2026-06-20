# ADR-0023: 食品バーコード検索・登録機能の設計

- ステータス: Accepted
- 日付: 2026-06-20
- 関連: [ADR-0002 DDD + Modular Monolith](0002-ddd-modular-monolith.md)、[ADR-0022 meal_template](0022-meal-templates.md)

## コンテキスト

バーコードスキャンおよび食品名検索から、カロリー・PFC を meal 記録フォームに自動入力したい。
MyFitnessPal を参考にしたハイブリッド設計：meal テーブルは引き続きカロリー・PFC を直接保持し、食品検索はあくまで入力補助として機能する。

将来的に「食品一覧」「栄養素確認」など meal 記録以外のユースケースも想定しているため、`meal` モジュールとは独立した `food` モジュールを新設する。

## 決定

### モジュール配置

`server/internal/food/` を新たな bounded context として作成する。
`meal` モジュールからの参照はなく、クライアント（mobile/web）が `food` API を叩いてフォームに値を詰めるだけで完結する。

```
server/internal/food/
  food.go                           -- Module 公開インターフェース
  dto/
    food_product_dto.go
  internal/
    domain/
      food_product.go               -- FoodProduct エンティティ
    usecase/
      search_by_barcode.go          -- バーコード検索（自社DB → Open Food Facts → 404）
      search_by_name.go             -- 名前検索（自社DB のみ）
      create_food_product.go        -- ユーザー登録
    infra/
      food_product_model.go         -- gorp モデル
      food_product_repository.go    -- DB 実装
      open_food_facts_client.go     -- 外部 API クライアント
    handler/
      food_handler.go
```

### テーブル設計

```sql
CREATE TABLE food_products (
  id              BINARY(16)    NOT NULL,
  barcode         VARCHAR(14)   NULL,           -- JAN/EAN/UPC。名前検索登録時は NULL
  name            VARCHAR(100)  NOT NULL,
  calories        INT           NOT NULL,       -- 1食分
  protein_g       DECIMAL(6,2)  NULL,
  fat_g           DECIMAL(6,2)  NULL,
  carbohydrate_g  DECIMAL(6,2)  NULL,
  register_source VARCHAR(20)   NOT NULL,       -- 'open_food_facts' | 'user'
  created_at      DATETIME(6)   NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  updated_at      DATETIME(6)   NOT NULL DEFAULT CURRENT_TIMESTAMP(6)
                                ON UPDATE CURRENT_TIMESTAMP(6),
  PRIMARY KEY (id),
  KEY idx_barcode (barcode),
  KEY idx_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
```

**設計上の決定**:
- `barcode` に UNIQUE 制約を設けない：同一バーコードの複数登録を許容し、ヒット時は候補リストを返す（MyFitnessPal と同様）
- `created_by_user_id` は持たない：全ユーザー共有の食品マスタとして扱い、編集・削除権限の管理を省略する
- カロリー・PFC は「1食分」の値を保存する（100g換算ではない）
- `register_source` で出典を区別し、Open Food Facts 由来データはキャッシュとして保存する

### 外部 API

**Open Food Facts** を使用する（無料・APIキー不要）。

```
GET https://world.openfoodfacts.org/api/v2/product/{barcode}.json
```

バーコード検索フロー:
1. 自社DB（`food_products`）照会 → ヒット → 候補リスト返却
2. 自社DBになければ Open Food Facts API 呼び出し
3. Open Food Facts でヒット → `register_source='open_food_facts'` で DB にキャッシュ → 返却
4. どちらもなし → 404（クライアント側でユーザー登録フローへ遷移）

### API 設計

全エンドポイントに認証（Bearer トーク）が必要。

| メソッド | パス | 説明 |
|---|---|---|
| `GET` | `/food_products?q={name}` | 名前検索（前方一致、自社DB） |
| `GET` | `/food_products/barcode/{code}` | バーコード検索 |
| `POST` | `/food_products` | ユーザーによる新規登録 |

### 値オブジェクト

| フィールド | VO |
|---|---|
| `name` | `valueobject.String100`（既存） |
| `barcode` | `valueobject.String14`（新規） |
| `calories` | `valueobject.NonNegativeInt`（既存） |
| `protein_g` / `fat_g` / `carbohydrate_g` | `*valueobject.NonNegativeDecimal`（既存・nullable） |
| `register_source` | `valueobject.RegisterSource`（新規） |

### meal との関係

`meals` テーブル・API に変更はない。
クライアントが `/food_products` で取得した値を `POST /meals` の body に詰めるだけで完結する。

## 代替案

- **`meal` モジュール内に food_products を置く**: meal 以外のユースケース（食品一覧等）が生じた際に `meal` への不自然な依存が発生するため不採用。
- **`meal_food_items` テーブルで食品単位の記録を持つ（フル MFP スタイル）**: `meals` テーブルの設計変更が大きく、現時点の要件に対してオーバースペックなため不採用。ハイブリッド設計で十分。
- **外部APIをキャッシュしない**: シンプルだが毎回レイテンシが発生し API 障害時に機能しなくなるため不採用。

## 影響

- migration: `000024_create_food_products.up.sql` を追加。既存テーブルへの変更はなし。
- `server/internal/food/` を新規作成し、`main.go` の `newMux()` に `foodModule` を追加。
- Mobile: `mobile_scanner` パッケージでバーコードスキャン UI を追加。
- Web: `@zxing/browser` でカメラ経由バーコードスキャン UI を追加。
