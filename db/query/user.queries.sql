-- name: CreateUser :execresult
INSERT INTO users (name, email) VALUES (?, ?);

-- name: GetUser :one
SELECT * FROM users
WHERE id = ?;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = ?;

-- name: ListUsers :many
SELECT * FROM users;