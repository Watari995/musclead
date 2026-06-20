CREATE TABLE food_products (
  id              BINARY(16)    NOT NULL,
  barcode         VARCHAR(14)   NULL,
  name            VARCHAR(100)  NOT NULL,
  calories        INT           NOT NULL,
  protein_g       DECIMAL(6,2)  NULL,
  fat_g           DECIMAL(6,2)  NULL,
  carbohydrate_g  DECIMAL(6,2)  NULL,
  register_source VARCHAR(20)   NOT NULL,
  created_at      DATETIME(6)   NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  updated_at      DATETIME(6)   NOT NULL DEFAULT CURRENT_TIMESTAMP(6)
                                ON UPDATE CURRENT_TIMESTAMP(6),
  PRIMARY KEY (id),
  KEY idx_barcode (barcode),
  KEY idx_name (name)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_0900_ai_ci;
