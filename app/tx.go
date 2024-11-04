package app

import (
	"context"

	"github.com/paulsonkoly/tracks/repository"
)

func (a *App) WithTx(ctx context.Context, blk func(h TXHandle) error) (err error) {
	return a.txdb.WithTx(ctx, blk)
}

func (a *App) Repo(h TXHandle) *repository.Queries { return a.txdb.Repo(h) }
