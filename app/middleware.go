package app

import (
	"log/slog"
	"net/http"

	"github.com/justinas/alice"
)

func (a *App) StandardChain() alice.Chain {
	return alice.New(a.Recover, a.SM.LoadAndSave, a.LogRequest, a.Headers)
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

