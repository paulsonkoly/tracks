package repository

import (
	"context"
	"database/sql"
	"errors"
)

func (q *Queries) Unique(ctx context.Context, filename string) (bool, error) {
	_, err := q.GetGPXFileByFilename(ctx, filename)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return true, nil
		}
		return false, err
	}
	return false, nil
}
