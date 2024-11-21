package handler_test

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"

	"github.com/paulsonkoly/tracks/app"
	"github.com/paulsonkoly/tracks/app/handler"
	"github.com/paulsonkoly/tracks/app/template"
	"github.com/paulsonkoly/tracks/repository"
	"github.com/paulsonkoly/tracks/repository/sqlc"
	"github.com/stretchr/testify/assert"
)

func TestViewUserLogin(t *testing.T) {
	withHandler(t, func(h *handler.Handler) {
		r := httptest.NewRequest(http.MethodGet, "/user/login", nil)
		w := httptest.NewRecorder()

		h.ViewUserLogin(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		rendersLoginPage(t, w)
	})
}

func TestPostUserLoginInvalidCredentials(t *testing.T) {
	withHandler(t, func(h *handler.Handler) {
		r := httptest.NewRequest(http.MethodPost, "/user/login", strings.NewReader("Username=admin&Password=adminpass"))
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()

		h.PostUserLogin(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		rendersLoginPage(t, w)
		showsInvalidCredentials(t, w)
	})
}

func TestPostUserLogin(t *testing.T) {
	withDB(t, func(db *sql.DB) {

		hash, err := bcrypt.GenerateFromPassword([]byte("adminpass"), 12)
		assert.NoError(t, err)

		_, err = db.Exec("INSERT INTO users (username, hashed_password, created_at) VALUES ('admin', $1, NOW())", hash)
		assert.NoError(t, err)

		queries := sqlc.New(db)
		repo := repository.New(queries, db)

		a := app.New(noLogger{}, &repo, noSession{}, template.New())
		h := handler.New(a)

		r := httptest.NewRequest(http.MethodPost, "/user/login", strings.NewReader("Username=admin&Password=adminpass"))
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()

		h.PostUserLogin(w, r)

		redirectsTo(t, w, "/")
	})
}

func rendersLoginPage(t *testing.T, w *httptest.ResponseRecorder) {
	body := w.Body.String()
	assert.Contains(t, body, "Login")
	assert.Contains(t, body, "Username")
	assert.Contains(t, body, "Password")
}

func showsInvalidCredentials(t *testing.T, w *httptest.ResponseRecorder) {
	body := w.Body.String()
	assert.Contains(t, body, "Invalid credentials.")
}

func redirectsTo(t *testing.T, w *httptest.ResponseRecorder, url string) {
	assert.Equal(t, http.StatusSeeOther, w.Code)
	assert.Equal(t, url, w.Result().Header.Get("Location"))
}
