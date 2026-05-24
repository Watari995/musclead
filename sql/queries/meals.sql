-- name: FindAllMealsByUserIDWithOffsetPagination :many
SELECT id, user_id, eaten_at, meal_type, calories, protein_g, fat_g, carbohydrate_g, memo, created_at, updated_at
FROM meals
WHERE user_id = ? ORDER BY eaten_at DESC LIMIT ? OFFSET ?;

-- name: FindMealByID :one
SELECT id, user_id, eaten_at, meal_type, calories, protein_g, fat_g, carbohydrate_g, memo, created_at, updated_at
FROM meals
WHERE id = ?;

-- name: FindMealPhotosByMealID :many
SELECT id, meal_id, image_path, display_order, created_at
FROM meal_photos
WHERE meal_id = ? ORDER BY display_order ASC;

-- name: UpsertMeal :exec
INSERT INTO meals (id, user_id, eaten_at, meal_type, calories, protein_g, fat_g, carbohydrate_g, memo, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
eaten_at = VALUES(eaten_at),
meal_type = VALUES(meal_type),
calories = VALUES(calories),
protein_g = VALUES(protein_g),
fat_g = VALUES(fat_g),
carbohydrate_g = VALUES(carbohydrate_g),
memo = VALUES(memo),
updated_at = VALUES(updated_at);

-- name: DeleteMealByID :exec
DELETE FROM meals
WHERE id = ?;

-- name: CreateMealPhoto :exec
INSERT INTO meal_photos (id, meal_id, image_path, display_order, created_at)
VALUES (?, ?, ?, ?, ?);

-- name: DeleteMealPhotosByMealID :exec
DELETE FROM meal_photos
WHERE meal_id = ?;

-- name: CountMealsByUserID :one
SELECT COUNT(*) FROM meals WHERE user_id = ?;