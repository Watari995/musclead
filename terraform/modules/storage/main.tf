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
# 2. Public Access Block: 「意図しない公開」 を完全防御
#    4 つすべて true = 公開 ACL も Bucket Policy も完全拒否
#    → presigned URL でしかアクセスできない状態を保証
# ─────────────────────────────────────────────────────────
resource "aws_s3_bucket_public_access_block" "images" {
  bucket = aws_s3_bucket.images.id

  block_public_acls       = true # public ACL を持つ object 作成を禁止
  block_public_policy     = true # public な bucket policy 自体を禁止
  ignore_public_acls      = true # 既に public ACL がある object でも無視
  restrict_public_buckets = true # public policy が万一あっても拒否
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
# 4. CORS 設定: ブラウザからの直接アクセスを許可
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
