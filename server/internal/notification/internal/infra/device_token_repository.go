package notificationinfra

// deviceTokenRepository: notificationdomain.DeviceTokenRepository / DeviceTokenQuery の実装。
//
// - Save: token の UNIQUE 制約を利用した ON DUPLICATE KEY UPDATE upsert
//   (shared/infra/outbox の outboxEventRepository.Save と同じパターン)
// - FindAllByUserID: 1 ユーザーが複数端末を持つ前提で []DeviceTokenView を返す
