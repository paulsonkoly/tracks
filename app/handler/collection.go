package handler

import (
	"net/http"

	"github.com/paulsonkoly/tracks/app/form"
)

// NewCollection renders a form to create a new track collection.
func (h *Handler) NewCollection(w http.ResponseWriter, r *http.Request) {
	a := h.app

	err := a.Render(w, "collection/new.html", a.BaseTemplate(r).WithForm(form.Collection{}))
	if err != nil {
		a.ServerError(w, err)
		return
	}
}
