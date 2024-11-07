-- name: InsertTrack :exec
insert into "public"."tracks" (gpxfile_id, type, name) values ($1, $2, $3);
