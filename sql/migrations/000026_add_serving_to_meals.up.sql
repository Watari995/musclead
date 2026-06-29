ALTER TABLE meals
  ADD COLUMN food_product_id BINARY(16)    NULL         AFTER memo,
  ADD COLUMN serving_count   DECIMAL(5, 2) NOT NULL
                                            DEFAULT 1.00 AFTER food_product_id,
  ADD CONSTRAINT fk_meals_food_product_id
    FOREIGN KEY (food_product_id) REFERENCES food_products(id);
