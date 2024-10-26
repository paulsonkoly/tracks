package handlers

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"

	"github.com/paulsonkoly/tracks/repository"
)

const currentUserID = "currentUserID"

func (h * Handler)ViewUserLogin(w http.ResponseWriter, _ *http.Request) {
  app := h.app

  err := app.Template.Render(w, "user/login.html", nil)
	if err != nil {
    app.ServerError(w, "template error", err)
    return
	}
}

func (h *Handler) PostUserLogin(w http.ResponseWriter, r *http.Request) {
	app := h.app

	if err := app.SM.RenewToken(r.Context()); err != nil {
		app.ServerError(w, "session renew token", err)
		return
	}

	user, err := app.AuthenticateUser(r.Context(), "Paul", "123456")
	if err != nil {
		app.ServerError(w, "authenticate user", err)
		return
	}

	if user == nil {
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	app.Logger.Info("user login", slog.String("username", user.Username), slog.Int("id", int(user.ID)))

	app.SM.Put(r.Context(), currentUserID, user.ID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) PostUserLogout(w http.ResponseWriter, r *http.Request) {
	app := h.app

	if err := app.SM.RenewToken(r.Context()); err != nil {
		app.ServerError(w, "session renew token", err)
		return
	}

	app.SM.Remove(r.Context(), currentUserID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

type TemplateDataUsers struct {
	CurrentUser *repository.User
	Users       []repository.User
}

func (h *Handler) ViewUsers(w http.ResponseWriter, r *http.Request) {
	app := h.app
	users, err := app.Repo.GetUsers(r.Context())
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		app.ServerError(w, "get users", err)
		return
	}

	td := TemplateDataUsers{}

	if app.SM.Exists(r.Context(), currentUserID) {
		uid := app.SM.GetInt32(r.Context(), currentUserID)

		user, err := app.Repo.GetUser(r.Context(), uid)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
      app.ServerError(w, "render error", err)
			return
		}

		td.CurrentUser = &user
	}

	td.Users = users

  err = app.Template.Render(w, "user/users.html", td)
	if err != nil {
    app.ServerError(w, "render error", err)
		return
	}
}

type TemplateDataUser struct {
	CurrentUser *repository.User
}

func (h *Handler) NewUser(w http.ResponseWriter, r *http.Request) {
	app := h.app
	td := TemplateDataUser{}

	if app.SM.Exists(r.Context(), currentUserID) {
		uid := app.SM.GetInt32(r.Context(), currentUserID)

		user, err := app.Repo.GetUser(r.Context(), uid)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			app.Logger.Error("current user", "error", err)
			http.Error(w, "user error", http.StatusInternalServerError)
			return
		}

		td.CurrentUser = &user
	}

  err := app.Template.Render(w, "user/new.html", td)
	if err != nil {
    app.ServerError(w, "render error", err)
		return
	}
}
