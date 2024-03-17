-- name: CreateLecturer :execlastid
INSERT INTO lecturers (name, email, description, gender, labId) 
			 VALUES (?, ?, ?, ?, ?);

-- name: DeleteLecturersByPk :exec
DELETE FROM lecturers
WHERE id = ?;

/*
null restriction nggak bisa bro
INSERT INTO lecturers(id, name, email, description, gender, labId)   
VALUES (?, ?, ?, ?, ?, ?)  
ON DUPLICATE KEY UPDATE 
		name = VALUES(name),
		email = VALUES(email),  
		description = VALUES(description),
		gender = VALUES(gender),
		labId = VALUES(labId);*/

-- name: UpdateLecturer :exec
UPDATE lecturers  
SET 
    name = COALESCE(sqlc.narg(name) , name),
    email = COALESCE(sqlc.narg(email), email),
    description = COALESCE(sqlc.narg(description), description),
    gender  = COALESCE(sqlc.narg(gender), gender),
    labId  = COALESCE(sqlc.narg(labId), labId)
WHERE id = sqlc.arg(id);

-- name: ListLecturers :many
SELECT *  FROM lecturers;

-- name: GetLecturersByPk :one
SELECT * FROM lecturers 
WHERE lecturers.id = ?;


