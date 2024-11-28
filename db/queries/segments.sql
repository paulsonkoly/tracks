-- name: InsertSegment :exec
insert into "public"."segments" (track_id, geometry) values ($1, $2);

-- name: GetSegmentPoints :many
SELECT latitude::float, longitude::float FROM points where segment_id=$1;
