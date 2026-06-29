ALTER TABLE meals
  DROP FOREIGN KEY fk_meals_food_product_id,
  DROP COLUMN food_product_id,
  DROP COLUMN serving_count;
