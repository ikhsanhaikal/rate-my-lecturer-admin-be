-- name: ListSubjects :many
SELECT * FROM subjects;

-- name: GetSubjectsByPk :one
SELECT * FROM subjects
WHERE subjects.id = ?;

