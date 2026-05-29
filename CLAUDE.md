# musclead — Project Context

> 筋トレ・食事・体重 一元管理 SaaS。 SODA入社準備のため Go + Connect-RPC + sqlc + DDD で実装。

---

## 📚 設計ドキュメント

- [ドメインモデル](docs/domain-model.md)
- [ADR](docs/adr/)

---

## 📐 Cursor Rules インデックス

`.cursor/rules/` に詳細ルール。

| ファイル | 内容 |
|---|---|
| [00-tech-stack.mdc](.cursor/rules/00-tech-stack.mdc) | 採用技術スタック一覧 |
| [01-architecture.mdc](.cursor/rules/01-architecture.mdc) | DDD + Modular Monolith |
| [02-value-object.mdc](.cursor/rules/02-value-object.mdc) | 値オブジェクト実装パターン |
| [03-coding-style.mdc](.cursor/rules/03-coding-style.mdc) | コーディング規約 |

---

## 👥 役割分担(AI と人間)

> **大原則: コードは基本的に人間が書く(学習目的)。 AI は勝手にコードを書かない。**
> AI がコードを書いてよいのは、 人間が明示的に「書いて」と依頼した時だけ。
> それ以外は「方針提示・ヒント・レビュー」に徹する。

| 種類 | 担当 |
|---|---|
| 設計・方針 | 議論で決定 |
| **すべての Go コード実装**(VO / Entity / UseCase / Repository / Handler / ヘルパー等) | **人間**(学習目的) |
| 設定ファイル / CI / Makefile / migration 等の非ロジック | AI 可(依頼ベース) |
| コードレビュー | **AI** |
| 詰まった時のガイド・ヒント | **AI**(原則コードは出さず方針提示。 明示依頼時のみ正解コード) |

---

## 🔄 進行ルール

- **7割で push**、 完璧主義禁止
- **詰まったら30分ルール** で質問
- **1機能 = 1 commit**(Conventional Commits)
- **すべて最新安定版**を採用
- メジャー更新は ADR 記録

## 💬 回答スタイル(重要)

- **回答は簡潔に**(目安: 通常 20 行以内、 深掘りでも 50 行以内)
- **結論先出し**、 冗長な前置き禁止
- 表やリストは **要点に絞る**(網羅性より要点)
- 余談・歴史的経緯・例文は **明示的に求められた時だけ**
