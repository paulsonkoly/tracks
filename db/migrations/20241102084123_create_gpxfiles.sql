-- migrate:up

create type filestatus as enum ('uploaded', 'processed', 'processing_failed');

create table "public"."gpxfiles" (
  "id" serial primary key,
  "filename" text not null,
  "filesize" bigint not null,
  "status" filestatus not null,
  "link" text not null default '',
  "created_at" timestamp with time zone not null);

-- otherwise we can't guarantee uploads not overwriting each other
alter table "public"."gpxfiles" add constraint gpxfiles_filename_key unique (filename);

comment on column "public"."gpxfiles"."link" is 'gpx metadata link field';

-- migrate:down
drop table "public"."gpxfiles";

drop type filestatus;
