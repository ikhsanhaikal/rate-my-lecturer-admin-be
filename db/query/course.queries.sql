-- name: ListCourseByLecturer :many
SELECT * FROM classes c 
WHERE c.lecturerId = ? 
LIMIT ? OFFSET ?;

-- name: CountCourseByLecturer :one
SELECT COUNT(*) FROM classes c 
WHERE c.lecturerId = ?;

-- name: CreateCourse :execlastid
INSERT INTO classes (lecturerId, subjectId, year, semester)
VALUES (?, ?, ?, ?);

-- name: GetCourseById :one
SELECT * FROM classes c
WHERE c.id = ?;


-- name: DeleteCourseById :execresult
DELETE FROM classes c
WHERE c.id = ?;

-- name: UpdateCourse :execlastid
UPDATE classes  
SET 
    subjectId = COALESCE(sqlc.narg(subjectId), subjectId),
    year = COALESCE(sqlc.narg(year), year),
    semester  = COALESCE(sqlc.narg(semester), semester)
WHERE id = sqlc.arg(id);