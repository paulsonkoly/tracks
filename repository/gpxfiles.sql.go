// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: gpxfiles.sql

package repository

import (
	"context"
)

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

const insertGPXFile = `-- name: InsertGPXFile :exec
insert into "public"."gpxfiles" (filename, filesize, link, status, created_at) values ($1, $2, $3, 'uploaded', Now())
`

type InsertGPXFileParams struct {
	Filename string
	Filesize int64
	Link     string
}

func (q *Queries) InsertGPXFile(ctx context.Context, arg InsertGPXFileParams) error {
	_, err := q.db.ExecContext(ctx, insertGPXFile, arg.Filename, arg.Filesize, arg.Link)
	return err
}
