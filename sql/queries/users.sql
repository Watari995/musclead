-- name: UpsertUser :exec
INSERT INTO users (id, name, email, password_hash, birthday, deleted_at, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
name = VALUES(name),
email = VALUES(email),
password_hash = VALUES(password_hash),
birthday = VALUES(birthday),
deleted_at = VALUES(deleted_at),
updated_at = VALUES(updated_at);

-- name: FindUserByID :one
SELECT id, name, email, password_hash, birthday, deleted_at, created_at, updated_at
FROM users
WHERE id = ? AND deleted_at IS NULL;

-- name: FindUserByEmail :one
SELECT id, name, email, password_hash, birthday, deleted_at, created_at, updated_at
FROM users
WHERE email = ? AND deleted_at IS NULL;