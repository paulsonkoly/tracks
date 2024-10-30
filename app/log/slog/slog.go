// Package slog uses log/slog back end for [pkg/github.com/paulsonkoly/tracks/app.Log] interface.
package slog

import (
	"log/slog"
	"os"
	"runtime/debug"
)

type SLog struct {
	logger *slog.Logger
}

func New() SLog {
	return SLog{logger: slog.New(slog.NewTextHandler(os.Stderr, nil))}
}

func (s SLog) ServerError(err error) {
	s.logger.Error("server error", "error", err.Error())
}

func (s SLog) ClientError(err error, status int) {
	s.logger.Debug("client error", "error", err.Error(), "status", status)
}

func (s SLog) Info(msg string, args ...any) {
	s.logger.Info(msg, args...)
}

func (s SLog) Panic(err any) {
	s.logger.Error("panic", "error", err, "stack", debug.Stack())
}
