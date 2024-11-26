package handler_test

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"

	"github.com/paulsonkoly/tracks/app"
	"github.com/paulsonkoly/tracks/app/handler"
	"github.com/stretchr/testify/assert"
)

func TestViewUserLogin(t *testing.T) {
	withHandler(t, func(_ mockSession, _ *app.App, h *handler.Handler) {
		r := httptest.NewRequest(http.MethodGet, "/user/login", nil)
		w := httptest.NewRecorder()

		h.ViewUserLogin(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		rendersLoginPage(t, w)
	})
}

func TestPostUserLoginInvalidCredentials(t *testing.T) {
	withHandler(t, func(session mockSession, _ *app.App, h *handler.Handler) {

		r := httptest.NewRequest(http.MethodPost, "/user/login", strings.NewReader("username=admin&password=adminpass"))
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()

		h.PostUserLogin(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		rendersLoginPage(t, w)
		showsInvalidCredentials(t, w)

		session.DoesntHave(t, app.SKCurrentUserID)
	})
}

func TestPostUserLogin(t *testing.T) {
	withHandler(t, func(session mockSession, a *app.App, h *handler.Handler) {

		hash, err := bcrypt.GenerateFromPassword([]byte("adminpass"), 12)
		assert.NoError(t, err)

		user, err := a.Q(context.Background()).InsertUser("admin", string(hash))
		assert.NoError(t, err)

		r := httptest.NewRequest(http.MethodPost, "/user/login", strings.NewReader("username=admin&password=adminpass"))
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()

		h.PostUserLogin(w, r)

		session.Has(t, app.SKCurrentUserID, user.ID)
		redirectsTo(t, w, "/")
	})
}

func TestPostUserLogout(t *testing.T) {
	withHandler(t, func(session mockSession, a *app.App, h *handler.Handler) {

		hash, err := bcrypt.GenerateFromPassword([]byte("adminpass"), 12)
		assert.NoError(t, err)

		user, err := a.Q(context.Background()).InsertUser("admin", string(hash))
		assert.NoError(t, err)

		session.Put(context.Background(), app.SKCurrentUserID, user.ID)

		r := httptest.NewRequest(http.MethodPost, "/user/logout", nil)
		w := httptest.NewRecorder()

		h.PostUserLogout(w, r)

		session.DoesntHave(t, app.SKCurrentUserID)
		redirectsTo(t, w, "/")
	})
}

func TestViewUsers(t *testing.T) {
	withHandler(t, func(_ mockSession, a *app.App, h *handler.Handler) {
		usernames := [...]string{"Alice", "Bob", "Charlie"}
		for _, username := range usernames {
			_, err := a.Q(context.Background()).InsertUser(username, "")
			assert.NoError(t, err)
		}

		r := httptest.NewRequest(http.MethodGet, "/users", nil)
		w := httptest.NewRecorder()

		h.ViewUsers(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		for _, username := range usernames {
			assert.Contains(t, w.Body.String(), username)
		}
	})
}

func TestNewUser(t *testing.T) {
	withHandler(t, func(_ mockSession, _ *app.App, h *handler.Handler) {
		r := httptest.NewRequest(http.MethodGet, "/user/new", nil)
		w := httptest.NewRecorder()

		h.NewUser(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Create New User")
		assert.Contains(t, w.Body.String(), "Username")
		assert.Contains(t, w.Body.String(), "Password")
		assert.Contains(t, w.Body.String(), "Confirm Password")
	})
}

func TestPostNewUser(t *testing.T) {
	withHandler(t, func(session mockSession, a *app.App, h *handler.Handler) {
		r := httptest.NewRequest(http.MethodPost, "/user/new", strings.NewReader("username=admin&password=adminpass&password_confirm=adminpass"))
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()

		h.PostNewUser(w, r)

		assert.Equal(t, http.StatusSeeOther, w.Code)

		session.HasFlashInfo(t, "User created.")

		user, err := a.Q(context.Background()).GetUserByName("admin")
		assert.NoError(t, err)

		assert.Equal(t, user.Username, "admin")
		assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte("adminpass")))
	})
}

func TestPostNewUserFailure(t *testing.T) {
	withHandler(t, func(_ mockSession, a *app.App, h *handler.Handler) {

		_, err := a.Q(context.Background()).InsertUser("admin", "")
		assert.NoError(t, err)

		r := httptest.NewRequest(http.MethodPost, "/user/new", strings.NewReader("username=admin&password=a&password_confirm=mismatch"))
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()

		h.PostNewUser(w, r)

		assert.Equal(t, http.StatusOK, w.Code)

		assert.Contains(t, w.Body.String(), "Username taken.")
		assert.Contains(t, w.Body.String(), "Password too short.")
		assert.Contains(t, w.Body.String(), "Passwords do not match.")
	})
}

func TestEditUser(t *testing.T) {
	withHandler(t, func(_ mockSession, a *app.App, h *handler.Handler) {
		user, err := a.Q(context.Background()).InsertUser("bob", "")
		assert.NoError(t, err)

		r := httptest.NewRequest(http.MethodGet, "/user/{id}/edit", nil)
		r.SetPathValue("id", strconv.Itoa(user.ID))
		w := httptest.NewRecorder()

		h.EditUser(w, r)

		assert.Equal(t, http.StatusOK, w.Code)

		assert.Contains(t, w.Body.String(), "Edit User")
		assert.Contains(t, w.Body.String(), "Username")
		assert.Contains(t, w.Body.String(), user.Username)
		assert.Contains(t, w.Body.String(), "Password")
		assert.Contains(t, w.Body.String(), "Confirm Password")
	})
}

func TestPostEditUser(t *testing.T) {
	withHandler(t, func(session mockSession, a *app.App, h *handler.Handler) {
		user, err := a.Q(context.Background()).InsertUser("bob", "")
		assert.NoError(t, err)

		r := httptest.NewRequest(http.MethodPost, "/user/{id}/edit", strings.NewReader("username=alice&password=1234567&password_confirm=1234567"))
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		r.SetPathValue("id", strconv.Itoa(user.ID))
		w := httptest.NewRecorder()

		h.PostEditUser(w, r)

		assert.Equal(t, http.StatusSeeOther, w.Code)

		modified, err := a.Q(context.Background()).GetUser(user.ID)
		assert.NoError(t, err)
		assert.Equal(t, "alice", modified.Username)
		assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(modified.HashedPassword), []byte("1234567")))

		session.HasFlashInfo(t, "User updated.")
	})
}

func TestPostEditUserUsernameOnly(t *testing.T) {
	withHandler(t, func(_ mockSession, a *app.App, h *handler.Handler) {
		hash, err := bcrypt.GenerateFromPassword([]byte("somepass"), 12)
		assert.NoError(t, err)

		user, err := a.Q(context.Background()).InsertUser("bob", string(hash))
		assert.NoError(t, err)

		r := httptest.NewRequest(http.MethodPost, "/user/{id}/edit", strings.NewReader("username=alice"))
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		r.SetPathValue("id", strconv.Itoa(user.ID))
		w := httptest.NewRecorder()

		h.PostEditUser(w, r)

		modified, err := a.Q(context.Background()).GetUser(user.ID)
		assert.NoError(t, err)
		assert.Equal(t, "alice", modified.Username)
		assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(modified.HashedPassword), []byte("somepass")))
	})
}

func TestPostEditUserFailure(t *testing.T) {
	withHandler(t, func(_ mockSession, a *app.App, h *handler.Handler) {
		user, err := a.Q(context.Background()).InsertUser("bob", "")
		assert.NoError(t, err)

		_, err = a.Q(context.Background()).InsertUser("alice", "")
		assert.NoError(t, err)

		r := httptest.NewRequest(http.MethodPost, "/user/{id}/edit", strings.NewReader("username=alice&password=1234567&password_confirm=1234567"))
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		r.SetPathValue("id", strconv.Itoa(user.ID))
		w := httptest.NewRecorder()

		h.PostEditUser(w, r)

		assert.Equal(t, http.StatusOK, w.Code)

		modified, err := a.Q(context.Background()).GetUser(user.ID)
		assert.NoError(t, err)
		assert.Equal(t, user, modified)

		assert.Contains(t, w.Body.String(), "Username taken.")
	})
}

func TestDeleteUser(t *testing.T) {
	withHandler(t, func(session mockSession, a *app.App, h *handler.Handler) {
		user, err := a.Q(context.Background()).InsertUser("bob", "")
		assert.NoError(t, err)

		r := httptest.NewRequest(http.MethodPost, "/user/{id}/delete", nil)
		r.SetPathValue("id", strconv.Itoa(user.ID))
		w := httptest.NewRecorder()

		h.DeleteUser(w, r)

		_, err = a.Q(context.Background()).GetUserByName("bob")
		assert.ErrorIs(t, err, sql.ErrNoRows)

		assert.Equal(t, http.StatusSeeOther, w.Code)

		session.HasFlashInfo(t, "User deleted.")
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
