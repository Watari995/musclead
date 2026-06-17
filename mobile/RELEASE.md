# 本番配信手順 (iOS / App Store)

> ルート: **Xcode で Archive → App Store Connect へ Upload**。
> 私（AI）側のリリース準備は完了済み。以下はあなたが Apple のコンソール/Xcode で行う作業。

## 前提
- Apple Developer Program 加入済み（有料）。Xcode の署名 team = `44675MS8C4` 設定済み・自動署名。
- `mobile/.fvmrc` の Flutter 3.41.7（`fvm`）。

## 準備済み（このリポジトリ）
- bundle id: `com.musclead.musclead` / 表示名: `musclead` / version: `1.0.0 (1)`（`pubspec.yaml` の `version: 1.0.0+1`）
- **release ビルドは自動で本番 API（`https://api.musclead.com`）に向く**（`AppConfig` が `kReleaseMode`→prod）
- 写真権限の説明文（Info.plist）/ 暗号化非該当申告（`ITSAppUsesNonExemptEncryption=false`）/ アプリアイコン一式
- アプリ内 **アカウント削除**（マイページ）= App Store 要件 5.1.1(v) 対応

### （任意）正式アイコンに差し替え
`mobile/assets/icon/app_icon.png` を 1024×1024 PNG に置換 →
```bash
cd mobile && fvm dart run flutter_launcher_icons
```

## 手順

### 1. App Store Connect でアプリ作成
1. https://appstoreconnect.apple.com → マイApp → ＋ → 新規App
2. iOS / 名前「musclead」/ 言語 日本語 / **bundle id `com.musclead.musclead`**（一覧に無ければ [Developer portal](https://developer.apple.com/account/resources/identifiers/list) で App ID 登録）/ SKU 任意

### 2. Xcode で Archive → Upload
```bash
cd mobile
fvm flutter build ios --release        # 署名なしで一度ビルド（pods/設定を反映）
open ios/Runner.xcworkspace             # ★ .xcworkspace を開く（.xcodeproj ではない）
```
Xcode で:
1. デバイスを **Any iOS Device (arm64)** に
2. **Product ▸ Archive**
3. Organizer ▸ **Distribute App ▸ App Store Connect ▸ Upload**（自動署名のまま進む）
4. 完了 → App Store Connect の TestFlight に「処理中」で出る（数分）

> CLI 派なら `fvm flutter build ipa` でも可。生成 ipa を Transporter / `xcrun altool` でアップロード。

### 3. TestFlight（実機テスト）
1. App Store Connect ▸ TestFlight ▸ 内部テスターに自分を追加
2. iPhone の TestFlight アプリでインストール（輸出コンプライアンスは「非該当」で自動回答済み）

### 4. App Store 審査へ提出（一般公開する場合）
1. 「App情報」「価格（無料）」を入力
2. **Appプライバシー**を申告: メール / プロフィール / トレーニング・食事・体重データを「アカウントに紐付けて収集」、**トラッキングなし**
3. スクリーンショット（6.7"・6.5" 必須）/ 説明 / キーワード / サポートURL
4. ビルドを選択 → 審査へ提出

## 審査の注意（本アプリ固有）
- **課金**: アプリ内に購入導線なし・サブスク状態を読むだけ（Guideline 3.1.1 準拠）。IAP を入れる場合は StoreKit + サーバ受信検証が別途必要。
- **Sign in with Apple**: 自前のメール/パスワードのみで第三者ソーシャルログイン不使用のため不要。
- **アカウント削除**: マイページに実装済み。
- **プライバシーマニフェスト**: Flutter/各プラグインが同梱。アップロードで required-reason API を指摘された場合のみ `Runner/PrivacyInfo.xcprivacy` を追加。

## 次回リリース時
`pubspec.yaml` の `version:` を上げる（例 `1.0.1+2`）→ 再 Archive。`+N` が build 番号（App Store Connect で重複不可）。
