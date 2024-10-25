package app

import (
	"database/sql"
	"log/slog"
)

type App struct {
	Logger *slog.Logger
	DB     *sql.DB
}

func NewApp(logger *slog.Logger, db *sql.DB) *App {
	return &App{Logger: logger, DB: db}
}
