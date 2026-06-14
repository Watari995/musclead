-- メール通知用の outbox (ADR 0015)。
-- Webhook 処理の TX 内で INSERT し、 TX 外で SNS publish する。
-- 即時 publish 失敗時は 1 分後の outbox-relay Lambda が拾う (failsafe)。
CREATE TABLE outbox_events (
  id            BINARY(16)    NOT NULL,
  event_type    VARCHAR(100)  NOT NULL,                            -- 'PaymentSucceeded' / 'PaymentFailed' / ...
  aggregate_id  BINARY(16)    NOT NULL,                            -- payment_id
  payload       JSON          NOT NULL,                            -- SNS / SQS に流す本文
  published_at  DATETIME(6)   NULL,                                -- NULL = 未配信
  publish_error TEXT          NULL,
  created_at    DATETIME(6)   NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  updated_at    DATETIME(6)   NOT NULL DEFAULT CURRENT_TIMESTAMP(6)
                                      ON UPDATE CURRENT_TIMESTAMP(6),
  PRIMARY KEY (id),
  KEY idx_unpublished (published_at, created_at),
  KEY idx_aggregate (aggregate_id)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_0900_ai_ci;
