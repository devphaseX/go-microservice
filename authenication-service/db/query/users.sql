-- name: GetAll :many

SELECT * FROM Users
OFFSET sqlc.arg('offset')
LIMIT sqlc.arg('limit');

-- name: GetUserByUnqiueID :one
SELECT * FROM Users
WHERE sqlc.narg('id') = id or sqlc.narg('email') ilike email
LIMIT 1;

-- name: CreateUser :one
INSERT INTO Users (
    first_name,
    last_name,
    email,
    password
) VALUES (
    sqlc.arg('first_name'),
    sqlc.arg('last_name'),
    sqlc.arg('email'),
    sqlc.arg('password')
) RETURNING *;

-- name: UpdateUser :one
UPDATE Users
SET first_name = sqlc.arg('first_name'),
    last_name = sqlc.arg('last_name'),
    user_active = sqlc.arg('user_active')
WHERE id = sqlc.arg('id')
RETURNING *;


-- name: UpdateUserEmail :one
UPDATE Users
SET email = sqlc.arg('email')
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteUsers :many
DELETE FROM Users
WHERE (sqlc.arg('ids')::uuid[] IS NOT NULL AND id = ANY(sqlc.arg('ids')::uuid[])) OR
      (sqlc.arg('emails')::text[] IS NOT NULL AND email = ANY(sqlc.arg('emails')::text[]))
RETURNING *;

-- name: DeleteUserByUnqiueID :one
DELETE FROM Users
WHERE sqlc.arg('id') = id OR
      sqlc.arg('email') ilike email
RETURNING *;
