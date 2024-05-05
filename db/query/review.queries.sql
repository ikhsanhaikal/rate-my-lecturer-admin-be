-- name: AverageRatingByLecturerId :one
SELECT AVG(rating) AS AverageRating from reviews r 
JOIN classes c ON r.classId = c.id  
WHERE lecturerId = ?;

-- name: GetCourseByReview :one
SELECT c.* FROM reviews r
JOIN classes c on r.classId = c.id
WHERE r.id = ?;

-- name: GetReviewById :one
SELECT * FROM reviews r
WHERE r.id = ?;

-- name: DeleteReviewById :execlastid
DELETE FROM reviews r
WHERE r.id = ?;

-- name: GetReviewsByLecturer :many
SELECT r.* FROM reviews r
JOIN classes c on c.id  = r.classId 
WHERE c.lecturerId = ?
LIMIT ? OFFSET ?;

-- name: CountReviewsByLecturer :one
SELECT COUNT(*) FROM reviews r
JOIN classes c on c.id  = r.classId 
WHERE c.lecturerId = ?;
