-- name: InsertGPXFile :exec
insert into "public"."gpxfiles" (filename, filesize, link, status, created_at) values ($1, $2, $3, 'uploaded', Now());

-- name: GetGPXFile :one
select * from "public"."gpxfiles" where id = $1;

-- name: GetGPXFileByFilename :one
select * from "public"."gpxfiles" where filename = $1;
