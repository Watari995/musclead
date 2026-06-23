---
name: ship-mobile
description: musclead iOS(Flutter)の TestFlight 配信 + App Store 審査提出を自動化。 バージョン採番 → pubspec 反映 → PR → GitHub Actions (testflight.yml) 手動実行 までを担当。
argument-hint: [version 例 1.0.1 (省略時は最新タグから patch+1) | --tag-only]
---

# ship-mobile — iOS リリース 自動化

`.github/workflows/testflight.yml` を `workflow_dispatch` で手動実行することで
**ビルド → 署名 → TestFlight 配信 → App Store 審査提出** を実行する。
このスキルはその「バージョン採番 → pubspec 更新 → main へのマージ」までを担当する。

| arg | 実行内容 |
|---|---|
| `1.0.1` 等 | そのバージョン名で採番(pubspec の version 名を更新) |
| 省略 | 最新 `mobile-v*` タグから patch +1 を自動採番 |
| `--tag-only` | バージョン bump / PR を行わず、 現在の main HEAD にタグだけ打つ |

---

## 前提(初回のみ・人間が用意済みであること)

このスキルが成立する前提。 不足を見つけたら即停止して user に伝える。

- **GitHub Secrets**: `IOS_DIST_CERT_BASE64` / `IOS_DIST_CERT_PASSWORD` / `IOS_PROFILE_BASE64` / `ASC_KEY_ID` / `ASC_KEY_BASE64` / `ASC_ISSUER_ID` / `SENTRY_DSN` が設定済み
- **App Store Connect 側**: アプリレコード作成済み、 **初回リリースのメタデータ(スクショ・説明文・プライバシー・年齢レーティング・輸出コンプライアンス)が一度入力済み**
  - これが無いと審査提出は失敗する。 2 回目以降の提出はメタデータが流用される
- branch 戦略は `ship` と同じ: `feature/* → develop → main`

---

## Phase A — 採番 & pre-flight

1. **次バージョンを決める**
   ```bash
   git tag --sort=-v:refname | grep '^mobile-v' | head -3   # 既存の iOS タグ
   grep '^version:' mobile/pubspec.yaml                       # 現在の version 名+ビルド
   ```
   - 引数があればそれを採用。 無ければ最新 `mobile-v*` から patch +1
   - 新機能リリースなら minor、 bug fix なら patch を user に確認(迷ったら質問)

2. **品質ゲート**(必ず pass させる。 失敗したら停止)
   ```bash
   cd mobile && fvm flutter analyze && fvm flutter test && cd ..
   ```

`--tag-only` の場合は Phase A の採番のみ行い B を skip して C へ。

---

## Phase B — version bump を main へ反映(PR 経由)

> **main / develop へ直接 push しない**(git memory `feedback_no_direct_push_to_protected.md`)。 必ず PR 経由。

1. feature ブランチに居なければ作る(例 `chore/mobile-release-vX.Y.Z`)
2. `mobile/pubspec.yaml` の `version:` 行を `X.Y.Z+<現状のビルド番号>` に更新(version 名のみ変更)
   - ビルド番号は GitHub Actions の `${{ github.run_number }}` で自動採番されるため `+N` は触らない
3. commit + push(feature branch への push は委任ルールで確認不要)
   ```bash
   git add mobile/pubspec.yaml
   git commit -m "chore(mobile): release vX.Y.Z"
   git push -u origin <branch>
   ```
4. develop へ PR(`--squash --delete-branch`)→ CI 通過待ち → merge
5. develop → main の deploy PR(`--merge`)→ CI 通過待ち → merge
6. `git checkout main && git pull --ff-only` で merge commit を取得

---

## Phase C — GitHub Actions で TestFlight 配信

1. **`testflight.yml` を workflow_dispatch で実行**
   ```bash
   gh workflow run testflight.yml --field version=X.Y.Z
   ```

2. user に伝える: 「GitHub Actions の testflight.yml を起動しました。ビルド → TestFlight 配信 → 審査提出が走ります」
   - Actions タブでビルド進捗を確認

---

## Phase D — 結果の確認(手動ポイントの明示)

GitHub Actions 完了後の確認観点を user に案内する(スキルは状態を断定しない):

- **TestFlight**: ASC → TestFlight に新ビルド(version X.Y.Z / build = run_number)が出たか
- **審査提出**: ASC → App Store の当該バージョンが「審査待ち(Waiting for Review)」になったか
  - もし「メタデータ不備」で提出失敗していたら、 不足項目を ASC で埋めて再実行

---

## エラー時の挙動

- analyze / test / CI / push のいずれかが失敗したら **即停止して user に報告**。 自己判断で revert / force push はしない
- destructive な操作は git memory(`feedback_ask_before_acting.md`)に従い必ず確認

---

## 簡易チェックリスト(冒頭で確認)

- [ ] 引数の version 指定有無(無ければ最新 `mobile-v*` から採番)
- [ ] 前提(GitHub Secrets / ASC メタデータ)が整っているか。 初回提出が未実施なら user に確認
- [ ] 現在ブランチ(feature から A→B→C か、 `--tag-only` で main に C のみか)
