# ADR 0002: Connect-RPC を採用する

## ステータス
採用 (2026-05-23)

## コンテキスト
API レイヤーの選定。 候補: REST + JSON / GraphQL / gRPC / Connect-RPC

## 決定
**Connect-RPC** を採用。

## 理由
1. メソッド名 = ユースケース名 で DDD フレンドリー
2. 言語間でスキーマを対等に共有(Protobuf)
3. HTTP/1.1 でも HTTP/2 でも動く(curl で叩ける)
4. SODA で実務利用想定、 学習価値高い

## 結果
- ✅ API スキーマが Protobuf で型安全に管理
- ✅ Connect-Web で TS 型自動生成
- ❌ REST より学習コストやや高い(許容範囲)
