# ADR 0009: プロフィール画像の保存設計(NOT NULL + default + public read + DELETE endpoint)

## ステータス
採用 (2026-06-07)

## コンテキスト

ユーザープロフィール画像のアップロード/表示/削除機能を追加する。

要件:
- アップロード: ブラウザから S3 へ直接(presigned URL 経由、 BE は帯域経由しない)
- 表示: BE が presigned GET URL を発行、 FE は `<img src>` で表示
- 削除: 「自分の画像を削除して default に戻す」 操作を提供

S3 bucket / Task Role / presigned URL 発行基盤は ADR 0007 系列で構築済(私たち内部)。

本 ADR は **DB / API / ドメインレベルの設計判断** を残す。
具体的には:
1. `profile_image_path` カラムは NULL を許すか
2. 「削除」 をどう表現するか(PATCH に統合 vs 独立 DELETE)

## 判断

### ⓪ 表示用 URL: **public bucket policy で素の URL**

S3 bucket の `profiles/` 配下だけ public read を許可し、 アバター画像は **`https://<bucket>.s3.<region>.amazonaws.com/profiles/...` の素の URL で直接配信** する。

```
profile_image_url = "https://musclead-images-204340689570.s3.ap-northeast-1.amazonaws.com/profiles/u1/abc.jpg"
                                                                                          ↑ 認証なし、 永続 URL
```

BE は **path に prefix を付けるだけ**(presigned GET URL は不要)。 アップロード(PUT)は引き続き presigned URL で認証する。

### ① スキーマ: **NOT NULL + DEFAULT 'profiles/default.png'**

```sql
ALTER TABLE users
  ADD COLUMN profile_image_path VARCHAR(255) NOT NULL DEFAULT 'profiles/default.png'
  AFTER birthday;
```

- 全 user に必ず path が入る
- 登録時にも `'profiles/default.png'` が自動セット
- default 画像本体は Terraform `aws_s3_object` で S3 に upload(repo の `terraform/modules/storage/assets/default.png` をソースに)

### ② 削除 API: **独立 DELETE endpoint**

```
DELETE /users/me/profile-image
```

実体は「path を `'profiles/default.png'` に戻す + 旧 S3 object を削除」。
PATCH は通常の field 更新専用に保つ。

### ③ ドメインモデル

```go
type User struct { ...; profileImagePath string }  // *string ではなく string

func (u *User) ProfileImagePath() string
func (u *User) SetProfileImagePath(path string)
```

nil の概念は持ち込まない。

## なぜ presigned GET URL ではなく public read を採用したか

代替案: bucket を完全 private にして、 BE が表示時に都度 presigned GET URL を発行する。

却下理由:
- BE 処理が毎回必要(URL 構築よりはるかに重い)
- URL が毎回変わるためブラウザ・CDN キャッシュが効かない → 表示遅延
- 業界標準(GitHub `avatars.githubusercontent.com` / Twitter / Slack)は public URL
- アバター画像は機密でない(本人が公開を選択している)
- UUID v7 ベースの path で推測不可

採用案(public read):
- BE は単純な string 連結のみ(`<bucket-url>/<path>`)
- URL 不変 → ブラウザでキャッシュ
- 公開範囲は `profiles/*` だけに限定(将来追加する非公開 prefix は引き続き private)
- アップロード(PUT)は presigned URL で認証する(書き込みは絞る)

## なぜ NULL 許可ではなく NOT NULL を採用したか

検討段階で「nullable + BE で default URL に変換」 案も俎上に乗ったが、 以下の理由で却下:

| 観点 | nullable | **NOT NULL + default** |
|---|---|---|
| API レスポンス | `profile_image_url: string \| null` | 常に `string`(FE はチェック不要) |
| ドメインの型 | `*string`(nil リスク) | `string`(シンプル) |
| BE コードパス | 各所で `if path == nil` 分岐 | 常に GenerateGetURL 1 行 |
| ドメイン真実 | 「user は avatar を**持つかもしれない**」 | 「user は **必ず avatar を持つ**(default 含む)」 |
| 業界 alignment | 少数派 | Twitter / GitHub / Slack / Notion 流儀 |
| schema 厳密性 | NULL を許す = 不整合余地 | 制約で強制 |

決定的だったのは **API 契約の一貫性** と **ドメインの真実認識**:
- 「default 画像は欠落値の代替」 ではなく、 **「default 画像も正規の avatar の 1 つ」**
- このフレーミングなら nullable は不自然 → NOT NULL 自然

## なぜ削除を PATCH ではなく DELETE にしたか

代替案: PATCH に統合し `{ "profile_image_path": null }` を「default 戻し」 と解釈する。

却下理由:
- null の意味が field 毎に違うと FE 開発者が混乱
  (birthday: 値クリア、 profile_image: default 復元)
- PATCH の責務が「設定」 と「削除→ default 戻し」 で重くなる
- 業界標準は独立 DELETE(GitHub の `DELETE /user/avatar` 等)

→ DELETE endpoint で意図を URL から明確にする方が綺麗。

## 影響

### スキーマ
- `users` テーブルに `profile_image_path VARCHAR(255) NOT NULL DEFAULT 'profiles/default.png'`
- 既存 user は DEFAULT 句のおかげで自動的に default で埋まる(migration 1 段階)

### API
```
POST   /users/me/profile-image/presigned-url   既存(発行)
PATCH  /users/me { profile_image_path: "..." } 既存(設定)
DELETE /users/me/profile-image                 新規(default に戻す)
GET    /users/me                               既存(URL を必ず返す)
```

### BE
- `User` domain: `profileImagePath string`(non-pointer)
- usecase: `DeleteProfileImage` 新規
- 共通ロジック: `if oldPath != defaultProfileImagePath { storageClient.DeleteObject(oldPath) }`
- default path 定数: `const defaultProfileImagePath = "profiles/default.png"`

### インフラ
- `terraform/modules/storage` に `aws_s3_object "default_profile_image"` 追加
- `terraform/modules/storage/assets/default.png` を repo に commit

### FE
- 「画像を選んでアップロード」 ボタン: 既存 PATCH フローでアップロード後 path 設定
- 「画像を削除」 ボタン: `DELETE /users/me/profile-image` 呼ぶ
- 表示: `<img src={user.profile_image_url}>` のみ、 null 考慮不要

## やらないこと

- gravatar / identicon の生成は採用しない(運営側で 1 枚の default を共有)
- 画像のサイズ/形式バリエーション(thumbnail 等) は当面なし(必要時に別 ADR)
- ban / 非公開 user 等の特殊 case 対応も別途

## 関連 ADR

- [ADR 0007](0007-infra-mvp-and-monorepo.md): インフラ / monorepo 構成
