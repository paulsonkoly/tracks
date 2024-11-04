package tx

import (
	"context"
	"database/sql"

	"github.com/paulsonkoly/tracks/repository"
)

type Handle = *sql.Tx

type TX struct {
	repo *repository.Queries
	db   *sql.DB
}

func New(repo *repository.Queries, db *sql.DB) *TX {
	return &TX{
		repo: repo,
		db:   db,
	}
}

func (tx *TX) WithTx(ctx context.Context, blk func(h Handle) error) (err error) {
	t, err := tx.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err = blk(t); err != nil {
		_ = t.Rollback()
		return err
	}

	return t.Commit()
}

func (tx *TX) Repo(h Handle) *repository.Queries {
	if h == nil {
		return tx.repo
	}

	return tx.repo.WithTx(h)
}
