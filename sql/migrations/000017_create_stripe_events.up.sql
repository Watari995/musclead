-- Stripe Webhook の受信記録 + 冪等性キー (ADR 0014)。
-- 受信した全イベントを raw payload で残し、 デバッグ / 監査 / 再処理に利用。
-- stripe_event_id (Stripe 側の evt_xxx) を UNIQUE 制約にすることで、 同一イベントの 2 回処理を物理的に防ぐ。
CREATE TABLE stripe_events (
  id                BINARY(16)    NOT NULL,
  stripe_event_id   VARCHAR(255)  NOT NULL,                            -- Stripe 側の evt_xxx
  event_type        VARCHAR(100)  NOT NULL,                            -- 'checkout.session.completed' 等
  payload           JSON          NOT NULL,                            -- Stripe event 全文
  processed_at      DATETIME(6)   NULL,                                -- NULL = 未処理 / 処理失敗
  processing_error  TEXT          NULL,                                -- エラー内容 (失敗時)
  created_at        DATETIME(6)   NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  updated_at        DATETIME(6)   NOT NULL DEFAULT CURRENT_TIMESTAMP(6)
                                          ON UPDATE CURRENT_TIMESTAMP(6),
  PRIMARY KEY (id),
  UNIQUE KEY uk_stripe_event_id (stripe_event_id),
  KEY idx_processed (processed_at),
  KEY idx_created (created_at)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_0900_ai_ci;
