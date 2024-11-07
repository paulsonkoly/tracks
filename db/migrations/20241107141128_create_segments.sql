-- this requires the postgis extension to be created by db superuser.

-- migrate:up
create table public."segments" (
  "id" serial primary key,
  "track_id" integer not null references public."tracks" on delete cascade,
  "geometry" geography(LINESTRING));

-- migrate:down
drop table public."segments";
