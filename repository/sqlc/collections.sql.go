// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: collections.sql

package sqlc

import (
	"context"

	"github.com/lib/pq"
)

const getCollection = `-- name: GetCollection :one
select name from collections where id = $1
`

func (q *Queries) GetCollection(ctx context.Context, id int32) (string, error) {
	row := q.db.QueryRowContext(ctx, getCollection, id)
	var name string
	err := row.Scan(&name)
	return name, err
}

const getCollectionByName = `-- name: GetCollectionByName :one
select id, name, user_id from collections where name = $1
`

func (q *Queries) GetCollectionByName(ctx context.Context, name string) (Collection, error) {
	row := q.db.QueryRowContext(ctx, getCollectionByName, name)
	var i Collection
	err := row.Scan(&i.ID, &i.Name, &i.UserID)
	return i, err
}

const getCollectionSegments = `-- name: GetCollectionSegments :many
select s.id from
collections c
inner join track_collections tc on tc.collection_id = c.id
inner join tracks t on tc.track_id = t.id
inner join segments s on s.track_id = t.id
where c.id = $1
`

func (q *Queries) GetCollectionSegments(ctx context.Context, id int32) ([]int32, error) {
	rows, err := q.db.QueryContext(ctx, getCollectionSegments, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []int32
	for rows.Next() {
		var id int32
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		items = append(items, id)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const insertCollection = `-- name: InsertCollection :exec
with c as (insert into collections (name, user_id) values ($1, $2) returning id)
insert into track_collections (collection_id, track_id)
(select c.id , unnest($3::integer[]) from c)
`

type InsertCollectionParams struct {
	Name     string
	UserID   int32
	TrackIds []int32
}

func (q *Queries) InsertCollection(ctx context.Context, arg InsertCollectionParams) error {
	_, err := q.db.ExecContext(ctx, insertCollection, arg.Name, arg.UserID, pq.Array(arg.TrackIds))
	return err
}
