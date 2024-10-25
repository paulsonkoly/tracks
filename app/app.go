package app

import "log/slog"

type App struct {
	Logger *slog.Logger
}

func NewApp(logger *slog.Logger) *App {
	return &App{Logger: logger}
}
