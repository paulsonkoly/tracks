-- migrate:up
create table sessions (
  token char(43) primary key,
  data bytea not null,
  expiry timestamp not null
);

create index sessions_expiry_idx on sessions (expiry);

-- migrate:down
drop table sessions;

