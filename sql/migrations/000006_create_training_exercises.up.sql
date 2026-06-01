CREATE TABLE training_exercises (
    id BINARY(16) NOT NULL,
    training_id BINARY(16) NOT NULL,
    name VARCHAR(50) NOT NULL,
    display_order INT NOT NULL,
    rest_seconds INT NULL,
    memo TEXT NULL,
    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),

    PRIMARY KEY(id),
    FOREIGN KEY(training_id) REFERENCES trainings(id) ON DELETE CASCADE,

    KEY idx_training_id_display_order (training_id, display_order)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_0900_ai_ci;