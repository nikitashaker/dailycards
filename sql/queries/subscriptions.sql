-- name: CreateSubscription :one
INSERT INTO subscriptions (user_id, pack_id)
VALUES ($1, $2)
RETURNING *;

-- name: ReadSubscription :one
SELECT * FROM subscriptions WHERE id = $1;

-- name: UpdateSubscription :one
UPDATE subscriptions
SET user_id = $1, pack_id = $2
WHERE id = $3
RETURNING *;

-- name: DeleteSubscription :exec
DELETE FROM subscriptions WHERE id = $1;
