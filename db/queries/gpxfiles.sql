-- name: InsertGPXFile :one
insert into "public"."gpxfiles" (filename, filesize,  status, user_id, created_at) values ($1, $2, 'uploaded', $3, Now()) returning id;

-- name: UpdateGPXFile :exec
update "public"."gpxfiles" set
version=$2,
creator=$3,
name=$4,
description=$5,
author_name=$6,
author_email=$7,
author_link=$8,
author_link_text=$9,
author_link_type=$10,
copyright=$11,
copyright_year=$12,
copyright_license=$13,
link=$14,
link_text=$15,
link_type=$16,
time=$17,
keywords=$18
where id = $1;

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
