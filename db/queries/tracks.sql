-- name: GetTracks :many
SELECT 
    sqlc.embed(t),
    sqlc.embed(u),
    SUM(ST_Length(s.geometry::geography))::double precision AS track_length_meters
FROM 
    public.tracks t
JOIN 
    public.segments s ON t.id = s.track_id
JOIN public.users u ON t.user_id = u.id
GROUP BY 
    t.id, u.id
ORDER BY 
    t.created_at desc;

-- name: InsertTrack :one
insert into "public"."tracks" (gpxfile_id, type, name, user_id) values ($1, $2, $3, $4) returning id;

-- name: GetTrack :one
select * from "public"."tracks" where id = $1;

