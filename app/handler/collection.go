package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

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

func (h *Handler) PostNewCollection(w http.ResponseWriter, r *http.Request) {
	a := h.app

	newCollectionForm := form.Collection{}
	if err := a.DecodeForm(&newCollectionForm, r); err != nil {
		a.ClientError(w, err, http.StatusBadRequest)
	}

	err := a.WithTx(r.Context(), func(ctx context.Context) error {
		ok, err := newCollectionForm.Validate(a.Q(ctx))
		if err != nil {
			return err
		}

		if !ok {
			// Do we want to somehow retain the set of tracks here? We would need to
			// load the tracks in the form for the ui to render, however is it worth
			// it?
			if err := a.Render(w, "collection/new.html", a.BaseTemplate(r).WithForm(newCollectionForm)); err != nil {
				return err
			}

			return nil
		}
		if err := a.Q(ctx).InsertCollection(newCollectionForm.Name, *a.CurrentUser(ctx), newCollectionForm.TrackIDs); err != nil {
			return err
		}

		a.FlashInfo(ctx, "Collection created.")
		a.LogAction(ctx, "collection created", "name", newCollectionForm.Name)

		http.Redirect(w, r, "/", http.StatusSeeOther)

		return nil
	})

	if err != nil {
		a.ServerError(w, err)
	}
}

func (h *Handler) ViewCollection(w http.ResponseWriter, r *http.Request) {
	a := h.app

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		a.ClientError(w, err, http.StatusBadRequest)
		return
	}

	col, err := a.Q(r.Context()).GetCollection(id)
	if err != nil {
		a.ServerError(w, err)
		return
	}

	if err := a.Render(w, "collection/collection.html", a.BaseTemplate(r).WithCollection(col)); err != nil {
		a.ServerError(w, err)
	}
}

func (h *Handler) ListCollectionPoints(w http.ResponseWriter, r *http.Request) {
	a := h.app

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		a.ClientError(w, err, http.StatusBadRequest)
		return
	}

	points, err := a.Q(r.Context()).GetCollectionPoints(id)
	if err != nil {
		a.ServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(points)
	if err != nil {
		a.ServerError(w, err)
		return
	}
}
