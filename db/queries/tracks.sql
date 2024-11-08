-- name: InsertTrack :one
insert into "public"."tracks" (gpxfile_id, type, name) values ($1, $2, $3) returning id;

-- name: GetTrack :one
select * from "public"."tracks" where id = $1;

