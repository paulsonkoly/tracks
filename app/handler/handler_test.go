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

type noSession struct{}

func (noSession) Get(_ context.Context, _ string) any        { return nil }
func (noSession) Put(_ context.Context, _ string, _ any)     {}
func (noSession) Remove(_ context.Context, _ string)         {}
func (noSession) Pop(_ context.Context, _ string) any        { return nil }
func (noSession) RenewToken(_ context.Context) error         { return nil }
func (noSession) LoadAndSave(next http.Handler) http.Handler { return next }

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

func withApp(t *testing.T, f func(a *app.App)) {
	withDB(t, func(db *sql.DB) {
		queries := sqlc.New(db)
		repo := repository.New(queries, db)

		a := app.New(noLogger{}, &repo, noSession{}, template.New())
		f(a)
	})
}

func withHandler(t *testing.T, f func(h *handler.Handler)) {
	withApp(t, func(a *app.App) {
		h := handler.New(a)
		f(h)
	})
}
