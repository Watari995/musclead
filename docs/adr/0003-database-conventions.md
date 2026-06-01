# ADR 0003: DB 共通設計規約

## ステータス
採用 (2026-05-23)

## コンテキスト
MySQL 8.0 / Aurora Serverless v2 を使用する自作SaaS で、 全テーブル共通の設計ルールを統一する必要がある。

---

## 決定事項

### ① ID 型: `BINARY(16)`

- **採用**: `BINARY(16)`
- **却下**: `CHAR(36)`
- **理由**:
  - 16バイトでインデックス効率最大
  - `CHAR(36)` だと容量 2.25倍、 B+treeのページ消費が増える
- **デメリット**: デバッグ時に hex 変換が必要(`HEX(id)` で確認)

### ② UUID バージョン: **v7**

- **採用**: UUID v7 (RFC 9562)
- **却下**: UUID v4 / CUID2
- **理由**:
  - 先頭48bit が timestamp = **B+treeフレンドリー**(末尾追記のみ)
  - ページスプリット発生せず、 INSERT高速・断片化なし
  - キャッシュ効率最大(ホットページが末尾のみ)
  - RFC 9562 で標準化済み
- **CUID2 を却下した理由**: 24文字必要、 標準化されてない、 タイムスタンプ漏洩リスクのメリットは個人開発では不要

### ③ タイムスタンプ列: `DATETIME(6)` + UTC運用

- **採用**: `DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6)`
- **却下**: `TIMESTAMP`
- **理由**:
  - 2038年問題なし(`TIMESTAMP` は 1970-2038年で破綻)
  - タイムゾーン変換が起きない(挙動が明確)
  - マイクロ秒精度でログ・トレースに有利
- **運用ルール**:
  - DB: UTC で保存(タイムゾーン情報なし)
  - アプリ層: `time.UTC` で統一
  - 表示時にユーザータイムゾーンに変換(FE 担当)
- **標準カラム**:
  ```sql
  created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6)
  updated_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6)
  ```

### ④ Soft Delete: users のみ採用

- **採用**: `users` テーブルに `deleted_at DATETIME(6) NULL`
- **却下**: 他テーブル(meals / trainings / weights / meal_photos 等)は物理削除
- **理由**:
  - `users`: アカウント誤削除の復旧ニーズ、 監査ログ的価値あり
  - `meals` 等: 個人記録、 履歴不要、 物理削除でシンプルに
- **クエリ規約**:
  ```sql
  SELECT * FROM users WHERE deleted_at IS NULL;
  ```
- **UNIQUE 制約の注意**:
  ```sql
  -- email を UNIQUE にする場合、 削除済みは除外したい
  -- → 部分ユニーク or 削除時に email を物理的に変更
  ```

### ⑤ Charset / Collation

- **採用**: `utf8mb4` / `utf8mb4_0900_ai_ci`(MySQL 8.0 デフォルト)
- **理由**:
  - `utf8mb4` = 4バイト UTF-8 完全実装、 **絵文字対応**(`utf8` は罠なので使わない)
  - `utf8mb4_0900_ai_ci` = Unicode 9.0 ベース、 アクセント・大文字小文字無視
  - メモ欄に絵文字、 email検索で大文字小文字気にしない、 など実用的
- **DATABASE 作成時に明示**:
  ```sql
  CREATE DATABASE musclead
    CHARACTER SET utf8mb4
    COLLATE utf8mb4_0900_ai_ci;
  ```

---

## 結果

- ✅ 全テーブル共通の物理設計が統一される
- ✅ パフォーマンス特性が予測可能
- ✅ 個人開発・SODA スケール双方で耐えうる設計
- ❌ BINARY(16) のデバッグはやや面倒(`HEX()` 必須)
- ❌ Soft Delete の UNIQUE制約は工夫が必要

---

## 関連

- 次のテーブル設計はこの規約に従う(`docs/db-schema.md` 参照、 まだ作成前)
- マイグレーションツール: goose(別 ADR で記録予定)
