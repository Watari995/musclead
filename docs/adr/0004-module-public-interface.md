# ADR 0004: モジュール間連携は公開 Command Interface 経由にする

## ステータス
採用 (2026-05-31)

## コンテキスト

[ADR 0002](./0002-ddd-modular-monolith.md) で Modular Monolith strict 構成を採用。
Go の `internal/` 制約により、 モジュール間で内部実装を直接 import できない。

例: auth モジュールがログイン処理で user の情報(email/password)が必要だが、
`internal/user/internal/domain/user.UserRepository` には auth から触れない。

## 決定

**「呼び出される側のモジュールが、 UseCase レベルの interface を `interface/publicfunctions/` 配下に公開する」** 形に統一する。

### 構造

```
internal/<module>/
├── interface/
│   └── publicfunctions/
│       └── command.go          ← <Module>Command interface + Request/Response + Error 型
├── internal/
│   ├── usecase/                ← interface 実装はここ
│   ├── domain/ ...
│   └── infra/ ...
└── <module>.go                 ← Module facade (NewModule + <Module>Command() ゲッター)
```

### 命名規約

| 種類 | 名前例 |
|---|---|
| Interface | `UserCommand`(モジュール名 + Command) |
| メソッド | `Authenticate(ctx, req) (*AuthenticateResponse, error)` (動詞) |
| Request 型 | `AuthenticateRequest`(prefix モジュール名は不要) |
| Response 型 | `AuthenticateResponse` |
| Error 型 | `ErrInvalidCredentials` 等のセンチネルエラー or 独自 Error 構造体 |

### 呼び出し側のパターン

呼び出し側(例: auth)は **自前で interface を定義**(Dependency Inversion、 mock 用)。
偶然 signature が一致するので、 main.go で `userMod.UserCommand()` をそのまま渡せる。

```go
// auth/internal/usecase/login.go
type userAuthenticator interface {  // ← 小文字、 パッケージ内のみ
    Authenticate(ctx context.Context, req userpublic.AuthenticateRequest) (*userpublic.AuthenticateResponse, error)
}
```

## 理由

1. **Repository 直公開を避ける**
   - 低レベル API を露出すると内部実装が漏れる(Entity / hash 等)
   - UseCase レベルで公開 = ビジネス操作を外に出す
2. **internal の壁を保ったまま連携可能**
   - `interface/publicfunctions/` は internal 外なので import 可
   - 中の `usecase`/`domain` は引き続き保護される
3. **エラーハンドリングまで含めた契約**
   - `errors.As` で型アサーションして呼び出し側が分岐できる
   - 内部のエラー(SQL エラー等)は公開エラーに詰め替えてから返す
4. **SODA(スニーカーダンク)の実装パターンと一致**
   - 実例: `ProductCatalogCommand.CreateProduct(...)` — モジュール名 + Command interface
   - 入社後すぐ馴染める命名/構造

## 結果

- ✅ モジュール境界を破らず連携できる
- ✅ 公開する処理単位が明確(UseCase/Command 単位、 Repository ではない)
- ✅ テスト時は呼び出し側が自前 interface をモック
- ❌ ボイラープレート増(interface + DTO + Module facade ゲッター + 実装)
- ❌ 「同じことを2回書く」 ように見える(public interface / internal usecase の2箇所)

## 適用例(auth が user を呼ぶケース)

1. user 側: `internal/user/interface/publicfunctions/command.go` に `UserCommand.Authenticate` 定義
2. user 側: `internal/user/internal/usecase/authenticate.go` に実装
3. user 側: `user.go` の Module に `UserCommand()` ゲッター追加
4. auth 側: `login.go` で自前 `userAuthenticator` interface 定義 + 呼び出し
5. main.go: `auth.NewModule(userMod.UserCommand(), ...)` で wiring
