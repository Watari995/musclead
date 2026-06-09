// Package dto は payment module の HTTP 入出力 DTO を定義する。
package dto

// TODO: HTTP リクエスト / レスポンスの DTO を定義
//   - 例 1: WebhookEvent (Stripe Webhook の生 payload を扱う用)
//   - 例 2: PaymentDTO (GET 系で payment の状態を返す用)
//
// 参考: internal/user/dto/user_dto.go
//
// Tips:
//   - DTO は struct + json タグだけ持つシンプルな型
//   - domain entity の private field にアクセスするための getter を呼んで構築する
//   - Map関数 (e.g. ToPaymentDTO(p *domain.Payment) PaymentDTO) を同ファイルに置く
