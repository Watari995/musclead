-- +goose Up
CREATE TABLE meal_photos (
  id             BINARY(16)   NOT NULL,
  meal_id        BINARY(16)   NOT NULL,
  image_path     VARCHAR(255) NOT NULL,
  display_order  INT          NOT NULL DEFAULT 0,
  created_at     DATETIME(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  PRIMARY KEY (id),
  KEY idx_meal (meal_id),
  CONSTRAINT fk_meal_photos_meal_id
    FOREIGN KEY (meal_id) REFERENCES meals(id) ON DELETE CASCADE
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_0900_ai_ci;

-- +goose Down
DROP TABLE meal_photos;
