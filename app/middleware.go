package app

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"

	"github.com/justinas/alice"
	"github.com/justinas/nosurf"
)

type ContextKey string

var CurrentUser = ContextKey("CurrentUser")

func (a *App) StandardChain() alice.Chain {
	return alice.New(a.Recover, a.Dynamic, a.LogRequest, a.Headers, a.NoSurf)
}

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

func (a *App) LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.Logger.Info("request", slog.String("method", r.Method), slog.String("url", r.URL.Path), slog.String("remote", r.RemoteAddr))
		next.ServeHTTP(w, r)
	})
}

func (a *App) Dynamic(next http.Handler) http.Handler {
	return alice.New(a.SM.LoadAndSave).ThenFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if a.SM.Exists(ctx, currentUserID) {
			uid := a.SM.GetInt32(ctx, currentUserID)

			user, err := a.Repo.GetUser(ctx, uid)
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

func (a *App) NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	})
	return csrfHandler
}

func (a *App) Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				a.Logger.Error("panic", "error", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

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
