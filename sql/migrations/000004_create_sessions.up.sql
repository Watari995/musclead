CREATE TABLE sessions (
    id BINARY(16) NOT NULL,
    user_id BINARY(16) NOT NULL,
    refresh_hash VARCHAR(255) NOT NULL,
    user_agent VARCHAR(255) NULL DEFAULT NULL,
    ip_address VARCHAR(45) NULL DEFAULT NULL,
    expires_at DATETIME(6) NOT NULL,
    revoked_at DATETIME(6) NULL DEFAULT NULL,
    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),

    PRIMARY KEY(id),
    FOREIGN KEY(user_id) REFERENCES users(id)
    ON DELETE CASCADE,
    UNIQUE KEY uq_sessions_refresh_hash(refresh_hash)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_0900_ai_ci;