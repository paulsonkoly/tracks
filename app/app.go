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

// Log defines the application logger.
type Log interface {

	// ServerError indicates that an internal server error happened.
	//
	// The handler should stop execution at this point, and internal server error
	// will be sent to the client and the error will be logged.
	ServerError(err error)

	// ClientError indicates that there was a problem with the client request.
	//
	// The handler should stop execution at this point, and the given status code
	// will be sent to the client and the error will be logged.
	ClientError(err error, status int)

	// Info creates a generic info level log message.
	//
	// args should be in pairs following the slog APIs.
	Info(msg string, args ...any)

	// Panic logs a recover() handling an uncought panic from a handler.
	//
	// dumps stack trace
	Panic(err any)
}

type App struct {
	logger   Log
	Repo     *repository.Queries
	SM       *scs.SessionManager
	Template *template.Cache
}

func New(logger Log, repo *repository.Queries, sm *scs.SessionManager, tmpl *template.Cache) *App {
	gob.Register(Flash{})
	return &App{logger: logger, Repo: repo, SM: sm, Template: tmpl}
}

func (a *App) ServerError(w http.ResponseWriter, err error) {
	a.logger.ServerError(err)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func (a *App) ClientError(w http.ResponseWriter, err error, status int) {
	a.logger.ClientError(err, status)
	http.Error(w, http.StatusText(status), status)
}

func (a *App) LogAction(ctx context.Context, action string, args ...any) {
	user := a.CurrentUser(ctx)
	if user != nil {
		args = append(args, slog.Int("actor id", int(user.ID)))
	}
	a.logger.Info(action, args...)
}
