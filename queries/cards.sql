-- name: CreateCard :one
INSERT INTO cards (question, answer, pack_id, rating)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: ReadCard :one
SELECT * FROM cards WHERE id = $1;

-- name: UpdateCard :one
UPDATE cards
SET question = $1, answer = $2, rating = 2
WHERE id = $3
RETURNING *;

-- name: ListCardsByPack :many
SELECT *
FROM cards
WHERE pack_id = $1
ORDER BY created_at DESC;

-- name: ListRepeatCards :many
SELECT * FROM cards
WHERE pack_id = $1
ORDER BY last_wrong DESC, created_at ASC;

-- name: MarkCardWrong :exec
UPDATE cards
SET last_wrong = $1
WHERE id = $2;

-- name: DeleteCard :exec
DELETE FROM cards WHERE id = $1;
