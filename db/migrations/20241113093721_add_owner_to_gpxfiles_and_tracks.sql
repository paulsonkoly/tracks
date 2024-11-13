-- migrate:up
alter table "public"."gpxfiles" add column user_id integer;
alter table "public"."tracks" add column user_id integer;

-- fill existing rows with first user's id 
update "public"."gpxfiles" set user_id = (select id from "public"."users" order by id asc limit 1);
update "public"."tracks" set user_id = (select id from "public"."users" order by id asc limit 1);

-- add foreign key constraints
alter table "public"."gpxfiles" add constraint "gpxfiles_user_id_fkey" foreign key ("user_id") references "public"."users"("id") on update cascade on delete cascade;
alter table "public"."gpxfiles" alter column "user_id" set not null;
alter table "public"."tracks" add constraint "tracks_user_id_fkey" foreign key ("user_id") references "public"."users"("id") on update cascade on delete cascade;
alter table "public"."tracks" alter column "user_id" set not null;

-- migrate:down
alter table "public"."gpxfiles" drop constraint "gpxfiles_user_id_fkey";
alter table "public"."tracks" drop constraint "tracks_user_id_fkey";
alter table "public"."gpxfiles" drop column "user_id";
alter table "public"."tracks" drop column "user_id";

