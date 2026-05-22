# ADR 0001: Monorepo を採用する

## ステータス
採用 (2026-05-23)

## コンテキスト
個人開発 SaaS の構成として、 Monorepo / Polyrepo を検討。
将来的に server / web / mobile の3クライアントを抱える想定。

## 決定
**Monorepo** を採用。

## 理由
1. Connect-RPC の proto を 3クライアント(server/web/mobile)で共有しやすい
2. ADR / 設計判断が一箇所に集約
3. 個人開発で3リポジトリ管理はオーバーヘッド過大
4. GitHub Actions の path filter で CI を分離可能

## 結果
- ✅ proto 共有が単純(同一ディレクトリ参照)
- ✅ 設計と実装が同一リポジトリで完結
- ❌ CI/CD は path filter で工夫が必要
- ❌ Polyrepo 固有の経験(buf schema registry 等)は別途必要
  → SODA 入社後に習得想定
