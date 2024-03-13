-- name: GetLab :one
SELECT * FROM labs
WHERE id = ?;

-- name: ListLabs :many
SELECT * FROM labs;

-- name: ListMembers :many
SELECT * FROM lecturers l1 
JOIN labs l2  ON l1.labId = l2.id
WHERE l2.id = ?;
