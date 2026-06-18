CREATE TABLE meal_templates (
  id             BINARY(16)    NOT NULL,
  user_id        BINARY(16)    NOT NULL,
  name           VARCHAR(100)  NOT NULL,
  display_order  INT           NOT NULL DEFAULT 0,
  meal_type      VARCHAR(20)   NOT NULL,
  calories       INT           NOT NULL,
  protein_g      DECIMAL(6, 2) NULL,
  fat_g          DECIMAL(6, 2) NULL,
  carbohydrate_g DECIMAL(6, 2) NULL,
  created_at     DATETIME(6)   NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  updated_at     DATETIME(6)   NOT NULL DEFAULT CURRENT_TIMESTAMP(6)
                                        ON UPDATE CURRENT_TIMESTAMP(6),
  PRIMARY KEY (id),
  KEY idx_meal_templates_user (user_id, display_order ASC, created_at ASC)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_0900_ai_ci;
