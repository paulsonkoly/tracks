package handler_test

import (
	"net/http"
	"os"
	"testing"

	_ "github.com/lib/pq"

	"github.com/paulsonkoly/tracks/app"
	"github.com/paulsonkoly/tracks/app/handler"
	"github.com/paulsonkoly/tracks/app/template"
	"github.com/paulsonkoly/tracks/repository"
	"github.com/paulsonkoly/tracks/repository/sqlc"
	"github.com/stretchr/testify/assert"
)

func TestViewUserLogin(t *testing.T) {
  // cd in project root
	err := os.Chdir("../..")
	assert.NoError(t, err)

	db := openDB(t)
	defer db.Close()

	queries := sqlc.New(db)
	repo := repository.New(queries, db)

	a := app.New(noLogger{}, &repo, noSession{}, template.New())
	h := handler.New(a)
	if assert.HTTPSuccess(t, h.ViewUserLogin, http.MethodGet, "/user/login", nil) {
		assert.HTTPBodyContains(t, h.ViewUserLogin, http.MethodGet, "/user/login", nil, "Login")
		assert.HTTPBodyContains(t, h.ViewUserLogin, http.MethodGet, "/user/login", nil, "Username")
		assert.HTTPBodyContains(t, h.ViewUserLogin, http.MethodGet, "/user/login", nil, "Password")
	}
}
