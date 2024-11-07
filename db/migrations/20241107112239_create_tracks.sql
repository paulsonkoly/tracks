-- migrate:up
create type tracktype as enum ('track', 'route');

create table public."tracks" (
  "id" serial not null primary key,
  "name" text not null default '',
  "type" tracktype not null,
  "gpxfile_id" integer not null references public."gpxfiles" on delete cascade);

-- migrate:down
drop table public."tracks";

drop type public."tracktype";
