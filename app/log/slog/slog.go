// Package slog uses log/slog back end for [pkg/github.com/paulsonkoly/tracks/app.Log] interface.
package slog

import (
	"log/slog"
	"os"
	"runtime/debug"
)

// SLog creates [app.Logger] based on [slog].
type SLog struct {
	logger *slog.Logger
}

// New creates a new [SLog].
func New() SLog {
	return SLog{logger: slog.New(slog.NewTextHandler(os.Stderr, nil))}
}

// ServerError logs a server error.
func (s SLog) ServerError(err error) {
	s.logger.Error("server error", "error", err.Error())
}

// ClientError logs a client error.
func (s SLog) ClientError(err error, status int) {
	s.logger.Debug("client error", "error", err.Error(), "status", status)
}

// Info logs some debug information.
func (s SLog) Info(msg string, args ...any) {
	s.logger.Info(msg, args...)
}

// Panic should not be used in handlers. It is invoked from the middleware that
// recovers panicing handlers.
func (s SLog) Panic(err any) {
	s.logger.Error("panic", "error", err, "stack", debug.Stack())
}
