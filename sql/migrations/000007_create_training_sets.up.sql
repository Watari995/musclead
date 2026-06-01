CREATE TABLE training_sets (
    id BINARY(16) NOT NULL,
    training_exercise_id BINARY(16) NOT NULL,
    set_number INT NOT NULL,
    -- 全体の桁数は6桁までで、小数2桁まで入る
    weight_kg DECIMAL(6,2) NOT NULL,
    reps INT NOT NULL,
    rest_seconds INT NULL, -- nullの時はtraining_exercises.rest_secondsを使用する
    memo TEXT NULL,
    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),

    PRIMARY KEY(id),
    FOREIGN KEY(training_exercise_id) REFERENCES training_exercises(id) ON DELETE CASCADE,

    -- 一種目でセット数はユニーク
    UNIQUE KEY uq_training_exercise_id_set_number (training_exercise_id, set_number)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_0900_ai_ci;