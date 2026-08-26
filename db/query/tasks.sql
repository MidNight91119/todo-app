-- name: CreateTask :one
INSERT INTO "tasks" (
    title
) VALUES (
    $1
) RETURNING *;