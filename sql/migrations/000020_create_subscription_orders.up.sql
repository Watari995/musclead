-- purchase 集約: 申込トリガー (1 回限り、 履歴)。
-- ADR 0013 で確定の「申込意思」 集約。 1 user が複数回申込めば複数 record。
--
-- 状態遷移:
--   pending  → 申込開始時
--   succeeded → Webhook で決済成功時
--   failed    → 決済失敗時 (将来)
--
-- 集約間参照: payment_id は payment 集約への参照 (DDD 流儀、 ID のみ)。
CREATE TABLE subscription_orders (
  id           BINARY(16)   NOT NULL,
  user_id      BINARY(16)   NOT NULL,
  plan         VARCHAR(20)  NOT NULL,                        -- 'pro' (将来 'pro_annual' 等)
  status       VARCHAR(20)  NOT NULL,                        -- 'pending' / 'succeeded' / 'failed'
  payment_id   BINARY(16)   NULL,
  succeeded_at DATETIME(6)  NULL,
  failed_at    DATETIME(6)  NULL,
  created_at   DATETIME(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  updated_at   DATETIME(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  PRIMARY KEY (id),
  KEY idx_user_created (user_id, created_at),
  KEY idx_payment (payment_id),
  CONSTRAINT fk_subscription_orders_user FOREIGN KEY (user_id) REFERENCES users(id)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_0900_ai_ci;
