-- name: InsertSegment :exec
insert into "public"."segments" (track_id, geometry) values ($1, $2);

