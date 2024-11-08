-- name: InsertSegment :exec
insert into "public"."segments" (track_id, geometry) values ($1, $2);

-- name: GetTrackSegmentPoints :many
SELECT * FROM points where track_id=$1;
