-- name: GetTracks :many
SELECT 
    sqlc.embed(t),
    SUM(ST_Length(s.geometry::geography))::double precision AS track_length_meters
FROM 
    public.tracks t
JOIN 
    public.segments s ON t.id = s.track_id
GROUP BY 
    t.id
ORDER BY 
    t.created_at desc;

-- name: InsertTrack :one
insert into "public"."tracks" (gpxfile_id, type, name) values ($1, $2, $3) returning id;

-- name: GetTrack :one
select * from "public"."tracks" where id = $1;

