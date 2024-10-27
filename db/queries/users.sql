-- name: GetUser :one
select * from users where id = $1;

-- name: GetUserByName :one
select * from users where username = $1;

-- name: GetUsers :many
select * from users;

-- name: InsertUser :one
insert into users (username, hashed_password, created_at) values ($1, $2, Now()) returning *;

-- name: DeleteUser :exec
delete from users where id = $1;
