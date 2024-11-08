// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: gpxfiles.sql

package repository

import (
	"context"
)

const deleteGPXFile = `-- name: DeleteGPXFile :one
delete from "public"."gpxfiles" where id = $1 returning filename
`

func (q *Queries) DeleteGPXFile(ctx context.Context, id int32) (string, error) {
	row := q.db.QueryRowContext(ctx, deleteGPXFile, id)
	var filename string
	err := row.Scan(&filename)
	return filename, err
}

const getGPXFile = `-- name: GetGPXFile :one
select id, filename, filesize, status, link, created_at from "public"."gpxfiles" where id = $1
`

func (q *Queries) GetGPXFile(ctx context.Context, id int32) (Gpxfile, error) {
	row := q.db.QueryRowContext(ctx, getGPXFile, id)
	var i Gpxfile
	err := row.Scan(
		&i.ID,
		&i.Filename,
		&i.Filesize,
		&i.Status,
		&i.Link,
		&i.CreatedAt,
	)
	return i, err
}

const getGPXFileByFilename = `-- name: GetGPXFileByFilename :one
select id, filename, filesize, status, link, created_at from "public"."gpxfiles" where filename = $1
`

func (q *Queries) GetGPXFileByFilename(ctx context.Context, filename string) (Gpxfile, error) {
	row := q.db.QueryRowContext(ctx, getGPXFileByFilename, filename)
	var i Gpxfile
	err := row.Scan(
		&i.ID,
		&i.Filename,
		&i.Filesize,
		&i.Status,
		&i.Link,
		&i.CreatedAt,
	)
	return i, err
}

const getGPXFiles = `-- name: GetGPXFiles :many
select id, filename, filesize, status, link, created_at from "public"."gpxfiles" order by created_at desc
`

func (q *Queries) GetGPXFiles(ctx context.Context) ([]Gpxfile, error) {
	rows, err := q.db.QueryContext(ctx, getGPXFiles)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Gpxfile
	for rows.Next() {
		var i Gpxfile
		if err := rows.Scan(
			&i.ID,
			&i.Filename,
			&i.Filesize,
			&i.Status,
			&i.Link,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const insertGPXFile = `-- name: InsertGPXFile :one
insert into "public"."gpxfiles" (filename, filesize, link, status, created_at) values ($1, $2, $3, 'uploaded', Now()) returning id
`

type InsertGPXFileParams struct {
	Filename string
	Filesize int64
	Link     string
}

func (q *Queries) InsertGPXFile(ctx context.Context, arg InsertGPXFileParams) (int32, error) {
	row := q.db.QueryRowContext(ctx, insertGPXFile, arg.Filename, arg.Filesize, arg.Link)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const setGPXFileStatus = `-- name: SetGPXFileStatus :exec
update "public"."gpxfiles" set status = $1 where id = $2
`

type SetGPXFileStatusParams struct {
	Status Filestatus
	ID     int32
}

func (q *Queries) SetGPXFileStatus(ctx context.Context, arg SetGPXFileStatusParams) error {
	_, err := q.db.ExecContext(ctx, setGPXFileStatus, arg.Status, arg.ID)
	return err
}
