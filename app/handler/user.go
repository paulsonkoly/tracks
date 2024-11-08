package handler

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

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

	if err := a.DecodeForm(&loginForm, r); err != nil {
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
	a.LogAction(r.Context(), "user login", slog.Int("id", int(user.ID)))
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

	a.LogAction(r.Context(), "user logout")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ViewUsers renders the user admin page with list of all users.
func (h *Handler) ViewUsers(w http.ResponseWriter, r *http.Request) {
	a := h.app

	users, err := a.Repo(nil).GetUsers(r.Context())
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

	if err := a.DecodeForm(&newUserForm, r); err != nil {
		a.ClientError(w, err, http.StatusBadRequest)
		return
	}

	err := a.WithTx(r.Context(), func(h app.TXHandle) error {

		ok, err := newUserForm.Validate(r.Context(), a.Repo(h))
		if err != nil {
			return err
		}
		if !ok {
			// if any errors
			if err := a.Render(w, "user/new.html", a.BaseTemplate(r).WithForm(newUserForm)); err != nil {
				return err
			}

			return nil
		}

		insert := repository.InsertUserParams{Username: newUserForm.Username}
		hash, err := bcrypt.GenerateFromPassword([]byte(newUserForm.Password), 12)
		if err != nil {
			return err
		}
		insert.HashedPassword = string(hash)

		if _, err := a.Repo(h).InsertUser(r.Context(), insert); err != nil {
			return err
		}

		a.FlashInfo(r.Context(), "User created.")
		a.LogAction(r.Context(), "user created", slog.String("username", insert.Username))
		http.Redirect(w, r, "/users", http.StatusSeeOther)

		return nil
	})
	if err != nil {
		a.ServerError(w, err)
	}
}

// EditUser renders the form for editing a user.
func (h *Handler) EditUser(w http.ResponseWriter, r *http.Request) {
	a := h.app

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		a.ClientError(w, err, http.StatusBadRequest)
		return
	}

	dbUser, err := a.Repo(nil).GetUser(r.Context(), int32(id))
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

	if err := a.DecodeForm(&form, r); err != nil {
		a.ClientError(w, err, http.StatusBadRequest)
		return
	}

	err = a.WithTx(r.Context(), func(h app.TXHandle) error {

		dbUser, err := a.Repo(h).GetUser(r.Context(), int32(id))
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.NotFound(w, r)
			} else {
				return err
			}
			return nil
		}
		form.ID = id

		ok, err := form.ValidateEdit(r.Context(), a.Repo(h))
		if err != nil {
			return err
		}
		if !ok {
			if err := a.Render(w, "user/edit.html", a.BaseTemplate(r).WithForm(form)); err != nil {
				return err
			}
			return nil
		}

		upd := repository.UpdateUserParams{Username: dbUser.Username, HashedPassword: dbUser.HashedPassword, ID: int32(id)}
		if form.Username != "" {
			upd.Username = form.Username
		}
		if form.Password != "" {
			hash, err := bcrypt.GenerateFromPassword([]byte(form.Password), 12)
			if err != nil {
				return err
			}
			upd.HashedPassword = string(hash)
		}

		if err := a.Repo(h).UpdateUser(r.Context(), upd); err != nil {
			return err
		}

		a.FlashInfo(r.Context(), "User updated.")
		a.LogAction(r.Context(), "user updated", slog.String("username", upd.Username), slog.Int("id", id))
		http.Redirect(w, r, "/users", http.StatusSeeOther)

		return nil
	})
	if err != nil {
		a.ServerError(w, err)
	}

}

// DeleteUser deletes a user.
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	a := h.app

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		a.ClientError(w, err, http.StatusBadRequest)
		return
	}

	if err := a.Repo(nil).DeleteUser(r.Context(), int32(id)); err != nil {
		a.ServerError(w, err)
		return
	}

	a.FlashInfo(r.Context(), "User deleted.")
	a.LogAction(r.Context(), "user deleted", slog.Int("id", id))
	http.Redirect(w, r, "/users", http.StatusSeeOther)
}
