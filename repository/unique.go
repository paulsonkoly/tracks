package repository

import (
	"context"
	"database/sql"
	"errors"
)

func (q *Queries) GPXFileUnique(ctx context.Context, filename string) (bool, error) {
	_, err := q.GetGPXFileByFilename(ctx, filename)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return true, nil
		}
		return false, err
	}
	return false, nil
}

func (q *Queries) UserUnique(ctx context.Context, username string) (bool, error) {
	_, err := q.GetUserByName(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return true, nil
		}
		return false, err
	}
	return false, nil
}

func (q *Queries) UserUniqueExceptID(ctx context.Context, id int, username string) (bool, error) {
	_, err := q.GetUserByNameNotID(ctx, GetUserByNameNotIDParams{ID: int32(id), Username: username})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return true, nil
		}
		return false, err
	}
	return false, nil
}
