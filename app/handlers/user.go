package handlers

import (
	"log/slog"
	"net/http"
)

const currentUserID = "currentUserID"

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
