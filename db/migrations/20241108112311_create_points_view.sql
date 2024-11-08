-- migrate:up
CREATE VIEW points AS (
SELECT 
  track_id,
  CAST(ST_X(geom) AS double precision) AS longitude,
  CAST(ST_Y(geom) AS double precision) AS latitude
FROM (SELECT (ST_DumpPoints(geometry::geometry)).geom, track_id FROM segments) AS S);

-- migrate:down
DROP VIEW IF EXISTS points;
