package handlers

import (
	"net/http"

	"github.com/paulsonkoly/tracks/app/template"
	"github.com/paulsonkoly/tracks/repository"
)

type TemplateData struct {
	CurrentUser *repository.User
}

func (h *Handler) ViewTrack(w http.ResponseWriter, r *http.Request) {
	app := h.app
	td := template.Data{}

	user := app.CurrentUser(r.Context())

  td.CurrentUser = user

  err := app.Template.Render(w, "track/track.html", td)
	if err != nil {
    app.ServerError(w, "render error", err)
		return
	}
}
