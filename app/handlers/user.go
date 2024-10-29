package handlers

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	fdecoder "github.com/go-playground/form/v4"
	"github.com/paulsonkoly/tracks/app"
	"github.com/paulsonkoly/tracks/app/form"
	"github.com/paulsonkoly/tracks/repository"
	"golang.org/x/crypto/bcrypt"
)

// ViewUserLogin renders the login page.
func (h *Handler) ViewUserLogin(w http.ResponseWriter, r *http.Request) {
	a := h.app

	err := a.Render(w, "user/login.html", a.BaseTemplate(r).WithForm(form.Login{}))
	if err != nil {
		a.ServerError(w, err)
		return
	}
}

// PostUserLogin processes the login form.
func (h *Handler) PostUserLogin(w http.ResponseWriter, r *http.Request) {
	a := h.app

	loginForm := form.Login{}

	decoder := fdecoder.NewDecoder()
	err := r.ParseForm()
	if err != nil {
		a.ClientError(w, err, http.StatusBadRequest)
		return
	}

	err = decoder.Decode(&loginForm, r.PostForm)
	if err != nil {
		a.ClientError(w, err, http.StatusBadRequest)
		return
	}

	user, err := a.AuthenticateUser(r.Context(), loginForm.Username, loginForm.Password)
	if err != nil && errors.Is(err, app.ErrAuthenticationFailed) {
		loginForm.AddError("Invalid credentials.")
		// remove data to prevent form filling in fields from previous submit
		loginForm.Username = ""
		loginForm.Password = ""

		err = a.Render(w, "user/login.html", a.BaseTemplate(r).WithForm(loginForm))
		if err != nil {
			a.ServerError(w, err)
			return
		}

		return
	} else if err != nil {
		a.ServerError(w, err)
		return
	}

	// succesfull login
	a.Logger.Info("user login", slog.String("username", user.Username), slog.Int("id", int(user.ID)))
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// PostUserLogout handles user logout.
func (h *Handler) PostUserLogout(w http.ResponseWriter, r *http.Request) {
	a := h.app

	err := a.ClearCurrentUser(r.Context())
	if err != nil {
		a.ServerError(w, err)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ViewUsers renders the user admin page with list of all users.
func (h *Handler) ViewUsers(w http.ResponseWriter, r *http.Request) {
	a := h.app

	users, err := a.Repo.GetUsers(r.Context())
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		a.ServerError(w, err)
		return
	}

	err = a.Render(w, "user/users.html", a.BaseTemplate(r).WithUsers(users))
	if err != nil {
		a.ServerError(w, err)
		return
	}
}

// NewUser renders a form to create new user.
func (h *Handler) NewUser(w http.ResponseWriter, r *http.Request) {
	a := h.app

	err := a.Render(w, "user/new.html", a.BaseTemplate(r).WithForm(form.User{}))
	if err != nil {
		a.ServerError(w, err)
		return
	}
}

// PostNewUser saves new user.
func (h *Handler) PostNewUser(w http.ResponseWriter, r *http.Request) {
	a := h.app

	newUserForm := form.User{}

	decoder := fdecoder.NewDecoder()
	err := r.ParseForm()
	if err != nil {
		a.ClientError(w, err, http.StatusBadRequest)
		return
	}

	err = decoder.Decode(&newUserForm, r.PostForm)
	if err != nil {
		a.ClientError(w, err, http.StatusBadRequest)
		return
	}

	newUserForm.Validate(r.Context(), a.Repo)

	if !newUserForm.Valid() {
		// if any errors
		err = a.Render(w, "user/new.html", a.BaseTemplate(r).WithForm(newUserForm))
		if err != nil {
			a.ServerError(w, err)
			return
		}

		return
	}

	insert := repository.InsertUserParams{Username: newUserForm.Username}
	hash, err := bcrypt.GenerateFromPassword([]byte(newUserForm.Password), 12)
	if err != nil {
		a.ServerError(w, err)
		return
	}
	insert.HashedPassword = string(hash)

	_, err = a.Repo.InsertUser(r.Context(), insert)
	if err != nil {
		a.ServerError(w, err)
		return
	}

	a.FlashInfo(r.Context(), "User created.")
	http.Redirect(w, r, "/users", http.StatusSeeOther)
}

// EditUser renders the form for editing a user.
func (h *Handler) EditUser(w http.ResponseWriter, r *http.Request) {
	a := h.app

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		a.ClientError(w, err, http.StatusBadRequest)
		return
	}

	dbUser, err := a.Repo.GetUser(r.Context(), int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
		} else {
			a.ServerError(w, err)
		}
		return
	}

	form := form.User{Username: dbUser.Username, ID: id}

	err = a.Render(w, "user/edit.html", a.BaseTemplate(r).WithForm(form))
	if err != nil {
		a.ServerError(w, err)
		return
	}
}

// PostEditUser updates the user with the updated data.
func (h *Handler) PostEditUser(w http.ResponseWriter, r *http.Request) {
	a := h.app

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		a.ClientError(w, err, http.StatusBadRequest)
		return
	}

	form := form.User{}

	decoder := fdecoder.NewDecoder()
	err = r.ParseForm()
	if err != nil {
		a.ClientError(w, err, http.StatusBadRequest)
		return
	}

	err = decoder.Decode(&form, r.PostForm)
	if err != nil {
		a.ClientError(w, err, http.StatusBadRequest)
		return
	}

	form.ValidateEdit()

	dbUser, err := a.Repo.GetUser(r.Context(), int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
		} else {
			a.ServerError(w, err)
		}
		return
	}
	form.ID = id

	if !form.Valid() {
		err = a.Render(w, "user/edit.html", a.BaseTemplate(r).WithForm(form))
		if err != nil {
			a.ServerError(w, err)
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
			a.ServerError(w, err)
			return
		}
		upd.HashedPassword = string(hash)
	}

	err = a.Repo.UpdateUser(r.Context(), upd)
	if err != nil {
		a.ServerError(w, err)
		return
	}

	a.FlashInfo(r.Context(), "User updated.")
	http.Redirect(w, r, "/users", http.StatusSeeOther)
}

// DeleteUser deletes a user.
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	app := h.app

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		app.ClientError(w, err, http.StatusBadRequest)
		return
	}

	err = app.Repo.DeleteUser(r.Context(), int32(id))
	if err != nil {
		app.ServerError(w, err)
		return
	}

	app.FlashInfo(r.Context(), "User deleted.")
	http.Redirect(w, r, "/users", http.StatusSeeOther)
}
