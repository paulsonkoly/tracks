// Package repository provides transactions and data manipulation.
package repository

import (
	"context"
	"database/sql"

	"github.com/paulsonkoly/tracks/repository/sqlc"
)

// Repository implements data repository and transactional access.
type Repository struct {
	sqlc *sqlc.Queries
	db   *sql.DB
}

// New creates new repository.
func New(sqlc *sqlc.Queries, db *sql.DB) Repository {
	return Repository{
		sqlc: sqlc,
		db:   db,
	}
}

const txHandle = ContextKey("tracksTxHandle")

type ContextKey string

// WithTx executes the given block in a transaction.
func (r Repository) WithTx(ctx context.Context, blk func(ctx context.Context) error) error {
	t, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	ctx = context.WithValue(ctx, ContextKey(txHandle), t)

	if err = blk(ctx); err != nil {
		_ = t.Rollback()
		return err
	}

	return t.Commit()
}

// Queries is the repository Query interface.
type Queries struct {
	// This allows automatically plugging ctx in the query methods on Query. This
	// is not an anti-pattern as Q should never be assigned to a variable and we
	// should always just call query methos on repository via Q:
	// repo.Q(ctx).GetUser(params). The only diff is it would be
	// repo.Q(ctx).GetUser(ctx, params).
	ctx  context.Context
	sqlc *sqlc.Queries
}

// Q returns the Query interface for repository access.
func (r Repository) Q(ctx context.Context) Queries {
	// the big debate here is passing the active transaction around in the
	// context or using goroutine local storage to store which transaction we are
	// in or not.
	//
	// For now we use the context.
	h, ok := ctx.Value(ContextKey(txHandle)).(*sql.Tx)
	if !ok {
		return Queries{ctx: ctx, sqlc: r.sqlc}
	}
	return Queries{ctx: ctx, sqlc: r.sqlc.WithTx(h)}
}

// Point is a pair of longitude and latitude.
type Point struct {
	Longitude, Latitude float64
}

// Segment is a line segment, a sequence of [Point].
type Segment []Point
