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

-- name: GetMatchingTracks :many
select t.id, t.name from "public"."tracks" t where t.name ilike $1;

-- name: InsertTrack :one
insert into "public"."tracks" (gpxfile_id, type, name, user_id) values ($1, $2, $3, $4) returning id;

-- name: GetTrack :one
select
   t.*,
   f.time,
  SUM(ST_Length(s.geometry::geography))::double precision AS track_length_meters
from "public"."tracks" t
join "public"."gpxfiles" f on t.gpxfile_id = f.id
join public.segments s ON t.id = s.track_id
where t.id = $1
group by t.id, f.time;

-- name: GetNonExistentTrackIDs :many
with s as (select unnest(@track_ids::integer[]) id)
select s.id::integer
from s
left join tracks as t on s.id = t.id
where t.id is null;

-- name: GetTrackSegments :many
select id from segments where track_id = $1;

