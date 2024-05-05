-- name: GetLabsByPk :one
SELECT * FROM labs
WHERE id = ?;

-- name: ListLabs :many
SELECT * FROM labs
ORDER BY id LIMIT ? OFFSET ?;

-- name: ListMembers :many
SELECT l1.* FROM lecturers l1 
JOIN labs l2  ON l1.labId = l2.id
WHERE l2.id = ?;

-- name: CountLabs :one
SELECT COUNT(*) FROM labs;

-- name: CreateLab :execresult
INSERT INTO labs (name, code, description) 
VALUES (?, ?, ?);

-- name: DeleteLab :execresult
DELETE FROM labs WHERE labs.id = ?;  
