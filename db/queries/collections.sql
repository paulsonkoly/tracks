-- name: InsertCollection :exec
with c as (insert into collections (name, user_id) values ($1, $2) returning id)
insert into track_collections (collection_id, track_id)
(select c.id , unnest(@track_ids::integer[]) from c);

-- name: GetCollection :one
select name from collections where id = $1;

-- name: GetCollectionByName :one
select * from collections where name = $1;

-- name: GetCollectionSegments :many
select s.id from
collections c
inner join track_collections tc on tc.collection_id = c.id
inner join tracks t on tc.track_id = t.id
inner join segments s on s.track_id = t.id
where c.id = $1;

-- name: GetCollections :many
select c.id, c.name from collections c order by c.id;
