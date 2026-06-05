# ADR 0008: RSC(React Server Components)移行の検討

## ステータス
保留 (2026-06-05) — 検討のみ、 実装は別 ADR で

## コンテキスト

`web/` は Next.js 16 App Router を使っているが、 **全ページが
`"use client"`** で実装されている。 これは Pages Router 時代の
クライアント主導 SPA に近い構造で、 App Router の主要メリット
(Server Components / Streaming SSR / 部分プリレンダー / RSC
Payload 経由のデータ転送)を **全部捨てている** 状態。

ADR 0007 で Next.js を Vercel + Fargate ではなく **Vercel のみ**
で運用する方針になり、 ホスティング側の RSC サポートは整っている。
コード側だけが Client-only パターンに張り付いている。

本 ADR は 「移行すべきか / すべきなら何が前提か / どこから始めるか」
を整理して、 実装着手前に意思決定する。

## 現状の特徴

| 観点 | 現状 |
|---|---|
| ページ実装 | 全 `"use client"` (login / register / meals / exercises / trainings / routines) |
| データ取得 | TanStack Query（クライアント） |
| 認証 | localStorage 互換の in-memory access token + refresh token (HttpOnly cookie) |
| ナビゲーション | `next/link` の client-side routing |
| SEO | ほぼ不要(ログイン必須 SaaS、 検索流入なし) |
| 初回表示 | Vercel の static prerender + client hydration |

## 検討する移行案

### 案 A: RSC + Cookie 認証への全面移行

サーバー側で `cookies()` から `refresh token` を読み、 BE を直接 fetch して
ページを SSR する。 page.tsx は async Server Component に書き換える。

**Pros**
- 初回表示が速くなる(JS bundle 待ちなし)
- BE トークン管理がサーバー側で完結 → XSS で access token を盗まれない
- TanStack Query を一部不要にできる(Server Component 内で fetch すれば良い)

**Cons**
- **BE 認証フローを Cookie ベースに変更必須**
  - 現状: `Authorization: Bearer <access_token>` (JS 管理)
  - 移行後: BE が `HttpOnly Cookie` 経由でも access token を受けるか、
    refresh cookie からの軽量再発行を server fetch で行う
- 認証セッションが「クライアント完結」 から 「サーバー経由」 に変わるので、
  CSRF 対策が必要(SameSite=Strict + Origin header 検証)
- React Query との二重管理が発生する(mutation はクライアント、 query は
  サーバー → invalidate 周りが複雑化)
- 既存ページ全部を async function に書き換え + Client Component を分離

### 案 B: 部分移行(エントリーページのみ RSC)

`/` (login への redirect)、 `/login`、 `/register` のような認証前ページだけ
Server Component 化し、 認証後のページは Client Component のまま残す。

**Pros**
- BE 改修不要(認証後のページは現状通り)
- SEO 必要なページ(ランディング、 SaaS 紹介)を追加するときに RSC を使える
- 投資が小さくリターンが計測しやすい

**Cons**
- 認証後のページの初回表示速度は改善しない
- 「半分 RSC 半分 Client」 の二重メンタルモデルが残る

### 案 C: 現状維持(SPA パターン)

`"use client"` 全乗せのまま、 React Query で完結させる。

**Pros**
- 追加開発コストゼロ
- メンタルモデルがシンプル(全部クライアント)
- Vercel の static prerender だけでも体感は十分速い

**Cons**
- App Router の旨味を出しきれない
- SODA 入社後 「なぜ App Router を選んだのに RSC を使わなかったか」 を
  説明する必要が出る可能性
- 将来 SEO 要件が出たら案 A or B にどのみち移行する

## 判断

**現時点では案 C(現状維持)を採用**。 ただし以下のトリガーで再検討する：

1. **SEO が必要なページが追加される**(ランディング、 ブログ等) → 案 B
2. **初回表示速度の改善要求が出る**(モバイル LCP 計測で 2.5s 超え等) → 案 A
3. **BE 認証フローを Cookie ベースに変更する別の動機が出る** → 案 A

## 理由

### a. 案 A を今やらない理由

- BE 改修が大きい:
  - `internal/auth/handler/auth_handler.go` で Cookie 経由の access token を扱う
  - CSRF middleware を追加
  - `internal/shared/middleware/auth_middleware.go` で Cookie / Bearer 両対応
- 現状の SPA + React Query で **機能要件は満たしている**
- 認証必須 SaaS なので **SEO の必要性がゼロに近い**
- 入社後の業務時間で覚える価値は高いが、 入社前の限られた時間を割く優先度は低い

### b. 案 B を今やらない理由

- 認証前ページ(login / register)は静的にプリレンダーされる現状で十分高速
- 部分移行は 「二重メンタルモデル」 のコストを払いつつ、 リターンは限定的
- ランディングページが必要になった時に追加で `app/(marketing)/` 配下に
  RSC で書き始めれば良い

### c. 案 C を採用する理由

- 直近で達成したい目的(SODA 練習、 本番 AWS デプロイ、 ドメイン取得)は
  全て現状アーキテクチャで満たせる
- RSC 移行の **判断材料を貯める時間** が確保できる
  (案 A をやるなら BE 認証設計から見直すのが筋)
- 後から移行しても致命的に難しくならない構造を 0001 / 0008(本 ADR)で
  既に整理した

## 移行する場合の前提条件(将来のためのメモ)

案 A を将来やる場合、 以下を **本 ADR より先に** 解決すること:

1. **BE の認証フローを Cookie 対応にする ADR**
   - access token を HttpOnly Cookie で受けるオプションを追加
   - CSRF 対策(SameSite=Strict + Origin header)
   - `refresh` endpoint の挙動を Server Component から呼べる形に整理

2. **React Query の「サーバー側プリフェッチ」 戦略を決める**
   - `dehydrate` / `HydrationBoundary` パターン
   - どのページで使うかの選定基準

3. **Client / Server Component の分割境界を文書化**
   - 「データ取得は Server」 「インタラクションは Client」 が基本
   - features/ 配下でどう構造化するか(例: `features/training/api/server.ts`
     に Server-side fetchers を分離)

4. **Vercel 上の動作検証**
   - RSC payload の転送サイズ
   - 初回表示と navigation の性能比較
   - Cloudwatch / Vercel Analytics で計測

## 不採用案

D. **Pages Router に戻す**
   - 棄却: App Router の他の利点(file-system routing for nested layouts,
     loading.tsx / error.tsx, parallel routes 等)は使えており、 戻すコスト >
     メリット。

E. **RSC + JWT を localStorage(現状互換)で運用**
   - 棄却: Server Component は localStorage を読めない。 JWT を Cookie に
     入れる時点で結局 Cookie 対応が必要 → 案 A に合流する。

F. **Remix / TanStack Start に移行**
   - 棄却: フレームワーク変更は学習コストと書き換えコストが Next.js 16 RSC
     移行の 5 倍以上。 SODA の練習目的にも噛み合わない。

## 結果

- 現状の `"use client"` 全乗せ構造を継続
- features-based のディレクトリ構成は **将来の RSC 移行でも崩れない設計**
  になっており、 移行コストは features 単位で段階的に払える
- 上記トリガー発生時に本 ADR を再評価 → 新 ADR で実装計画を立てる
