CREATE TABLE trainings (
    id BINARY(16) NOT NULL,
    user_id BINARY(16) NOT NULL,
    started_at DATETIME(6) NOT NULL,
    ended_at DATETIME(6) NULL DEFAULT NULL,
    memo TEXT NULL DEFAULT NULL,
    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),

    PRIMARY KEY(id),
    FOREIGN KEY(user_id) REFERENCES users(id)
    ON DELETE CASCADE,

    -- 一覧クエリ用のIndexを貼る
    KEY idx_user_id_started_at (user_id, started_at DESC)
 ) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_0900_ai_ci;