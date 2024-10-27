package app

import (
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

func (a *App) ServerError(w http.ResponseWriter, msg string, err error) {
	a.Logger.Error(msg, "error", err.Error())
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

