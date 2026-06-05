# GitHub Actions が AWS リソースを操作するための OIDC 認証
#
# 仕組み(なぜ鍵いらないか):
#   1. GitHub Actions が workflow 実行時に「私は Watari995/musclead リポ
#      の main ブランチで動いている workflow です」 という署名付き token
#      (OIDC token = JWT) を発行
#   2. workflow 内で「この token で AWS の Role を引き受けたい」 と要求
#   3. AWS が token を GitHub に問い合わせて検証(GitHub の公開鍵で署名検証)
#   4. 検証 OK → 一時的な credentials (1時間有効) を発行
#   5. workflow がそれで AWS API を叩く
#
# 静的キー方式との違い:
#   - 静的: アクセスキーを GitHub Secrets に保存 → 漏洩リスク + rotate 必要
#   - OIDC: 鍵を持ち運ばない、 都度発行 → 漏洩リスクゼロ + rotate 不要

# ─────────────────────────────────────────────────────────
# 1. OIDC Provider: GitHub と AWS の信頼関係を AWS アカウントに登録
#    アカウントに 1 つ作れば、 全リポ・全 module から共通で使える
# ─────────────────────────────────────────────────────────
resource "aws_iam_openid_connect_provider" "github" {
  # GitHub OIDC の発行元 URL(固定値)
  url = "https://token.actions.githubusercontent.com"

  # client_id_list = audience
  # 「この OIDC を使う先は AWS STS だよ」 という宣言(固定値)
  client_id_list = ["sts.amazonaws.com"]

  # thumbprint: GitHub の TLS 証明書の指紋
  # AWS が「GitHub からの応答である」 と検証するために使う
  # → AWS 公式ドキュメントの推奨値(2 つ並べて将来の証明書ローテに耐える)
  thumbprint_list = [
    "6938fd4d98bab03faadb97b34396831e3780aea1",
    "1c58a3a8518e8759bf075b76b750d4f2df264fcd",
  ]
}

# ─────────────────────────────────────────────────────────
# 2. IAM Role: GitHub Actions が引き受ける(assume する) Role
#    workflow がこの Role の権限で AWS を操作する
# ─────────────────────────────────────────────────────────
resource "aws_iam_role" "github_actions" {
  name = "musclead-github-actions"

  # Trust policy: 「誰が」 この Role を引き受けられるか
  # ここで「Watari995/musclead の main ブランチからの workflow だけ」 と限定
  # → これを緩めると他リポからも乗っ取られるので最重要設定
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"

      # 引き受け元 = 上で作った OIDC Provider
      Principal = {
        Federated = aws_iam_openid_connect_provider.github.arn
      }

      # OIDC 経由の Role 引き受け Action
      Action = "sts:AssumeRoleWithWebIdentity"

      # 条件: どんな OIDC token なら引き受けを許すか
      Condition = {
        StringEquals = {
          # audience は AWS STS(固定)
          "token.actions.githubusercontent.com:aud" = "sts.amazonaws.com"
        }
        StringLike = {
          # sub claim で「どのリポ・どのブランチからか」 を制限
          # 形式: repo:<org>/<repo>:ref:refs/heads/<branch>
          # 例: repo:Watari995/musclead:ref:refs/heads/main
          "token.actions.githubusercontent.com:sub" = "repo:${var.github_repo}:ref:refs/heads/${var.allowed_branch}"
        }
      }
    }]
  })
}

# ─────────────────────────────────────────────────────────
# 3. この Role に「何ができるか」 の権限を付与
#    最小権限: ECR push + ECS deploy + PassRole のみ
#    → workflow から「DB 削除」 等の余計な操作はできない
# ─────────────────────────────────────────────────────────
resource "aws_iam_role_policy" "github_actions" {
  name = "musclead-github-actions-deploy"
  role = aws_iam_role.github_actions.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      # ── (1) ECR token 取得 ─────────────────
      # ECR への docker login に必要
      # この権限は ARN 単位で絞れない仕様なので Resource = "*" になる
      {
        Effect   = "Allow"
        Action   = "ecr:GetAuthorizationToken"
        Resource = "*"
      },

      # ── (2) ECR push / read(musclead-server repo のみ) ─────────────────
      # docker push でレイヤー + マニフェストをアップロード、
      # Buildx のマルチアーキビルドでは既存マニフェスト確認(HEAD)も発生する
      {
        Effect = "Allow"
        Action = [
          "ecr:BatchCheckLayerAvailability", # 既存レイヤー確認
          "ecr:InitiateLayerUpload",         # アップロード開始
          "ecr:UploadLayerPart",             # チャンク送信
          "ecr:CompleteLayerUpload",         # 完了通知
          "ecr:PutImage",                    # マニフェスト登録
          "ecr:BatchGetImage",               # マニフェスト HEAD / pull(Buildx に必須)
          "ecr:DescribeImages",              # image 情報取得
        ]
        Resource = var.ecr_repository_arn
      },

      # ── (3) ECS deploy ─────────────────
      # Task Definition の新 revision 登録 + Service 更新
      # ECS の Task Def は ARN が rev 番号で変わるので Resource = "*"
      {
        Effect = "Allow"
        Action = [
          "ecs:DescribeServices",       # 現在の Service 状態確認
          "ecs:UpdateService",          # 新 Task Def に切替
          "ecs:DescribeTaskDefinition", # 現 Task Def の中身取得
          "ecs:RegisterTaskDefinition", # 新 revision 登録
        ]
        Resource = "*"
      },

      # ── (4) PassRole ─────────────────
      # ECS Task Def を登録する時、 「この Task は <execution role> を使います」
      # と宣言する。 そのために「この role を pass(渡す) する権限」 が必要
      # → 該当 role 以外を pass できないよう ARN を絞る
      {
        Effect   = "Allow"
        Action   = "iam:PassRole"
        Resource = var.task_execution_role_arn
      },
    ]
  })
}
