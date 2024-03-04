-- name: GetLab :one
SELECT * FROM labs
WHERE id = ?;