-- name: ListSubjects :many
SELECT * FROM subjects
ORDER BY id
LIMIT ? OFFSET ?;

-- name: GetSubjectsByPk :one
SELECT * FROM subjects
WHERE subjects.id = ?;

-- name: CountSubjects :one
SELECT COUNT(*) FROM subjects;

-- name: DeleteSubjectByPk :execresult
DELETE FROM subjects WHERE subjects.id = ?;

-- name: CreateSubject :execlastid
INSERT INTO subjects (name, description, editorId)
VALUES (?, ?, ?);

-- name: UpdateSubject :exec
UPDATE subjects  
SET 
    name = COALESCE(sqlc.narg(name) , name),
    description = COALESCE(sqlc.narg(description), description)
WHERE id = sqlc.arg(id);