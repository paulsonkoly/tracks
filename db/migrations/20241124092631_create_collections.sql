-- migrate:up
create table "public".collections (
  "id" serial primary key,
  "name" text not null unique,
  "user_id" integer not null references public."users" on delete cascade
);

create table "public".track_collections (
  "track_id" integer not null references public."tracks" on delete cascade,
  "collection_id" integer not null references public."collections" on delete cascade,
  unique("track_id", "collection_id")
);

-- migrate:down
drop table public."track_collections";
drop table public."collections";

