package app

import (
	"context"

	"github.com/paulsonkoly/tracks/repository"
)

// WithTx executes blk inside a transaction. If blk fails with error the
// transaction is rolled back and the error is returned. The transaction handle
// is passed to blk and the blk code is supposed to pass this handle to [app.Repo].
func (a *App) WithTx(ctx context.Context, blk func(h TXHandle) error) (err error) {
	return a.txdb.WithTx(ctx, blk)
}

// Repo returns the database accessors / app data repository.
func (a *App) Repo(h TXHandle) *repository.Queries { return a.txdb.Repo(h) }
