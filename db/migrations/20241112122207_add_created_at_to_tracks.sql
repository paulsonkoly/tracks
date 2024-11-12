-- migrate:up
alter table "public"."tracks" add column "created_at" timestamp with time zone not null default now();

-- migrate:down
alter table "public"."tracks" drop column "created_at";

