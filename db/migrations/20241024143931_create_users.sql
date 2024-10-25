-- migrate:up
create table users (
  id serial primary key,
  username varchar(255) not null,
  hashed_password varchar(255) not null,
  created_at timestamp with time zone not null);

alter table users add constraint users_username_key unique (username);

-- migrate:down
drop table users;

