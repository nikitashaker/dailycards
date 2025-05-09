-- name: CreateLog :one
INSERT INTO logs (user_id, pack_id, rating_improved, rating_worsen, cards_learned, cards_mastered)
VALUES ($1, $2, 3, 1, 5, 2)
RETURNING *;

-- name: ReadLog :one
SELECT * FROM logs WHERE id = $1;

-- name: UpdateLog :one
UPDATE logs
SET rating_improved = 4, cards_mastered = 3
WHERE id = $1
RETURNING *;

-- name: DeleteLog :exec
DELETE FROM logs WHERE id = $1;
