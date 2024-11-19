package handler

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/paulsonkoly/tracks/app"
	"github.com/paulsonkoly/tracks/app/form"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
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
	a.LogAction(r.Context(), "user login", slog.Int("id", user.ID))
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

	users, err := a.Q(r.Context()).GetUsers()
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

	err := a.WithTx(r.Context(), func(ctx context.Context) error {

		ok, err := newUserForm.Validate(a.Q(ctx))
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

		hash, err := bcrypt.GenerateFromPassword([]byte(newUserForm.Password), 12)
		if err != nil {
			return err
		}

		if _, err := a.Q(ctx).InsertUser(newUserForm.Username, string(hash)); err != nil {
			return err
		}

		a.FlashInfo(ctx, "User created.")
		a.LogAction(ctx, "user created", slog.String("username", newUserForm.Username))
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

	dbUser, err := a.Q(r.Context()).GetUser(id)
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

	err = a.WithTx(r.Context(), func(ctx context.Context) error {

		dbUser, err := a.Q(ctx).GetUser(id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.NotFound(w, r)
			} else {
				return err
			}
			return nil
		}
		form.ID = id

		ok, err := form.ValidateEdit(a.Q(ctx))
		if err != nil {
			return err
		}
		if !ok {
			if err := a.Render(w, "user/edit.html", a.BaseTemplate(r).WithForm(form)); err != nil {
				return err
			}
			return nil
		}

		if form.Username != "" {
			dbUser.Username = form.Username
		}
		if form.Password != "" {
			hash, err := bcrypt.GenerateFromPassword([]byte(form.Password), 12)
			if err != nil {
				return err
			}
			dbUser.HashedPassword = string(hash)
		}

		if err := a.Q(ctx).UpdateUser(id, dbUser.Username, dbUser.HashedPassword); err != nil {
			return err
		}

		a.FlashInfo(ctx, "User updated.")
		a.LogAction(ctx, "user updated", slog.String("username", dbUser.Username), slog.Int("id", id))
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

	if err := a.Q(r.Context()).DeleteUser(id); err != nil {
		a.ServerError(w, err)
		return
	}

	a.FlashInfo(r.Context(), "User deleted.")
	a.LogAction(r.Context(), "user deleted", slog.Int("id", id))
	http.Redirect(w, r, "/users", http.StatusSeeOther)
}
