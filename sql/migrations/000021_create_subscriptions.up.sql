-- purchase 集約: 継続的な権利状態。
-- ADR 0013 で確定の「Pro である権利」 集約。 Pro 機能 gate (ProGate middleware) は本テーブルを見る。
--
-- 状態遷移 (ADR 0017):
--   active   → 課金成功時に INSERT (Webhook 'checkout.session.completed')
--   canceled → 解約予約時 (Customer Portal で「期末解約」)、 ただし期末まで Pro 機能利用可
--   expired  → 期末経過後の最終状態 (Webhook 'customer.subscription.deleted')
--
-- Pro 判定: status に関係なく expires_at > NOW() なら Pro 扱い (ADR 0017)。
--
-- 集約間参照:
--   - subscription_order_id: 起源の order (履歴追跡用、 admin 手動作成では NULL)
--   - payment_id: 紐付く payment 集約への参照 (DDD: ID のみ)
CREATE TABLE subscriptions (
  id                    BINARY(16)   NOT NULL,
  user_id               BINARY(16)   NOT NULL,
  plan                  VARCHAR(20)  NOT NULL,                        -- 'pro'
  status                VARCHAR(20)  NOT NULL,                        -- 'active' / 'canceled' / 'expired'
  subscription_order_id BINARY(16)   NULL,                            -- 起源の申込 (admin 手動作成は NULL)
  payment_id            BINARY(16)   NOT NULL,                        -- payment 集約への参照
  activated_at          DATETIME(6)  NOT NULL,                        -- Pro 開始時刻
  expires_at            DATETIME(6)  NOT NULL,                        -- Pro 期限 (outbox 経由で payment.current_period_end を反映)
  canceled_at           DATETIME(6)  NULL,                            -- 解約予約時刻
  created_at            DATETIME(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  updated_at            DATETIME(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  PRIMARY KEY (id),
  KEY idx_user_status (user_id, status),                              -- Pro gate / 解約予約状態の検索用
  KEY idx_payment (payment_id),
  KEY idx_order (subscription_order_id),
  CONSTRAINT fk_subscriptions_user FOREIGN KEY (user_id) REFERENCES users(id)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_0900_ai_ci;
