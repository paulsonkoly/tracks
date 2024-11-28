-- migrate:up
drop view if exists points;

CREATE VIEW points AS (
SELECT 
  id as segment_id,
  CAST(ST_X(geom) AS double precision) AS longitude,
  CAST(ST_Y(geom) AS double precision) AS latitude
FROM (SELECT (ST_DumpPoints(geometry::geometry)).geom, id FROM segments) AS S);


-- migrate:down
drop view if exists points;

CREATE VIEW points AS (
SELECT 
  track_id,
  CAST(ST_X(geom) AS double precision) AS longitude,
  CAST(ST_Y(geom) AS double precision) AS latitude
FROM (SELECT (ST_DumpPoints(geometry::geometry)).geom, track_id FROM segments) AS S);
