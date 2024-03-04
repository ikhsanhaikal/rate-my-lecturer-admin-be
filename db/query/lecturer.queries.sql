-- name: CreateLecturer :execresult
INSERT INTO lecturers (name, email, description, gender, labId) 
			 VALUES (?, ?, ?, ?, ?);

-- name: DeleteLecturer :exec
DELETE FROM lecturers
WHERE id = ?;

-- name: ListLecturers :many
SELECT * FROM lecturers;


