---
name: ship
description: musclead の本番リリース workflow を自動化。 feature ブランチから develop PR(スクリーンショット付き) → main への deploy PR → タグ + GitHub Release → 本番動作確認 までを 1 コマンドで実行。
argument-hint: [all (省略時) | feature | deploy | release | verify]
---

# ship — musclead 本番リリース 自動化

引数で実行範囲を制御します(省略時は `all` = 全 phase 実行)。

| arg | 実行内容 |
|---|---|
| `feature` | Phase A のみ(develop PR 作成 + screenshot + merge) |
| `deploy` | Phase B のみ(deploy PR 作成 + merge) |
| `release` | Phase C のみ(tag + release notes) |
| `verify` | Phase D のみ(本番 backend + frontend の動作確認) |
| `all` / 省略 | A → B → C → D を順次実行 |

各 phase の最後に `TaskUpdate` で進捗を記録すること(複数 phase 実行時は phase 単位)。

---

## 前提と慣習

- branch 戦略: `feature/* → develop → main`
- main へは **squash 禁止**、 必ず `--merge`(deploy PR は merge commit で develop 履歴を残す)
- develop へは `--squash --delete-branch`
- tag は semver、 既存最新を `gh release list --limit 1` で確認してから次バージョンを決める
- 本番 URL:
  - frontend: `https://app.musclead.com`
  - backend:  `https://api.musclead.com`
- frontend は Vercel が main の push で自動 deploy、 backend は `.github/workflows/deploy-server.yml` が ECS rolling update
- 本番 DB マイグレーションは server container の `entrypoint.sh` が自動実行

---

## Phase A — feature → develop PR(スクリーンショット付き)

### 前提条件
現在のブランチが feature ブランチで、 一通り commit が積まれていること。 ユーザに「タイトルを何にするか」確認するか、 直近 commit のメッセージから自動推測する。

### 手順

1. **branch 確認**
   ```bash
   git rev-parse --abbrev-ref HEAD  # feature/* を期待
   ```

2. **lint / test / build を pass させる**
   ```bash
   cd web && npx tsc --noEmit && npx eslint --max-warnings 0 src/ && npx vitest run
   npm run build  # /settings 等の新 route が静的生成されていることも確認
   ```

3. **未 commit があれば commit + push**(ユーザに必ず確認、 push 先は当該 feature branch)

4. **スクリーンショット撮影**(UI 変更を伴う PR の場合のみ。 非 UI なら skip)
   - 既存ユーザ(or 一時ユーザを `POST /users` で作成)で login
   - `docs/screenshots/<feature-slug>/` 配下に png 保存(命名規則: `NN-title.png`、 数字 prefix で並び順を担保)
   - 撮影 viewport: desktop = 1280x800, mobile = 390x844、 両方 `deviceScaleFactor: 2`
   - dev server 起動方法: `cd web && npm run dev`(既存サーバが起動してたら使う)、 backend は `ADDR=:8081 make run` で別ポートに上げ frontend を `NEXT_PUBLIC_API_BASE_URL=http://localhost:8081` で起動するのが安全
   - 必要なら `node` script で playwright を叩く(`web/node_modules/playwright` を利用):
     ```bash
     cd web && cat > scripts/_shots.mjs <<EOF
     import { chromium } from "playwright";
     ...
     EOF
     node scripts/_shots.mjs && rm scripts/_shots.mjs
     ```
   - 撮ったら `git add docs/screenshots/<feature-slug>/` で commit + push

5. **PR 作成**
   - title: `feat(scope): description`(Conventional Commits)
   - base: `develop`
   - body には以下を含める:
     - `## Summary` — 何をやったか 3 行以内
     - `## Screenshots` — `raw.githubusercontent.com/Watari995/musclead/<branch>/docs/screenshots/<slug>/X.png` を `<img>` で参照(table で light/dark や desktop/mobile 並列)
     - `## 主要な変更` — UI / API / 内部設計 等の section
     - `## 動作確認 (manual)` — 表形式
     - `## Test plan` — checkbox リスト
     - 末尾に `🤖 Generated with [Codex](https://Codex.com/Codex)`
   ```bash
   gh pr create --base develop --title "..." --body "$(cat <<'EOF'
   ...
   EOF
   )"
   ```

6. **CI 通過待ち** — `gh pr checks <num>` で polling、 mergeStateStatus が CLEAN になるまで待つ。 Vercel preview は CI gating 外なので blocking しない
   ```bash
   until [ "$(gh pr view <N> --json mergeStateStatus --jq .mergeStateStatus)" = "CLEAN" ] || [ "$(gh pr view <N> --json mergeStateStatus --jq .mergeStateStatus)" = "BLOCKED" ]; do sleep 30; done
   ```

7. **merge** — squash + branch 削除
   ```bash
   gh pr merge <N> --squash --delete-branch
   git checkout develop && git pull --ff-only
   ```

---

## Phase B — develop → main deploy PR

1. **diff 確認**
   ```bash
   git log --oneline main..develop | head -20
   ```
   含まれる PR 番号を控える(release notes で参照)

2. **次バージョンを決める**
   ```bash
   gh release list --limit 1  # 最新 release の tag (例: v0.1.5)
   git tag --sort=-v:refname | head -5  # 既存 tag(orphan tag に注意)
   ```
   semver で minor / patch を判断:
   - 新機能(新 endpoint / 新ページ等) → minor (`v0.X.0`)
   - bug fix / 微修正 → patch (`v0.X.Y`)
   - 既存 tag と衝突する場合は次の空きを使う

3. **deploy PR 作成**
   ```bash
   gh pr create --base main --head develop --title "Deploy prod vX.Y.Z: <短い要約>" --body "$(cat <<'EOF'
   ## Summary
   ...
   ## 含まれる主要な変更
   ### Backend
   - PR #NN ...
   ### Frontend
   - PR #MM ...
   ## マイグレーション
   (有無を明記)
   ## Test plan
   ...
   EOF
   )"
   ```

4. **CI 待ち + merge**(squash ではなく `--merge` で merge commit)
   ```bash
   gh pr merge <N> --merge
   git checkout main && git pull --ff-only
   ```

---

## Phase C — tag + GitHub Release

1. **tag 作成 + push**(annotated)
   ```bash
   git tag -a vX.Y.Z -m "vX.Y.Z — <短い要約>" <merge-sha>
   git push origin vX.Y.Z
   ```

2. **GitHub Release 作成**
   - title: `vX.Y.Z — <短い要約>`
   - body 構造:
     - `## ハイライト` — 大カテゴリ別に絵文字 + 概要
     - 各カテゴリで含まれる変更点を箇条書き
     - `## マイグレーション` — 自動/手動どちらかを明記
     - `## API 互換性` — 破壊的変更があれば必ず書く
     - `## 関連 PR` — 番号 + タイトル
     - `**Full Changelog**: https://github.com/Watari995/musclead/compare/<prev-tag>...<new-tag>`
   ```bash
   gh release create vX.Y.Z --title "..." --notes "$(cat <<'EOF'
   ...
   EOF
   )"
   ```

---

## Phase D — 本番動作確認

### 待機
1. **server deploy 完了待ち**(ECS rolling)
   ```bash
   until [ "$(gh run list --branch main --workflow 'Deploy Server' --limit 1 --json status --jq '.[0].status')" = "completed" ]; do sleep 20; done
   ```
2. **Vercel deploy 完了待ち**
   - GitHub status `Vercel` が `success` になるまで polling、 ただし **pending のまま 30 分以上止まる** ことがある(過去事例 #35)。 30 分超えたら user に Vercel ダッシュボード確認を依頼すること
   - 並行して、 home page の HTML に `bg-white`(旧コード)が消えて `bg-[var(--color-surface)]` に置き換わったかで build 反映を判定するのが確実:
     ```bash
     curl -s https://app.musclead.com/ | grep -oE "bg-white|bg-\[var\(--color-surface\)\]" | sort -u
     ```

### Smoke test(backend、 curl で)
```bash
EMAIL="prod-smoke-$(date +%s)@example.com"
PASS="ProdSmoke12345"
# 1) register
curl -sf -X POST https://api.musclead.com/users -H "Content-Type: application/json" \
  -d "{\"name\":\"ProdSmoke\",\"email\":\"$EMAIL\",\"password\":\"$PASS\"}"
# 2) login
TOKEN=$(curl -s -X POST https://api.musclead.com/auth/login -H "Content-Type: application/json" \
  -d "{\"email\":\"$EMAIL\",\"password\":\"$PASS\"}" | python3 -c "import json,sys;print(json.load(sys.stdin)['access_token'])")
# 3) me
curl -s -X GET https://api.musclead.com/users/me -H "Authorization: Bearer $TOKEN" | python3 -m json.tool
# 4) (機能関連の追加 endpoint があれば curl で確認)
```

### Frontend スクリーンショット(任意、 UI 変更があった場合)
playwright で本番 URL に対して `app.musclead.com/login` → ログイン → 主要画面の screenshot を撮る。 取った画像は手元保存のみ(commit しない)。

### 報告
- Phase D 完了時に、 backend / frontend それぞれの動作確認結果を表で報告する
- 失敗箇所があれば即停止し、 user に判断を仰ぐ(自動 rollback はしない)

---

## エラー時の挙動

- どの phase でも、 想定外のエラー(test 失敗、 CI 失敗、 push 失敗、 deploy 失敗等)が出たら **即座に停止して user に報告**。 自己判断で revert / force push / tag 移動などは行わない
- 既存の git memory(`feedback_ask_before_acting.md`)を尊重し、 destructive な操作は必ず確認をとる
- 委任作業中の `feedback_auto_push_delegated.md` ルールにより、 feature branch への push は確認なしで OK

---

## 簡易チェックリスト(冒頭で確認)

skill 起動時、 user に何も問わずに以下を一度に確認:

- [ ] 現在のブランチが `feature/*` か(Phase A から始める前提)、 もしくは既に `develop` か(Phase B から)、 `main` か(Phase D のみ)
- [ ] 引数で phase 指定があるか(`feature` / `deploy` / `release` / `verify` / `all`)
- [ ] 直近 commit 内容から PR title / リリースタイプ(minor or patch)を推測できるか

判断に迷ったら user に質問。
