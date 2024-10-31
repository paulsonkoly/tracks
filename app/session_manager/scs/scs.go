// Package scs uses alexedwards/scs for session management.
//
// This implementation fulfills the [pkg/github.com/paulsonkoly/tracks/app.SessionManager] interface.
package scs

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
)

type TokenRenewalError struct {
	scsErr error
}

func (err TokenRenewalError) Error() string {
	return err.scsErr.Error()
}

type SessionManager struct {
	scs *scs.SessionManager
}

func New(db *sql.DB) SessionManager {
	scs := scs.New()
	scs.Store = postgresstore.New(db)
	return SessionManager{
		scs: scs,
	}
}

func (sm SessionManager) Get(ctx context.Context, key string) any { return sm.scs.Get(ctx, key) }

func (sm SessionManager) Put(ctx context.Context, key string, value any) { sm.scs.Put(ctx, key, value) }

func (sm SessionManager) Remove(ctx context.Context, key string) { sm.scs.Remove(ctx, key) }

func (sm SessionManager) Pop(ctx context.Context, key string) any { return sm.scs.Pop(ctx, key) }

func (sm SessionManager) RenewToken(ctx context.Context) error {
	err := sm.scs.RenewToken(ctx)
	if err != nil {
		return fmt.Errorf("renew token error %w", TokenRenewalError{err})
	}
	return nil
}

func (sm SessionManager) LoadAndSave(next http.Handler) http.Handler { return sm.scs.LoadAndSave(next) }
