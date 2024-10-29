package handlers

import (
	"net/http"

	"github.com/paulsonkoly/tracks/repository"
)

type TemplateData struct {
	CurrentUser *repository.User
}

func (h *Handler) ViewTrack(w http.ResponseWriter, r *http.Request) {
	a := h.app

	err := a.Render(w, "track/track.html", a.BaseTemplate(r))
	if err != nil {
		a.ServerError(w, err)
		return
	}
}
