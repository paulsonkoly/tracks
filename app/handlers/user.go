package handlers

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"

	fdecoder "github.com/go-playground/form/v4"
	"github.com/paulsonkoly/tracks/app/form"
	"github.com/paulsonkoly/tracks/repository"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) ViewUserLogin(w http.ResponseWriter, r *http.Request) {
	app := h.app

	err := app.Render(w, "user/login.html", app.BaseTemplate(r))
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

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) PostUserLogout(w http.ResponseWriter, r *http.Request) {
	app := h.app

	if err := app.SM.RenewToken(r.Context()); err != nil {
		app.ServerError(w, "session renew token", err)
		return
	}

	app.ClearCurrentUser(r.Context())
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) ViewUsers(w http.ResponseWriter, r *http.Request) {
	app := h.app
	users, err := app.Repo.GetUsers(r.Context())
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		app.ServerError(w, "get users", err)
		return
	}

	err = app.Render(w, "user/users.html", app.BaseTemplate(r).WithUsers(users))
	if err != nil {
		app.ServerError(w, "render error", err)
		return
	}
}

func (h *Handler) NewUser(w http.ResponseWriter, r *http.Request) {
	app := h.app
	newUserForm := form.NewUserForm{}

	err := app.Render(w, "user/new.html", app.BaseTemplate(r).WithForm(newUserForm))
	if err != nil {
		app.ServerError(w, "render error", err)
		return
	}
}

func (h *Handler) PostNewUser(w http.ResponseWriter, r *http.Request) {
	app := h.app

	newUserForm := form.NewUserForm{}

	decoder := fdecoder.NewDecoder()
	err := r.ParseForm()
	if err != nil {
		// client
		app.ServerError(w, "parse form error", err)
		return
	}

	err = decoder.Decode(&newUserForm, r.PostForm)
	if err != nil {
		// client
		app.ServerError(w, "decode form error", err)
		return
	}

	newUserForm.Validate(r.Context(), app.Repo)

	if !newUserForm.Valid() {
		// if any errors

		err = app.Render(w, "user/new.html", app.BaseTemplate(r).WithForm(newUserForm))
		if err != nil {
			app.ServerError(w, "render error", err)
			return
		}

		return
	}

	insert := repository.InsertUserParams{Username: newUserForm.Username}
	hash, err := bcrypt.GenerateFromPassword([]byte(newUserForm.Password), 12)
	if err != nil {
		app.ServerError(w, "bcrypt error", err)
		return
	}
	insert.HashedPassword = string(hash)

	_, err = app.Repo.InsertUser(r.Context(), insert)
	if err != nil {
		app.ServerError(w, "render error", err)
		return
	}

	http.Redirect(w, r, "/users", http.StatusSeeOther)
}
