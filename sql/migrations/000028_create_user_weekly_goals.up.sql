CREATE TABLE user_weekly_goals (
  id               BINARY(16)    NOT NULL,
  user_id          BINARY(16)    NOT NULL,
  training_count   INT           NULL,
  calorie_average  INT           NULL,
  weight_change_kg DECIMAL(4, 1) NULL,
  created_at       DATETIME(6)   NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  updated_at       DATETIME(6)   NOT NULL DEFAULT CURRENT_TIMESTAMP(6)
                                          ON UPDATE CURRENT_TIMESTAMP(6),
  PRIMARY KEY (id),
  UNIQUE KEY uk_user_weekly_goals_user_id (user_id),
  CONSTRAINT fk_user_weekly_goals_user_id
    FOREIGN KEY (user_id) REFERENCES users (id)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_0900_ai_ci;
