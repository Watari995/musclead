CREATE TABLE IF NOT EXISTS device_tokens (
  id         BINARY(16)   NOT NULL,
  user_id    BINARY(16)   NOT NULL,
  token      VARCHAR(255) NOT NULL,
  created_at DATETIME(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  updated_at DATETIME(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6)
                                   ON UPDATE CURRENT_TIMESTAMP(6),
  PRIMARY KEY (id),
  UNIQUE KEY uk_token (token),
  -- user_idと降順(一番最新のもの)を取ると思ったので複合インデックス
  KEY idx_user_id_created_at (user_id, created_at DESC),
  CONSTRAINT fk_device_tokens_user_id
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_0900_ai_ci;
