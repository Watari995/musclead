
ALTER TABLE training_exercises
DROP COLUMN name,
ADD COLUMN exercise_id BINARY(16) NOT NULL AFTER training_id,
ADD CONSTRAINT fk_training_exercises_exercises FOREIGN KEY (exercise_id) REFERENCES exercises (id) ON DELETE RESTRICT; -- 使用中のExerciseは削除不可にする