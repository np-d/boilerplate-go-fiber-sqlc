-- name: GetUser :one
SELECT *
FROM users
WHERE id = $1
  AND deleted_at IS NULL LIMIT 1;
-- name: GetUserByUsername :one 
SELECT *
FROM users
WHERE username = $1
  AND deleted_at IS NULL LIMIT 1;
-- name: GetUserByEmail :one 
SELECT *
FROM users
WHERE email = $1
  AND deleted_at IS NULL LIMIT 1;
-- name: GetUserDisplayName :one 
SELECT display_name
FROM users
WHERE id = $1
  AND deleted_at IS NULL;
-- name: CreateUser :one 
INSERT INTO users (display_name, username, email, password)
VALUES ($1, $2, $3, $4) RETURNING *;
-- name: UpdateUser :exec 
UPDATE users
SET display_name = $1,
    username     = $2,
    email        = $3
WHERE id = $4
  AND deleted_at IS NULL;
-- name: UpdateUserPassword :exec 
UPDATE users
SET password = $1
WHERE id = $2
  AND deleted_at IS NULL;
-- name: DeleteUser :exec 
UPDATE users
SET deleted_at = CURRENT_TIMESTAMP
WHERE id = $1;