---
name: ship-mobile
description: musclead iOS(Flutter)の TestFlight 配信 + App Store 審査提出を自動化。 バージョン採番 → pubspec 反映 → PR → mobile-v* タグ push で Codemagic を起動し、 TestFlight 配信と審査提出まで通す。
argument-hint: [version 例 1.0.1 (省略時は最新タグから patch+1) | --tag-only]
---

# ship-mobile — iOS リリース 自動化

`mobile-v*` タグの push をトリガに Codemagic(`codemagic.yaml` の `ios-testflight`)が
**ビルド → 署名 → TestFlight 配信 → App Store 審査提出(`submit_to_app_store`)** を実行する。
このスキルはその「トリガを正しく作る」までを担当する。 配信・審査提出そのものは Codemagic 側。

| arg | 実行内容 |
|---|---|
| `1.0.1` 等 | そのバージョン名で採番(pubspec の version 名を更新) |
| 省略 | 最新 `mobile-v*` タグから patch +1 を自動採番 |
| `--tag-only` | バージョン bump / PR を行わず、 現在の main HEAD にタグだけ打つ |

---

## 前提(初回のみ・人間が用意済みであること)

このスキルが成立する前提。 不足を見つけたら即停止して user に伝える。

- **Codemagic UI 側**: App Store Connect API キー連携(名前 `musclead_asc`)+ 環境変数グループ `app_store_connect` 設定済み(`codemagic.yaml` 冒頭コメント参照)
- **App Store Connect 側**: アプリレコード作成済み、 **初回リリースのメタデータ(スクショ・説明文・プライバシー・年齢レーティング・輸出コンプライアンス)が一度入力済み**
  - これが無いと `submit_to_app_store` は失敗する。 2 回目以降の提出はメタデータが流用されるので自動で通る
- branch 戦略は `ship` と同じ: `feature/* → develop → main`、 タグは **main の merge commit** に打つ
- ビルド番号は Codemagic が `$BUILD_NUMBER` で自動採番するため pubspec の `+N` は触らない(version 名のみ更新)

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
3. commit + push(feature branch への push は委任ルールで確認不要)
   ```bash
   git add mobile/pubspec.yaml
   git commit -m "chore(mobile): release vX.Y.Z"
   git push -u origin <branch>
   ```
4. develop へ PR(`--squash --delete-branch`)→ CI 通過待ち → merge
5. develop → main の deploy PR(`--merge`)→ CI 通過待ち → merge
   （`ship` Phase A/B と同じ手順。 mobile のみの変更なら deploy PR の本文は簡潔で良い）
6. `git checkout main && git pull --ff-only` で merge commit を取得

---

## Phase C — タグ push で Codemagic 起動

1. **annotated タグを main の HEAD(merge commit)に打つ**
   ```bash
   git tag -a mobile-vX.Y.Z -m "iOS vX.Y.Z — <短い要約>"
   git push origin mobile-vX.Y.Z
   ```
   → これで Codemagic の `ios-testflight` workflow が起動する。

2. user に伝える: 「タグ push 済み。 Codemagic でビルド → TestFlight 配信 → 審査提出が走る」
   - Codemagic ダッシュボードでビルド進捗を確認(このスキルからは監視しない。 API トークン連携があれば将来自動化可)

---

## Phase D — 結果の確認(手動ポイントの明示)

Codemagic 完了後の確認観点を user に案内する(スキルは状態を断定しない):

- **TestFlight**: ASC → TestFlight に新ビルド(version X.Y.Z / build = Codemagic のビルド番号)が出たか
- **審査提出**: ASC → App Store の当該バージョンが「審査待ち(Waiting for Review)」になったか
  - もし「メタデータ不備」で提出失敗していたら、 不足項目を ASC で埋めて再タグ(または ASC で手動提出)
- `release_type` は `codemagic.yaml` で `AFTER_APPROVAL`(承認後 自動公開)。 手動公開にしたい場合は `MANUAL` に変更

---

## エラー時の挙動

- analyze / test / CI / push のいずれかが失敗したら **即停止して user に報告**。 自己判断で revert / force push / タグ移動はしない
- タグ名が既存と衝突したら停止して次の空きバージョンを user に確認(打ち直しの強制移動はしない)
- destructive な操作は git memory(`feedback_ask_before_acting.md`)に従い必ず確認

---

## 簡易チェックリスト(冒頭で確認)

- [ ] 引数の version 指定有無(無ければ最新 `mobile-v*` から採番)
- [ ] 前提(Codemagic 連携 / ASC メタデータ)が整っているか。 初回提出が未実施なら user に確認
- [ ] 現在ブランチ(feature から A→B→C か、 `--tag-only` で main に C のみか)
