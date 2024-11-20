package handler_test

import (
	"database/sql"
	"net/http"
	"os"
	"testing"

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

func recreateDB(t *testing.T) {
	db, err := sql.Open("postgres", "user=tracks_test password=1234567 dbname=postgres sslmode=disable")
	assert.NoError(t, err)

	defer db.Close()

	// drop database if exists
	_, err = db.Exec("DROP DATABASE IF EXISTS tracks_test")
	assert.NoError(t, err)

	// create database
	_, err = db.Exec("CREATE DATABASE tracks_test")
	assert.NoError(t, err)
}

func openDB(t *testing.T) *sql.DB {
	// drop and create
	recreateDB(t)

	db, err := sql.Open("postgres", "user=tracks_test password=1234567 dbname=tracks_test sslmode=disable")
	assert.NoError(t, err)

	// load schema file
	schema, err := os.ReadFile("db/schema.sql")
	assert.NoError(t, err)
	_, err = db.Exec(string(schema))
	assert.NoError(t, err)

	return db
}
