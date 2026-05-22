-- +goose Up
CREATE TABLE meals (
  id              BINARY(16)    NOT NULL,
  user_id         BINARY(16)    NOT NULL,
  eaten_at        DATETIME(6)   NOT NULL,
  meal_type       VARCHAR(20)   NOT NULL,
  calories        INT           NOT NULL,
  protein_g       DECIMAL(6, 2) NULL,
  fat_g           DECIMAL(6, 2) NULL,
  carbohydrate_g  DECIMAL(6, 2) NULL,
  memo            VARCHAR(1000) NULL,
  created_at      DATETIME(6)   NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  updated_at      DATETIME(6)   NOT NULL DEFAULT CURRENT_TIMESTAMP(6)
                                         ON UPDATE CURRENT_TIMESTAMP(6),
  PRIMARY KEY (id),
  KEY idx_user_eaten (user_id, eaten_at DESC),
  CONSTRAINT fk_meals_user_id
    FOREIGN KEY (user_id) REFERENCES users(id)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_0900_ai_ci;

-- +goose Down
DROP TABLE meals;
