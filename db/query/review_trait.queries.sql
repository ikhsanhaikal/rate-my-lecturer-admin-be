-- name: GetTagsByReview :many
SELECT tr.* FROM tags t 
JOIN reviews r  ON r.id  = t.reviewId
JOIN traits tr ON tr.id  = t.traitId 
WHERE r.id = ?;

-- name: GetSummaryTagsByLecturer :many
SELECT DISTINCT tr.* FROM tags t 
JOIN reviews r  ON r.id  = t.reviewId
JOIN traits tr ON tr.id  = t.traitId 
JOIN classes c on c.id = r.classId 
WHERE c.lecturerId = ?;


