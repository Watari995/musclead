-- payment 集約の本体テーブル。
-- 1 record = 1 Stripe Subscription (ADR 0017)。
-- stripe_subscription_id に UNIQUE 制約を貼って、 ユーザー再加入時のサブスク重複を防ぐ。
CREATE TABLE payments (
  id                         BINARY(16)    NOT NULL,
  user_id                    BINARY(16)    NOT NULL,
  currency                   VARCHAR(3)    NOT NULL,                      -- 'JPY'
  status                     VARCHAR(20)   NOT NULL,                      -- 'pending' / 'succeeded' / 'failed' / 'canceled'
  stripe_customer_id         VARCHAR(255)  NULL,
  stripe_subscription_id     VARCHAR(255)  NULL,
  stripe_checkout_session_id VARCHAR(255)  NULL,
  checkout_url               VARCHAR(500)  NULL,                          -- Stripe Checkout のリダイレクト先
  current_period_end         DATETIME(6)   NULL,                          -- Stripe subscription の現在の課金期間終了
  succeeded_at               DATETIME(6)   NULL,
  failed_at                  DATETIME(6)   NULL,
  failure_reason             VARCHAR(255)  NULL,
  created_at                 DATETIME(6)   NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  updated_at                 DATETIME(6)   NOT NULL DEFAULT CURRENT_TIMESTAMP(6)
                                                    ON UPDATE CURRENT_TIMESTAMP(6),
  PRIMARY KEY (id),
  UNIQUE KEY uk_stripe_subscription (stripe_subscription_id),
  KEY idx_user (user_id),
  CONSTRAINT fk_payments_user
    FOREIGN KEY (user_id) REFERENCES users(id)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_0900_ai_ci;
