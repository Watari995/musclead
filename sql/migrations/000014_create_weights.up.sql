CREATE TABLE weights (
  id                   BINARY(16)    NOT NULL,
  user_id              BINARY(16)    NOT NULL,
  weight_kg            DECIMAL(5, 2) NOT NULL,
  body_fat_percentage  DECIMAL(4, 2) NULL,
  skeletal_muscle_kg   DECIMAL(5, 2) NULL,
  measured_at          DATETIME(6)   NOT NULL,
  created_at           DATETIME(6)   NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  updated_at           DATETIME(6)   NOT NULL DEFAULT CURRENT_TIMESTAMP(6)
                                              ON UPDATE CURRENT_TIMESTAMP(6),
  PRIMARY KEY (id),
  KEY idx_user_measured (user_id, measured_at DESC),
  CONSTRAINT fk_weights_user_id
    FOREIGN KEY (user_id) REFERENCES users(id)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_0900_ai_ci;
