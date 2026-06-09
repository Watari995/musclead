// Package payment is the public facade of the payment module.
// Modular Monolith (strict) のため、 外部からは Module 経由でのみアクセス可能。
//
// 設計: ADR 0013 (purchase / payment 分離) + ADR 0014 (Webhook 同期処理)
// 依存方向: purchase → payment (本 module は purchase を知らない)
package payment

// TODO: NewModule の実装
//   - 必要な引数: dbmap, stripeClient (SDK), retryStrategy などを受け取る
//   - 内部で repository / usecase を組み立てる
//   - Module struct を返す
//
// 参考: internal/weight/weight.go, internal/user/user.go
//   - dbmap.AddTableWithName でモデルを gorp に登録
//   - 各 repository を NewXxxRepository で生成
//   - 各 usecase を Newxxx で組み立て
//   - Module struct に Handler / Command を持たせる
//
// Module struct の構成 (案):
//   type Module struct {
//       Handler  http.Handler              // POST /payment/webhook を返す
//       command  publicfunctions.PaymentCommand
//       query    publicfunctions.PaymentQuery
//   }
//
//   func (m *Module) Command() publicfunctions.PaymentCommand { return m.command }
//   func (m *Module) Query() publicfunctions.PaymentQuery     { return m.query }
