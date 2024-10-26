package app

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/paulsonkoly/tracks/app/template"
	"github.com/paulsonkoly/tracks/repository"
)

type App struct {
	Logger   *slog.Logger
	Repo     *repository.Queries
	SM       *scs.SessionManager
	Template *template.Template
}

func New(logger *slog.Logger, repo *repository.Queries, sm *scs.SessionManager, tmpl *template.Template) *App {
	return &App{Logger: logger, Repo: repo, SM: sm, Template: tmpl}
}

func (a *App) ServerError(w http.ResponseWriter, msg string, err error) {
	a.Logger.Error(msg, "error", err.Error())
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func (a *App) AuthenticateUser(ctx context.Context, name, password string) (*repository.User, error) {
	user, err := a.Repo.GetUserByName(ctx, name)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	// TODO not hashed
	if err != nil || user.HashedPassword != password {
		return nil, err
	}
	return &user, nil
}
