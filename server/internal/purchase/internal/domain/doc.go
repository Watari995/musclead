// Package purchasedomain は purchase 集約の domain 層を定義する。
//
// 含むもの:
//   - entity: SubscriptionOrder, Subscription
//   - VO: SubscriptionOrderStatus, SubscriptionStatus (内部で使う、 enum 風)
//   - interface: SubscriptionOrderRepository, SubscriptionRepository
//
// 設計参考:
//   - ADR 0013 (purchase / payment 分離)
//   - ADR 0017 (Pro 判定は status に関係なく expires_at > NOW())
//   - 既存 internal/payment/internal/domain/ の entity / repository パターン
package purchasedomain
