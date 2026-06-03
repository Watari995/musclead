CREATE TABLE routine_exercises (
    id BINARY(16) NOT NULL,
    routine_id BINARY(16) NOT NULL,
    exercise_id BINARY(16) NOT NULL,
    display_order INT NOT NULL,
    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),

    PRIMARY KEY(id),
    FOREIGN KEY(routine_id) REFERENCES routines(id) ON DELETE CASCADE,
    FOREIGN KEY(exercise_id) REFERENCES exercises(id) ON DELETE RESTRICT,

    UNIQUE KEY uk_routine_id_exercise_id (routine_id, exercise_id),
    KEY idx_routine_id_display_order (routine_id, display_order)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_0900_ai_ci;