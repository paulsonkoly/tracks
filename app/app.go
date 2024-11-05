package app

import (
	"context"
	"database/sql"
	"encoding/gob"
	"io"
	"log/slog"
	"net/http"

	"github.com/paulsonkoly/tracks/repository"
)

const TMPDir = "tmp"

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

// SessionManager provides access to the session store.
type SessionManager interface {

	// Get returns the value associated with the given key.
	//
	// If the key is not found, false is returned.
	Get(ctx context.Context, key string) any

	// Put associates the given value with the given key.
	Put(ctx context.Context, key string, value any)

	// Remove removes the value associated with the given key.
	Remove(ctx context.Context, key string)

	// Pop removes and returns the value associated with the given key.
	Pop(ctx context.Context, key string) any

	// RenewToken updates the session data to have a new session token while
	// retaining the current session data. The session lifetime is also reset and
	// the session data status will be set to Modified.
	//
	// The old session token and accompanying data are deleted from the session store.
	//
	// To mitigate the risk of session fixation attacks, it's important that you
	// call RenewToken before making any changes to privilege levels (e.g. login
	// and logout operations). See
	// [https://github.com/OWASP/CheatSheetSeries/blob/master/cheatsheets/Session\_Management\_Cheat\_Sheet.md#renew-the-session-id-after-any-privilege-level-change](https://github.com/OWASP/CheatSheetSeries/blob/master/cheatsheets/Session_Management_Cheat_Sheet.md#renew-the-session-id-after-any-privilege-level-change)
	// for additional information.
	RenewToken(ctx context.Context) error

	// LoadAndSave loads the session data from the session store.
	LoadAndSave(next http.Handler) http.Handler
}

type TXHandle = *sql.Tx

// Transaction manages per request database transaction.
type Transaction interface {
	// WithTx runs blk in a transaction. error will be returned if either there
	// was a problem with the transaction or blk returned an error. blk must pass
	// h to Repo() to use transaction.
	WithTx(ctx context.Context, blk func(h TXHandle) error) (err error)

	// Repo returns the query interface.
	// Pass nil if not within a transaction.
	Repo(h TXHandle) *repository.Queries
}

// Template provides html page rendering.
type Template interface {
	// Render produces html content identified by name and writes it to w. The data
	// carries page specific sideband data required on the page.
	Render(w io.Writer, name string, data any) error
}

type App struct {
	logger   Log
	txdb     Transaction
	sm       SessionManager
	template Template
	decoder  fDecoder
}

func New(logger Log, dbtx Transaction, sm SessionManager, tmpl Template) *App {
	gob.Register(Flash{})
	return &App{logger: logger, txdb: dbtx, sm: sm, template: tmpl, decoder: newDecoder()}
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
