package handler_test

import (
	"database/sql"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/paulsonkoly/tracks/app"
	"github.com/paulsonkoly/tracks/app/handler"
	"github.com/paulsonkoly/tracks/app/template"
	"github.com/paulsonkoly/tracks/repository"
	"github.com/paulsonkoly/tracks/repository/sqlc"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

type noLogger struct{}

func (noLogger) ServerError(_ error)        {}
func (noLogger) ClientError(_ error, _ int) {}
func (noLogger) Info(_ string, _ ...any)    {}
func (noLogger) Panic(_ any)                {}

type mockSession struct {
	store map[any]any
}

func newMockSession() mockSession {
	return mockSession{store: make(map[any]any)}
}

func (s mockSession) Get(_ context.Context, key string) any {
	val, ok := s.store[key]
	if ok {
		return val
	}
	return nil
}

func (s mockSession) Put(_ context.Context, key string, val any) {
	s.store[key] = val
}

func (s mockSession) Remove(_ context.Context, key string) {
	delete(s.store, key)
}

func (s mockSession) Pop(_ context.Context, key string) any {
	val, ok := s.store[key]
	if ok {
		delete(s.store, key)
		return val
	}
	return nil
}

func (s mockSession) Has(t *testing.T, key string, expected any) {
	actual, ok := s.store[key]
	assert.True(t, ok, "session expected to contain key %q, but it was not found", key)
	if ok && expected != nil {
		assert.Equal(t, expected, actual)
	}
}

func (s mockSession) DoesntHave(t *testing.T, key string) {
	assert.NotContains(t, key, s.store)
}

func (s mockSession) HasFlashInfo(t *testing.T, msg string) {
	flash, ok := s.store[app.SKFlash]
	assert.True(t, ok, "session expected to contain flash, but it was not found")

	assert.IsType(t, app.Flash{}, flash)

	if flash, ok := flash.(app.Flash); ok {
		assert.Contains(t, flash, "info")
		infos, ok := flash["info"]
		assert.True(t, ok, "flash expected to contain infos, but it was not found")

		assert.Contains(t, infos, msg)
	}
}

func (mockSession) RenewToken(_ context.Context) error         { return nil }
func (mockSession) LoadAndSave(next http.Handler) http.Handler { return next }

func init() {
	// templates cannot be picked up if we are not in project root
	_, filename, _, _ := runtime.Caller(0)
	err := os.Chdir(filepath.Join(filepath.Dir(filename), "..", ".."))
	if err != nil {
		panic(err)
	}
}

func withDB(t *testing.T, f func(*sql.DB)) {
	// connect to postgres cluster
	db, err := sql.Open("postgres", "user=tracks_test password=1234567 dbname=postgres sslmode=disable")
	assert.NoError(t, err)

	// drop database if exists
	_, err = db.Exec("DROP DATABASE IF EXISTS tracks_test")
	assert.NoError(t, err)

	// create database
	_, err = db.Exec("CREATE DATABASE tracks_test")
	assert.NoError(t, err)

	// close connection to cluster
	db.Close()

	// open connection to test db
	db, err = sql.Open("postgres", "user=tracks_test password=1234567 dbname=tracks_test sslmode=disable")
	assert.NoError(t, err)

	defer db.Close()

	// load schema file
	schema, err := os.ReadFile("db/schema.sql")
	assert.NoError(t, err)

	_, err = db.Exec(string(schema))
	assert.NoError(t, err)

	_, err = db.Exec("SET search_path TO public")
	assert.NoError(t, err)

	// validate DB connection / schema loaded correctly
	_, err = db.Exec("SELECT * FROM users")
	assert.NoError(t, err)

	f(db)
}

func withApp(t *testing.T, f func(session mockSession, a *app.App)) {
	withDB(t, func(db *sql.DB) {
		queries := sqlc.New(db)
		repo := repository.New(queries, db)
		session := newMockSession()

		a := app.New(noLogger{}, &repo, session, template.New())
		f(session, a)
	})
}

func withHandler(t *testing.T, f func(session mockSession, a *app.App, h *handler.Handler)) {
	withApp(t, func(s mockSession, a *app.App) { f(s, a, handler.New(a)) })
}
