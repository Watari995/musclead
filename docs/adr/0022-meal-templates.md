# ADR-0022: meal_template 機能の設計

- ステータス: Accepted
- 日付: 2026-06-18
- 関連: [ADR-0002 DDD + Modular Monolith](0002-ddd-modular-monolith.md)

## コンテキスト

よく食べる食事をテンプレートとして保存し、ワンタップで記録できる機能を追加する。ルーティン食の入力コストを削減するのが目的。

## 決定

### テーブル設計

```sql
CREATE TABLE meal_templates (
  id             BINARY(16)    NOT NULL,
  user_id        BINARY(16)    NOT NULL,
  name           VARCHAR(100)  NOT NULL,
  display_order  INT           NOT NULL DEFAULT 0,
  meal_type      VARCHAR(20)   NOT NULL,
  calories       INT           NOT NULL,
  protein_g      DECIMAL(6, 2) NULL,
  fat_g          DECIMAL(6, 2) NULL,
  carbohydrate_g DECIMAL(6, 2) NULL,
  created_at     DATETIME(6)   NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  updated_at     DATETIME(6)   NOT NULL DEFAULT CURRENT_TIMESTAMP(6)
                                        ON UPDATE CURRENT_TIMESTAMP(6),
  PRIMARY KEY (id),
  KEY idx_meal_templates_user (user_id, display_order ASC, created_at ASC),
  CONSTRAINT fk_meal_templates_user_id
    FOREIGN KEY (user_id) REFERENCES users(id)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_0900_ai_ci;
```

### `meals` テーブルへの `meal_template_id` FK は追加しない

テンプレートは「コピー元」であり、記録との永続的な参照関係を持たない。テンプレートを更新・削除しても過去の記録に影響しないことを保証するため、コピー・オン・ユース（記録時にデータをコピー）とする。

### API 設計

| メソッド | パス | 説明 |
|---|---|---|
| `GET` | `/meal_templates` | 一覧（全フィールドを返す。詳細エンドポイントは不要） |
| `POST` | `/meal_templates` | 作成 |
| `PUT` | `/meal_templates/{id}` | 更新 |
| `DELETE` | `/meal_templates/{id}` | 削除 |
| `PUT` | `/meal_templates/{id}/order` | 表示順変更 |

一覧はネストした子リソースを持たないため、詳細エンドポイントは設けない。ソート順は `display_order ASC, created_at ASC`。

テンプレートから記録する専用エンドポイントは設けない。モバイルがテンプレートのデータを `POST /meals` の body に詰めるだけで完結する。

### 値オブジェクト

| フィールド | VO |
|---|---|
| `name` | `valueobject.String100`（既存） |
| `meal_type` | `valueobject.String20`（既存・meals と統一） |
| `calories` | `valueobject.NonNegativeInt`（既存） |
| `protein_g` / `fat_g` / `carbohydrate_g` | `*valueobject.NonNegativeDecimal`（既存・nullable） |

### モジュール配置

既存 meal モジュールのパターンに従い `server/internal/meal/` 配下に追加する。外部からは `meal.Module` 経由でのみアクセスできる。

```
server/internal/meal/
  dto/meal_template_dto.go
  internal/
    domain/meal_template.go
    infra/meal_template_models.go
    infra/meal_template_repository.go
    usecase/create_meal_template.go
    usecase/list_meal_templates.go
    usecase/update_meal_template.go
    usecase/delete_meal_template.go
    usecase/reorder_meal_template.go
    handler/meal_handler.go  # 既存に追記
```

## 代替案

- **`meal_template_id` を meals に持たせる**: テンプレート更新時に過去の記録が書き変わる、または参照整合性のために削除制限が発生する。今回の要件（「ルーティン食を秒で入力」）には不要なため不採用。
- **テンプレート専用 module を新設**: テンプレートは meal の付随機能であり独立した集約を持つほどの複雑さがないため、meal module 内に留める。

## 影響

- migration: `000X_create_meal_templates.up.sql` を追加。`meals` テーブルへの変更はなし。
- `meal.Module` の `NewModule()` に `MealTemplateRepository` と関連 usecase を追加。
- モバイル: テンプレート一覧画面 + 選択でフォームにプリフィル or ワンタップ記録の UI を追加。
