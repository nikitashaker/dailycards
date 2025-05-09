-- name: CreatePack :one
INSERT INTO packs (name, category)
VALUES ($1, $2)
RETURNING *;

-- name: ReadPack :one
SELECT * FROM packs WHERE id = $1;

-- name: UpdatePack :one
UPDATE packs
SET name = $1, category = $2
WHERE id = $3
RETURNING *;

-- name: ListPacks :many
SELECT id, name, category, created_at, updated_at
FROM packs
ORDER BY created_at DESC;

-- name: DeletePack :exec
DELETE FROM packs WHERE id = $1;

