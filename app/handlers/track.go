package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/paulsonkoly/tracks/repository"
)

type TemplateData struct {
	CurrentUser *repository.User
}

func (h *Handler) ViewTrack(w http.ResponseWriter, r *http.Request) {
	// mimic loading up current user from DB for now
	app := h.app
	td := TemplateData{}

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

  err := app.Template.Render(w, "track/track.html", td)
	if err != nil {
    app.ServerError(w, "render error", err)
		return
	}
}
