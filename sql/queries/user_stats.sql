-- name: CreateUserStats :exec
INSERT INTO user_stats (user_id) VALUES ($1)
ON CONFLICT DO NOTHING;

-- name: AddUserRating :exec
UPDATE user_stats
SET rating = rating + $1
WHERE user_id = $2;

-- name: IncPacksCreated :exec
UPDATE user_stats
SET packs_created = packs_created + 1
WHERE user_id = $1;

-- name: IncPacksMastered :exec
UPDATE user_stats
SET packs_mastered = packs_mastered + 1
WHERE user_id = $1;

-- name: GetUserStats :one
SELECT rating, packs_created, packs_mastered
FROM user_stats
WHERE user_id = $1;
