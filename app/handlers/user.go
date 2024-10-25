package handlers

import (
	"log/slog"
	"net/http"
)

const currentUserID = "currentUserID"

func (h *Handler) PostUserLogin(w http.ResponseWriter, r *http.Request) {
  app := h.app
  err := app.SM.RenewToken(r.Context())
  if err!= nil {
    app.Logger.Error("RenewToken", "error", err)
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  app.Logger.Info("user login", slog.Int("uid", 1))

  app.SM.Put(r.Context(), currentUserID, 1)
  http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) PostUserLogout(w http.ResponseWriter, r *http.Request) {

  app := h.app
  err := app.SM.RenewToken(r.Context())
  if err!= nil {
    app.Logger.Error("RenewToken", "error", err)
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  app.SM.Remove(r.Context(), currentUserID)
  http.Redirect(w, r, "/", http.StatusSeeOther)
}
