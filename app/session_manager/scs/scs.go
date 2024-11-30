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

// TokenRenewalError indicates and error while renewing the session token.
type TokenRenewalError struct {
	scsErr error
}

// Error is the error message.
func (err TokenRenewalError) Error() string {
	return err.scsErr.Error()
}

// SessionManager is a session managaer piggybacking on the scs package.
type SessionManager struct {
	scs *scs.SessionManager
}

// New creates a session manager.
func New(db *sql.DB) SessionManager {
	scs := scs.New()
	scs.Store = postgresstore.New(db)
	return SessionManager{
		scs: scs,
	}
}

// Get retrieves an object from the session.
func (sm SessionManager) Get(ctx context.Context, key string) any { return sm.scs.Get(ctx, key) }

// Put puts an object in the session.
func (sm SessionManager) Put(ctx context.Context, key string, value any) { sm.scs.Put(ctx, key, value) }

// Remove removes an object from the session.
func (sm SessionManager) Remove(ctx context.Context, key string) { sm.scs.Remove(ctx, key) }

// Pop removes an object from the session and returns it.
func (sm SessionManager) Pop(ctx context.Context, key string) any { return sm.scs.Pop(ctx, key) }

// RenewToken renews the session token.
func (sm SessionManager) RenewToken(ctx context.Context) error {
	err := sm.scs.RenewToken(ctx)
	if err != nil {
		return fmt.Errorf("renew token error %w", TokenRenewalError{err})
	}
	return nil
}

// LoadAndSave is a middleware that loads the session.
func (sm SessionManager) LoadAndSave(next http.Handler) http.Handler { return sm.scs.LoadAndSave(next) }
