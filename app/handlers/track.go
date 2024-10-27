package handlers

import (
	"net/http"

	"github.com/paulsonkoly/tracks/repository"
)

type TemplateData struct {
	CurrentUser *repository.User
}

func (h *Handler) ViewTrack(w http.ResponseWriter, r *http.Request) {
	app := h.app

  err := app.Render(w, "track/track.html", app.BaseTemplate(r))
	if err != nil {
    app.ServerError(w, "render error", err)
		return
	}
}
