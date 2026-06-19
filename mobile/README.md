# musclead — iOS (Flutter)

筋トレ・食事・体重 一元管理 SaaS の iOS ネイティブアプリ。設計判断は [ADR-0021](../docs/adr/0021-adopt-flutter-for-ios.md)。

## 必要環境

- [fvm](https://fvm.app/) — Flutter バージョンは `.fvmrc` に固定（**3.41.7 / Dart 3.11.5**）
- Xcode（iOS ビルド時）

## セットアップ

```bash
cd mobile
fvm install                 # .fvmrc のバージョンを用意
fvm flutter pub get
fvm dart run build_runner build   # Freezed / json_serializable 生成（生成物は gitignore）
```

## 実行

```bash
# dev（API: http://localhost:8080）
fvm flutter run --dart-define=FLAVOR=dev

# 本番 API（https://api.musclead.com）
fvm flutter run --dart-define=FLAVOR=prod
```

base URL は `--dart-define=API_BASE_URL=...` で上書き可能（`lib/core/config/app_config.dart`）。

## コード生成

DTO は Freezed + json_serializable。スキーマの単一ソースは `server/docs/swagger.yaml`（`field_rename: snake` でバックエンドの snake_case に整合）。

```bash
fvm dart run build_runner watch   # 変更を監視して再生成
```

## テスト

```bash
fvm flutter test            # unit + widget
fvm flutter test --coverage
```

E2E は MagicPod（実機 / シミュレータ）。

## アーキテクチャ

Feature-first + 軽量3層。

```
lib/
  core/        config / theme(Liquid Glass + accent) / api(dio+interceptor) /
               auth / router(go_router) / error / util / widgets / providers
  features/<x>/ data(DTO+Repository) / application(Riverpod) / presentation(Screens)
  bootstrap/   Firebase 初期化
```

- 状態管理: Riverpod v2（手書き Notifier）+ hooks_riverpod
- 認証: バックエンド JWT(access) は `flutter_secure_storage`、refresh は HttpOnly Cookie を `PersistCookieJar` で永続化。401 は interceptor が refresh→再試行（競合は単一化）
- デザイン: Liquid Glass + ニュートラル基調 + アクセント1トークン（`accentProvider` で差し替え可）。UI プレビュー: `preview/index.html`

## Firebase（Crashlytics / Analytics）

実プロジェクトは各自で用意する（リポジトリには機密を含めない）。

```bash
dart pub global activate flutterfire_cli
flutterfire configure            # GoogleService-Info.plist と lib/firebase_options.dart を生成（gitignore 済）
```

未設定でもアプリは起動する（`bootstrap/firebase_bootstrap.dart` が例外を握りつぶす）。FCM/Push はバックエンド通知基盤が整ってから（Phase 3）。

## 課金（重要 / Phase 3）

App Store Review 3.1.1 のため、**アプリ内に購入導線・外部購入リンクを置かない**。Pro は Web 購入のみで、アプリは `GET /purchase/subscription` の状態を読むだけ。StoreKit IAP（`in_app_purchase`）とサーバ側レシート/Server Notifications V2 検証は Phase 3。

## CI/CD

- GitHub Actions（`.github/workflows/ci-mobile.yml`）: PR で `analyze` + `test`
- GitHub Actions（`.github/workflows/testflight.yml`）: `workflow_dispatch` で version 指定 → iOS 署名 → TestFlight → App Store 審査提出。必要な Secrets は `IOS_DIST_CERT_BASE64` / `IOS_DIST_CERT_PASSWORD` / `IOS_PROFILE_BASE64` / `ASC_KEY_ID` / `ASC_KEY_BASE64` / `ASC_ISSUER_ID` / `SENTRY_DSN`

## iOS ローカルビルドの前提

`flutter build ios` は CocoaPods を使う。ローカル環境で初回ビルドする前に、以下を満たすこと（GitHub Actions の self-hosted runner では不要）。

```bash
gem install rexml          # xcodeproj が要求（環境により未導入のことがある）
pod repo update            # もしくは CocoaPods CDN を利用
```

Podfile は Firebase iOS SDK 要件に合わせ `platform :ios, '15.0'`（deployment target も 15.0 に統一済み）。
