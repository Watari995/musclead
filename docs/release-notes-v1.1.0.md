# v1.1.0 リリースノート — meal_templates 機能

**リリース日**: 2026-06-18

---

## 概要

よく食べる食事をテンプレートとして保存し、ワンタップで記録に転用できる機能を追加。  
Web・iOS(Flutter) 両方に対応。

---

## 新機能

### テンプレート保存 (CRUD)
- テンプレート名・食事タイプ・カロリー・P/F/C を保存できる
- 作成時に `display_order` が自動採番される（MAX+1）

### テンプレートから記録をプリフィル
- **Web**: 食事ページ右側の「テンプレート」カードをクリック → 記録フォームに自動入力
- **iOS**: 食事画面のブックマークアイコン → テンプレート一覧シート → タップで記録シートへプリフィル

### テンプレート並び替え
- `POST /meal_templates/reorder` で表示順を一括更新

---

## API

| メソッド | パス | 説明 |
|---|---|---|
| GET | `/meal_templates` | 一覧取得（offset pagination） |
| POST | `/meal_templates` | 作成 |
| PUT | `/meal_templates/{id}` | 更新 |
| DELETE | `/meal_templates/{id}` | 削除 |
| POST | `/meal_templates/reorder` | 並び替え |

---

## 技術ノート

- `meal_template_id` は `meals` テーブルに**持たない** (copy-on-use)。テンプレートはプリフィル専用。
- DB: `meal_templates` テーブル追加 (migration 000023)。`user_id` に FK。
- インデックス: `(user_id, display_order ASC, created_at ASC)` — filesort なしでユーザー別一覧取得。
- ADR: [0022-meal-templates](adr/0022-meal-templates.md)
