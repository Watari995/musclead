# E2E (mobile UI guard)

iPhone SE viewport (375 × 667) でモバイル UI 崩れを自動検出する Playwright スイート。
詳細な背景は [.cursor/rules/11-mobile-responsive.mdc](../../.cursor/rules/11-mobile-responsive.mdc) を参照。

## 検査内容

| spec | 検査 |
|---|---|
| [`specs/mobile-overflow.spec.ts`](./specs/mobile-overflow.spec.ts) | 水平スクロール + 各要素が viewport の左右に出ていないか |
| [`specs/a11y.spec.ts`](./specs/a11y.spec.ts) | axe-core で `critical` / `serious` の WCAG 2.1 AA 違反 |

両 spec ともに [`helpers/pages.ts`](./helpers/pages.ts) の `TARGET_PAGES` を巡回。

## 新ページを追加したら

1. `helpers/pages.ts` の `TARGET_PAGES` に 1 行追加
2. Popover やドロワーなど「開かないと出てこない UI」は `interact:` でその操作を書く
3. ローカルで `npm run e2e` が通ることを確認

## ローカル実行

事前にスタックを起動しておく:

```bash
# repo root で
make db-up           # MySQL
make migrate-up      # スキーマ適用
make run &           # API server (別ターミナル推奨)
(cd web && npm run dev &)   # web
```

その後:

```bash
cd web
npm run e2e:install   # 初回のみ: Playwright Chromium をインストール
npm run e2e           # 実行
npm run e2e:ui        # UI モードでデバッグ
```

失敗時は `web/playwright-report/index.html` を開いて trace / screenshot / video を確認。

## CI

`.github/workflows/e2e-web.yml` で `web/`・`server/`・`sql/migrations/` 変更時に自動実行。

## 認証

`auth.setup.ts` が POST `/users` でランダム email のユーザを作成し、 `/login` 経由で
`e2e/.auth/user.json` に storageState を保存。 各テストはこの cookie で `AuthBootstrap`
の refresh を通過し、 認証済み状態で起動する。

`e2e/.auth/` は gitignore。
