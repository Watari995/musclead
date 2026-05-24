-- name: CreateUser :exec
INSERT INTO users (id, name, email, password_hash, birthday, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: 