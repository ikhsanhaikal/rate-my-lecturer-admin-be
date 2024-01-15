-- name: GetTodo :one
SELECT * FROM todos
WHERE id = ?;

-- name: ListTodos :many
SELECT * FROM todos
ORDER BY created_at;

-- name: CreateTodo :execresult
INSERT INTO todos (task, user_id) VALUES (?, ?);

-- name: DeleteTodo :exec
DELETE FROM todos
WHERE id = ?;

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