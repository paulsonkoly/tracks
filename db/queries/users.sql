-- name: GetUser :one
select * from users where id = $1;

-- name: GetUserByName :one
select * from users where username = $1;
