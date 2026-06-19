# ADR-0021: iOS ネイティブアプリに Flutter を採用する

- ステータス: Accepted
- 日付: 2026-06-17
- 関連: [ADR-0017 Stripe 連携](0017-stripe-integration-details.md), [ADR-0013 purchase/payment 分離](0013-purchase-payment-separation.md)

## コンテキスト

Web(Next.js)に続き、iOS ネイティブアプリを提供する。フロントエンド実装は本プロジェクトの学習主眼（Go バックエンド）ではないため AI に委任しつつ、SODA 実務で通用する品質を要件とする。

バックエンドは **REST/JSON API**（`server/docs/swagger.yaml` を `swag` で生成、Web は `openapi-fetch` で型生成）。CLAUDE.md にあった「Connect-RPC」は実態と異なるため本 ADR で訂正する（proto / connect-go は不在）。

## 決定

**Flutter(Dart)を採用**し、monorepo の `mobile/` に配置する。

| 領域 | 採用 | 理由 |
|---|---|---|
| 状態管理 | Riverpod v2（手書き Notifier）+ hooks_riverpod | モダン・テスト容易。`riverpod_generator`/`riverpod_lint` は **freezed 3 系と依存衝突**するため不使用（コード生成に頼らず手書き） |
| HTTP | dio + dio_cookie_manager + cookie_jar | interceptor で Bearer 付与 / 401→refresh / Cookie 永続化 |
| モデル | Freezed 3 + json_serializable（`field_rename: snake`） | バックエンドの snake_case JSON に自動整合。**小数は文字列で授受**されるため `Decimal` 変換コンバータを用意 |
| ルーティング | go_router（auth リダイレクトガード） | 宣言的・ディープリンク対応 |
| 認証 | バックエンド JWT(access 15分) + refresh は HttpOnly Cookie を流用 | `flutter_secure_storage`(access) + `PersistCookieJar`(refresh)。**バックエンド改修ゼロ**。**Firebase Auth は使わない**（ID 二重管理回避） |
| デザイン | Liquid Glass(iOS 26) + ニュートラル基調 + **accent 1トークン**（`ColorScheme.primary`、既定ブランド赤、差し替え可） | プレビュー(`mobile/preview/`)準拠 |
| Firebase | Crashlytics / Analytics のみ | 可観測性。FCM はバックエンド通知基盤が無いため後続 |
| CI/CD | GitHub Actions（analyze/test ゲート + `testflight.yml` で iOS 署名/TestFlight） | — |
| テスト | flutter_test（unit/widget）+ MagicPod（E2E） | refresh 競合制御など要所を単体テスト |

### API クライアントの生成方針

当初 `swagger_parser` での自動生成を検討したが、Dart 3.11 / Freezed 3 系の最新環境での生成物の安定性とコンパイル信頼性を優先し、**実装する範囲の DTO は Freezed で手書き**する。スキーマの単一ソースは引き続き `server/docs/swagger.yaml`。将来、生成の安定性が確認できれば codegen へ移行可能（DTO は同形）。

### iOS 課金（重要）

App Store Review Guideline 3.1.1 により、**アプリ内に購入導線・外部購入リンクを置かない**。Pro は Web 購入のみとし、アプリは `GET /purchase/subscription` の状態を読むだけ（Pro 機能の解放判定）。StoreKit IAP + サーバ側レシート/Server Notifications V2 検証は **Phase 3** のスコープ。

## 代替案

- **React Native / Swift ネイティブ**: Flutter は Web のデザイン/機能を最短で移植でき、単一コードベースで iOS/Android 将来対応も可能なため採用。
- **Firebase Auth 全面採用**: バックエンドが既に JWT で ID を所有しており二重管理になるため不採用。

## 影響

- `mobile/` を新設。既存の `server/` `web/` と同 repo・同 PR フロー。
- 認証はバックエンド無改修で流用可能。
- Phase 3 で StoreKit / FCM のためのサーバ側追加実装（レシート検証・device token・配信）が必要。
- Crashlytics/Analytics/IAP/FCM の最終疎通には Apple Developer・Firebase プロジェクト・実機が必要（手順は `mobile/README.md`）。
