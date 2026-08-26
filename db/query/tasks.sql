-- name: CreateTask :one
INSERT INTO "tasks" (
    title
) VALUES (
    $1
) RETURNING *;

-- name: ListTasks :many
SELECT * FROM "tasks"
WHERE deleted_at IS NULL
ORDER BY id;

-- name: GetTask :one
SELECT * FROM "tasks"
WHERE id = $1 
AND deleted_at IS NULL;

-- name: DeleteTask :exec
UPDATE "tasks"
SET deleted_at = now(), updated_at = now()
WHERE id = $1
AND deleted_at IS NULL;

-- name: UpdateTask :one
UPDATE "tasks"
SET title = $1, updated_at = now()
WHERE id = $2
AND deleted_at IS NULL
RETURNING *;

-- name: CompleteTask :one
UPDATE "tasks"
SET completed = $1, updated_at = now()
WHERE id = $2
AND deleted_at IS NULL
RETURNING *;
