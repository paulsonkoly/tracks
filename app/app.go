package app

import (
	"context"
	"encoding/gob"
	"log/slog"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/paulsonkoly/tracks/app/template"
	"github.com/paulsonkoly/tracks/repository"
)

const currentUserID = "currentUserID"

type App struct {
	Logger   *slog.Logger
	Repo     *repository.Queries
	SM       *scs.SessionManager
	Template *template.Template
}

func New(logger *slog.Logger, repo *repository.Queries, sm *scs.SessionManager, tmpl *template.Template) *App {
	gob.Register(Flash{})
	return &App{Logger: logger, Repo: repo, SM: sm, Template: tmpl}
}

func (a *App) ServerError(w http.ResponseWriter, err error) {
	a.Logger.Error("server error", "error", err.Error())
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func (a *App) ClientError(w http.ResponseWriter, err error, status int) {
	a.Logger.Debug("client error", "error", err.Error(), "status", status)
	http.Error(w, http.StatusText(status), status)
}

func (a *App) LogAction(ctx context.Context, action string, args ...any) {
	user := a.CurrentUser(ctx)
	if user != nil {
		args = append(args, slog.Int("actor id", int(user.ID)))
	}
	a.Logger.Info(action, args...)
}
