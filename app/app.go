package app

import (
	"database/sql"
	"log/slog"

	"github.com/alexedwards/scs/v2"
)

type App struct {
	Logger *slog.Logger
	DB     *sql.DB
	SM     *scs.SessionManager
}

func New(logger *slog.Logger, db *sql.DB, sm *scs.SessionManager) *App {
	return &App{Logger: logger, DB: db, SM: sm}
}
