-- migrate:up
alter table "public"."gpxfiles" add column "version" text;
alter table "public"."gpxfiles" add column "creator" text;
alter table "public"."gpxfiles" add column "name" text;
alter table "public"."gpxfiles" add column "description" text;
alter table "public"."gpxfiles" add column "author_name" text;
alter table "public"."gpxfiles" add column "author_email" text;
alter table "public"."gpxfiles" add column "author_link" text;
alter table "public"."gpxfiles" add column "author_link_text" text;
alter table "public"."gpxfiles" add column "author_link_type" text;
alter table "public"."gpxfiles" add column "copyright" text;
alter table "public"."gpxfiles" add column "copyright_year" text;
alter table "public"."gpxfiles" add column "copyright_license" text;
alter table "public"."gpxfiles" add column "link_text" text;
alter table "public"."gpxfiles" add column "link_type" text;
alter table "public"."gpxfiles" add column "time" timestamp with time zone;
alter table "public"."gpxfiles" add column "keywords" text;

-- migrate:down
alter table "public"."gpxfiles" drop column "version";
alter table "public"."gpxfiles" drop column "creator";
alter table "public"."gpxfiles" drop column "name";
alter table "public"."gpxfiles" drop column "description";
alter table "public"."gpxfiles" drop column "author_name";
alter table "public"."gpxfiles" drop column "author_email";
alter table "public"."gpxfiles" drop column "author_link";
alter table "public"."gpxfiles" drop column "author_link_text";
alter table "public"."gpxfiles" drop column "author_link_type";
alter table "public"."gpxfiles" drop column "copyright";
alter table "public"."gpxfiles" drop column "copyright_year";
alter table "public"."gpxfiles" drop column "copyright_license";
alter table "public"."gpxfiles" drop column "link_text";
alter table "public"."gpxfiles" drop column "link_type";
alter table "public"."gpxfiles" drop column "time";
alter table "public"."gpxfiles" drop column "keywords";
