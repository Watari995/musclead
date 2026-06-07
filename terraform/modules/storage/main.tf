# 画像保存用 S3 bucket とアクセス制御
#
# 用途: ユーザーのプロフィール画像、 種目アイコン、 食事写真など
# 構成: 1 bucket 内に prefix で type を分離(profiles/、 exercises/ など)
#
# 設計方針:
#   - 完全 private (public access block 全有効)
#   - 平文保存はせず SSE-S3 で暗号化(置いた瞬間に暗号化される)
#   - CORS: ブラウザから presigned URL 経由の PUT/GET を許可
#   - バージョニングは無効(コスト削減、 image は上書き不要 = 新 key で運用)

# ─────────────────────────────────────────────────────────
# 1. S3 bucket 本体
#    実体は「ファイル置き場」、 名前は全 AWS で一意である必要
# ─────────────────────────────────────────────────────────
resource "aws_s3_bucket" "images" {
  # bucket 名 = "musclead-images-<account_id>"
  # account_id を含めることで衝突回避
  bucket = "musclead-images-${var.account_id}"

  tags = {
    Name = "musclead-images"
  }
}

# ─────────────────────────────────────────────────────────
# 2. Public Access Block: ACL ベースの公開は禁止、 Policy ベースは許可
#    profiles/ 配下のアバター画像は Policy で public read にする(下参照)
#    ACL ベースの object 単位公開は引き続き禁止(誤公開防止)
# ─────────────────────────────────────────────────────────
resource "aws_s3_bucket_public_access_block" "images" {
  bucket = aws_s3_bucket.images.id

  block_public_acls       = true  # public ACL を持つ object 作成を禁止(ACL での誤公開防止)
  block_public_policy     = false # Policy ベースの public 許可は受け付ける(profiles/ 公開のため)
  ignore_public_acls      = true  # 既に public ACL がある object でも無視
  restrict_public_buckets = false # Policy で public 指定された prefix は実際に公開する
}

# ─────────────────────────────────────────────────────────
# 2.1 Bucket Policy: profiles/ 配下だけ public read 許可
#     アバター用途で公開 URL での直接配信を可能にする(CDN ライク)。
#     他 prefix(将来追加されるもの)は引き続き private のまま。
# ─────────────────────────────────────────────────────────
resource "aws_s3_bucket_policy" "images_public_profiles" {
  bucket = aws_s3_bucket.images.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Sid       = "PublicReadProfiles"
      Effect    = "Allow"
      Principal = "*"
      Action    = "s3:GetObject"
      Resource  = "${aws_s3_bucket.images.arn}/profiles/*"
    }]
  })

  # public_access_block を緩める更新が反映される前に policy を入れると拒否されるので
  # 明示的に依存させる
  depends_on = [aws_s3_bucket_public_access_block.images]
}

# ─────────────────────────────────────────────────────────
# 3. サーバー側暗号化 (SSE-S3)
#    AWS 管理の鍵で「ディスク書き込み時に自動暗号化、 読む時に自動復号」
#    料金: 無料、 性能影響なし、 設定しない理由がない
# ─────────────────────────────────────────────────────────
resource "aws_s3_bucket_server_side_encryption_configuration" "images" {
  bucket = aws_s3_bucket.images.id

  rule {
    apply_server_side_encryption_by_default {
      # AES256 = SSE-S3 (AWS 管理鍵)
      # KMS は鍵管理ができる代わりに $1/月 + リクエスト課金、 今回 overkill
      sse_algorithm = "AES256"
    }
  }
}

# ─────────────────────────────────────────────────────────
# 4. デフォルトプロフィール画像を S3 にアップロード
#    アバター未設定 / 削除時の fallback として全 user が参照する。
#    ソースは repo 同梱 (terraform/modules/storage/assets/default.png)、
#    Terraform apply のタイミングで S3 に push される。
#    画像差し替え時は repo の PNG を上書きして apply するだけ。
# ─────────────────────────────────────────────────────────
resource "aws_s3_object" "default_profile_image" {
  bucket       = aws_s3_bucket.images.id
  key          = "profiles/default.png"
  source       = "${path.module}/assets/default.png"
  content_type = "image/png"

  # ソースの中身が変わった時だけ差分検出 → 自動 upload
  # filemd5 はファイル内容のハッシュ
  etag = filemd5("${path.module}/assets/default.png")
}

# ─────────────────────────────────────────────────────────
# 5. CORS 設定: ブラウザからの直接アクセスを許可
#    BE 経由ではなくブラウザが直接 S3 に PUT/GET するので、
#    S3 側でも CORS preflight を通す必要がある
# ─────────────────────────────────────────────────────────
resource "aws_s3_bucket_cors_configuration" "images" {
  bucket = aws_s3_bucket.images.id

  cors_rule {
    # 許可する origin (FE)
    # 開発時は localhost も足したい場合はここに追加
    allowed_origins = var.allowed_origins

    # PUT: presigned URL でアップロード時
    # GET: presigned URL で画像表示時
    allowed_methods = ["PUT", "GET", "HEAD"]

    # presigned URL に含まれる X-Amz-* ヘッダ、 Content-Type 等を許可
    allowed_headers = ["*"]

    # ブラウザに公開するレスポンスヘッダ
    expose_headers = ["ETag"]

    # preflight 結果のキャッシュ秒数 (1 時間)
    max_age_seconds = 3600
  }
}
