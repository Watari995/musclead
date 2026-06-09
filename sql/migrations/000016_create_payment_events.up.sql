-- payment 集約の監査ログ (append-only)。
-- 状態遷移 (initiated / succeeded / failed / canceled / renewed) を時系列で記録。
-- updated_at なし、 状態を書き換えない (= 追記のみ)。
-- 将来の決済 SaaS 追加に備え、 決済サービス非依存の抽象表現にする (ADR 0014)。
CREATE TABLE payment_events (
  id          BINARY(16)    NOT NULL,
  payment_id  BINARY(16)    NOT NULL,
  event_type  VARCHAR(50)   NOT NULL,                              -- 'payment.initiated' / 'payment.succeeded' / ...
  metadata    JSON          NOT NULL,                              -- 状態遷移時のスナップショット
  created_at  DATETIME(6)   NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  PRIMARY KEY (id),
  KEY idx_payment_created (payment_id, created_at),
  CONSTRAINT fk_payment_events_payment
    FOREIGN KEY (payment_id) REFERENCES payments(id)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_0900_ai_ci;
