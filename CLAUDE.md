# musclead — Project Context

> 筋トレ・食事・体重 一元管理 SaaS。 SODA入社準備のため、 バックエンドは Go(net/http REST + OpenAPI/swag)+ gorp + DDD、 iOS は Flutter で実装。

---

## 📚 設計ドキュメント

- [ドメインモデル](docs/domain-model.md)
- [ADR](docs/adr/)（iOS=Flutter 採用: [ADR-0021](docs/adr/0021-adopt-flutter-for-ios.md)）
- iOS アプリ: `mobile/`（Flutter。 実装は AI 担当。 **v1.0.0 App Store 審査中** as of 2026-06-18）

---

## 📐 Cursor Rules インデックス

`.cursor/rules/` に詳細ルール。

| ファイル | 内容 |
|---|---|
| [00-tech-stack.mdc](.cursor/rules/00-tech-stack.mdc) | 採用技術スタック一覧 |
| [01-architecture.mdc](.cursor/rules/01-architecture.mdc) | DDD + Modular Monolith |
| [02-value-object.mdc](.cursor/rules/02-value-object.mdc) | 値オブジェクト実装パターン |
| [03-coding-style.mdc](.cursor/rules/03-coding-style.mdc) | コーディング規約 |
| [05-shared-first.mdc](.cursor/rules/05-shared-first.mdc) | shared/ 優先、 局所ヘルパー重複定義の禁止 |
| [06-industry-standard.mdc](.cursor/rules/06-industry-standard.mdc) | 業界標準・SODA 流儀を踏襲、 略語/独自手法を避ける |
| [07-review-rule.mdc](.cursor/rules/07-review-rule.mdc) | **AI レビュー作法**: PR push、 慣習整合、 Go 慣習優先、 代替案提示、 簡潔指摘 |
| [10-web-design-system.mdc](.cursor/rules/10-web-design-system.mdc) | Web UI デザインシステム(snkrdunk テイスト)、 web/ 編集時必読 |
| [11-mobile-responsive.mdc](.cursor/rules/11-mobile-responsive.mdc) | モバイル対応(iPhone SE 基準、 ハンバーガー必須)、 web/ 編集時必読 |

---

## 👥 役割分担(AI と人間)

> **大原則: コードは基本的に人間が書く(学習目的)。 AI は勝手にコードを書かない。**
> AI がコードを書いてよいのは、 人間が明示的に「書いて」と依頼した時だけ。
> それ以外は「方針提示・ヒント・レビュー」に徹する。

| 種類 | 担当 |
|---|---|
| 設計・方針 | 議論で決定 |
| **すべての Go コード実装**(VO / Entity / UseCase / Repository / Handler / ヘルパー等) | **人間**(学習目的) |
| **iOS(Flutter / `mobile/`)の実装** | **AI**(委任。 Go バックエンドが学習対象のため) |
| 設定ファイル / CI / Makefile / migration 等の非ロジック | AI 可(依頼ベース) |
| コードレビュー | **AI** |
| 詰まった時のガイド・ヒント | **AI**(原則コードは出さず方針提示。 明示依頼時のみ正解コード) |

---

## 🔄 進行ルール

- **詰まったら30分ルール** で質問
- **1機能 = 1 commit**(Conventional Commits)
- **すべて最新安定版**を採用
- メジャー更新は ADR 記録
- **shipping を急いだ局所最適コードは禁止**(詳細は [05-shared-first.mdc](.cursor/rules/05-shared-first.mdc))
  - 新しいヘルパー関数は `shared/` 配下を grep してから書く
  - 既存パッケージに収まる時は必ずそこに追加、 既存命名規約を踏襲

## 💬 回答スタイル(重要)

- **回答は簡潔に**(目安: 通常 20 行以内、 深掘りでも 50 行以内)
- **結論先出し**、 冗長な前置き禁止
- 表やリストは **要点に絞る**(網羅性より要点)
- 余談・歴史的経緯・例文は **明示的に求められた時だけ**
