CREATE TABLE notifications (
  id                BINARY(16)   NOT NULL,
  user_id           BINARY(16)   NOT NULL,
  notification_type VARCHAR(50)  NOT NULL,
  metadata          JSON         NOT NULL,
  read_at           DATETIME(6)  NULL     DEFAULT NULL,
  created_at        DATETIME(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  PRIMARY KEY (id),
  KEY idx_notifications_user_id (user_id),
  CONSTRAINT fk_notifications_user_id
    FOREIGN KEY (user_id) REFERENCES users (id)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_0900_ai_ci;
