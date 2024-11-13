-- name: InsertGPXFile :one
insert into "public"."gpxfiles" (filename, filesize, link, status, user_id, created_at) values ($1, $2, $3, 'uploaded', $4, Now()) returning id;

-- name: GetGPXFile :one
select * from "public"."gpxfiles" where id = $1;

-- name: GetGPXFileByFilename :one
select * from "public"."gpxfiles" where filename = $1;

-- name: GetGPXFiles :many
select
  f.id, f.filename, f.filesize, f.link, f.status, f.created_at,
  sqlc.embed(u)
from "public"."gpxfiles" f
join "public"."users" u on f."user_id" = u."id"
order by f.created_at desc;

-- name: DeleteGPXFile :one
delete from "public"."gpxfiles" where id = $1 returning filename;

-- name: SetGPXFileStatus :exec
update "public"."gpxfiles" set status = $1 where id = $2;
