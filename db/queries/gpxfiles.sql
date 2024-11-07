-- name: InsertGPXFile :one
insert into "public"."gpxfiles" (filename, filesize, link, status, created_at) values ($1, $2, $3, 'uploaded', Now()) returning id;

-- name: GetGPXFile :one
select * from "public"."gpxfiles" where id = $1;

-- name: GetGPXFileByFilename :one
select * from "public"."gpxfiles" where filename = $1;

-- name: GetGPXFiles :many
select * from "public"."gpxfiles" order by created_at desc;

-- name: DeleteGPXFile :one
delete from "public"."gpxfiles" where id = $1 returning filename;

-- name: SetGPXFileStatus :exec
update "public"."gpxfiles" set status = $1 where id = $2;
