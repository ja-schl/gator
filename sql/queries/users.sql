-- name: CreateUser :one
INSERT INTO users (
	id, created_at, updated_at, name
) VALUES ( $1, $2, $3, $4 ) returning *;

-- name: GetUser :one
SELECT * FROM users
where name = $1;

-- name: GetUsers :many
Select * FROM users;

-- name: GetUserById :one
select * from users where id = $1;

-- name: DeleteUsers :exec
DELETE FROM users;