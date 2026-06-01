# ADR 0002: DDD + Modular Monolith を採用する

## ステータス
採用 (2026-05-23)

## コンテキスト
アーキテクチャ設計の選定。 候補: Layered / DDD / Hexagonal / Modular Monolith / Microservices

## 決定
**DDD + Modular Monolith** を採用。

## 理由
1. Bounded Context 単位でモジュール分離(User / Meal / Training / Weight)
2. 単一バイナリで運用が単純(個人開発)
3. 将来マイクロサービス分割の余地を残せる
4. SODA で実務利用想定、 学習価値高い

## 結果
- ✅ ドメインの境界が明確
- ✅ 単一プロセスで開発・デプロイが楽
- ❌ モジュール間依存の規律が必要(linter で防御)
