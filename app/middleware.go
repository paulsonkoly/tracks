package app

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/justinas/alice"
	"github.com/justinas/nosurf"
)

// ContextKey is used to store data in the request context.
type ContextKey string

// CurrentUser is the context key for the currently logged in user if any. By
// having the current user in the context we avoid loading it many times
// throughout the lifespan of the request.
var CurrentUser = ContextKey("CurrentUser")

// StandardChain is a middleware to be used for every request that does not require further special treatment.
func (a *App) StandardChain() alice.Chain {
	return alice.New(a.Recover, a.Dynamic, a.LogRequest, a.Headers, a.NoSurf)
}

// Headers is a middleware that adds common http headers.
func (a *App) Headers(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(
			"Content-Security-Policy",
			"default-src 'self'; "+
				"style-src 'self'; "+
				"font-src fonts.gstatic.com; img-src 'self' a.tile.opentopomap.org b.tile.opentopomap.org c.tile.opentopomap.org")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		w.Header().Set("Server", "Go")
		next.ServeHTTP(w, r)
	})
}

// LogRequest is a middleware that logs the request.
func (a *App) LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.logger.Info("request", "method", r.Method, "url", r.URL.Path, "remote", r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

// Dynamic is a middleware that enables sessions, and loads the current user.
func (a *App) Dynamic(next http.Handler) http.Handler {
	return alice.New(a.sm.LoadAndSave).ThenFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if uid, ok := a.sm.Get(ctx, SKCurrentUserID).(int); ok {
			user, err := a.Q(ctx).GetUser(uid)
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				a.ServerError(w, err)
				return
			}

			if err == nil {
				ctx = context.WithValue(ctx, CurrentUser, user)
				r = r.WithContext(ctx)
			}
		}

		next.ServeHTTP(w, r)
	})
}

// NoSurf is a middleware that gives CSRF protection.
func (a *App) NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	})
	return csrfHandler
}

// Recover is a middleware that catches any uncaught panic, and logs a backtrace.
func (a *App) Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				a.logger.Panic(err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// RequiresLogIn is a middleware that enforces the user to be logged in.
func (a *App) RequiresLogIn(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if nil == a.CurrentUser(r.Context()) {
			a.FlashError(r.Context(), "Requires login")
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}
