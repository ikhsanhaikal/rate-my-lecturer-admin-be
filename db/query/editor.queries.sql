-- name: CreateEditor :execlastid
INSERT INTO editors (username, email) VALUES (?, ?);

-- name: GetEditor :one
SELECT * FROM editors
WHERE id = ?;

-- name: GetEditorByEmail :one
SELECT * FROM editors
WHERE email = ?;