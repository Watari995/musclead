CREATE TABLE healthplanet_tokens (
  id              BINARY(16)    NOT NULL,
  user_id         BINARY(16)    NOT NULL,
  access_token    VARCHAR(512)  NOT NULL,
  refresh_token   VARCHAR(512)  NOT NULL,
  expires_at      DATETIME(6)   NOT NULL,
  last_synced_at  DATETIME(6)   NULL,
  created_at      DATETIME(6)   NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  updated_at      DATETIME(6)   NOT NULL DEFAULT CURRENT_TIMESTAMP(6)
                                         ON UPDATE CURRENT_TIMESTAMP(6),
  PRIMARY KEY (id),
  UNIQUE KEY uk_user_id (user_id),
  CONSTRAINT fk_healthplanet_tokens_user_id
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_0900_ai_ci;
