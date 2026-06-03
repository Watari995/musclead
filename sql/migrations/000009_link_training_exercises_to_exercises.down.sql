ALTER TABLE training_exercises
DROP FOREIGN KEY fk_training_exercises_exercises,
DROP COLUMN exercise_id,
ADD COLUMN name VARCHAR(50) NOT NULL AFTER training_id;