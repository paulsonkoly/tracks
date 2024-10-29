package handlers

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	fdecoder "github.com/go-playground/form/v4"
	"github.com/paulsonkoly/tracks/app/form"
	"github.com/paulsonkoly/tracks/repository"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) ViewUserLogin(w http.ResponseWriter, r *http.Request) {
	app := h.app

	loginForm := form.LoginForm{}

	err := app.Render(w, "user/login.html", app.BaseTemplate(r).WithForm(loginForm))
	if err != nil {
		app.ServerError(w, "template error", err)
		return
	}
}

func (h *Handler) PostUserLogin(w http.ResponseWriter, r *http.Request) {
	app := h.app

	loginForm := form.LoginForm{}

	decoder := fdecoder.NewDecoder()
	err := r.ParseForm()
	if err != nil {
		// client
		app.ServerError(w, "parse form error", err)
		return
	}

	err = decoder.Decode(&loginForm, r.PostForm)
	if err != nil {
		// client
		app.ServerError(w, "decode form error", err)
		return
	}

	user, err := app.AuthenticateUser(r.Context(), loginForm.Username, loginForm.Password)
	if err != nil {
		app.ServerError(w, "authenticate user error", err)
		return
	}

	if user == nil {
		loginForm.AddError("Invalid credentials.")
		// remove data to prevent form filling in fields from previous submit
		loginForm.Username = ""
		loginForm.Password = ""

		err = app.Render(w, "user/login.html", app.BaseTemplate(r).WithForm(loginForm))
		if err != nil {
			app.ServerError(w, "render error", err)
			return
		}

		return
	}

	// succesfull login
	app.Logger.Info("user login", slog.String("username", user.Username), slog.Int("id", int(user.ID)))
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) PostUserLogout(w http.ResponseWriter, r *http.Request) {
	app := h.app

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
	newUserForm := form.User{}

	err := app.Render(w, "user/new.html", app.BaseTemplate(r).WithForm(newUserForm))
	if err != nil {
		app.ServerError(w, "render error", err)
		return
	}
}

func (h *Handler) PostNewUser(w http.ResponseWriter, r *http.Request) {
	app := h.app

	newUserForm := form.User{}

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

	app.FlashInfo(r.Context(), "User created.")
	http.Redirect(w, r, "/users", http.StatusSeeOther)
}

func (h *Handler) EditUser(w http.ResponseWriter, r *http.Request) {
	app := h.app

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		app.ServerError(w, "decoding id", err)
		return
	}

	dbUser, err := app.Repo.GetUser(r.Context(), int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
		} else {
			app.ServerError(w, "getting user", err)
		}
		return
	}

	form := form.User{Username: dbUser.Username, ID: id}

	err = app.Render(w, "user/edit.html", app.BaseTemplate(r).WithForm(form))
	if err != nil {
		app.ServerError(w, "rendering template", err)
		return
	}
}

func (h *Handler) PostEditUser(w http.ResponseWriter, r *http.Request) {
	app := h.app

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		app.ServerError(w, "decoding id", err)
		return
	}

	form := form.User{}

	decoder := fdecoder.NewDecoder()
	err = r.ParseForm()
	if err != nil {
		// client
		app.ServerError(w, "parse form error", err)
		return
	}

	err = decoder.Decode(&form, r.PostForm)
	if err != nil {
		// client
		app.ServerError(w, "decode form error", err)
		return
	}

	form.ValidateEdit()

	dbUser, err := app.Repo.GetUser(r.Context(), int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
		} else {
			app.ServerError(w, "getting user", err)
		}
		return
	}
	form.ID = id

	if !form.Valid() {
		err = app.Render(w, "user/edit.html", app.BaseTemplate(r).WithForm(form))
		if err != nil {
			app.ServerError(w, "render error", err)
		}
		return
	}

	upd := repository.UpdateUserParams{Username: dbUser.Username, HashedPassword: dbUser.HashedPassword, ID: int32(id)}
	if form.Username != "" {
		upd.Username = form.Username
	}
	if form.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(form.Password), 12)
		if err != nil {
			app.ServerError(w, "generate hash from password", err)
			return
		}
		upd.HashedPassword = string(hash)
	}

	err = app.Repo.UpdateUser(r.Context(), upd)
	if err != nil {
		app.ServerError(w, "updating user", err)
		return
	}

	app.FlashInfo(r.Context(), "User updated.")
	http.Redirect(w, r, "/users", http.StatusSeeOther)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	app := h.app

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		app.ServerError(w, "decoding id", err)
		return
	}

	err = app.Repo.DeleteUser(r.Context(), int32(id))
	if err != nil {
		app.ServerError(w, "deleting user", err)
		return
	}

	app.FlashInfo(r.Context(), "User deleted.")
	http.Redirect(w, r, "/users", http.StatusSeeOther)
}
