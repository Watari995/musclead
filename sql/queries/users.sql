-- name: CreateUser :exec
INSERT INTO users (id, name, email, password_hash, birthday, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: FindUserByID :one
SELECT id, name, email, password_hash, birthday, deleted_at, created_at, updated_at
FROM users
WHERE id = ? AND deleted_at IS NULL;

-- name: FindUserByEmail :one
SELECT id, name, email, password_hash, birthday, deleted_at, created_at, updated_at
FROM users
WHERE email = ? AND deleted_at IS NULL;

-- name: UpdateUser :exec
UPDATE users SET name = ?, email = ?, password_hash = ?, birthday = ?, deleted_at = ?, updated_at = ? WHERE id = ? AND deleted_at IS NULL;