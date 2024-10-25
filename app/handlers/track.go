package handlers

import (
	"net/http"
	"text/template"

	"github.com/paulsonkoly/tracks/repository"
)

type TemplateData struct {
  CurrentUser *repository.User
}

func (h *Handler)ViewTrack(w http.ResponseWriter, r *http.Request) {
  // mimic loading up current user from DB for now
  app := h.app
  td := TemplateData{}

  if app.SM.Exists(r.Context(), currentUserID) {
    uid := app.SM.GetInt32(r.Context(), currentUserID)

    td.CurrentUser = &repository.User{ID: uid, Username: "Paul"}
  }

	t, err := template.ParseFiles("ui/html/base.html", "ui/html/partials/navbar.html", "ui/html/track/track.html")
	if err != nil {
    app.Logger.Error("template err", "error", err)
    http.Error(w, "template error", http.StatusInternalServerError)
    return
	}

	err = t.Execute(w, td)
	if err != nil {
    app.Logger.Error("template err", "error", err)
    http.Error(w, "template error", http.StatusInternalServerError)
    return
	}
}
