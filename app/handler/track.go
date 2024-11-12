package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func (h *Handler) ViewTracks(w http.ResponseWriter, r *http.Request) {
	a := h.app

	tracks, err := a.Repo(nil).GetTracks(r.Context())
	if err != nil {
		a.ServerError(w, err)
		return
	}

	if err := a.Render(w, "track/tracks.html", a.BaseTemplate(r).WithTracks(tracks)); err != nil {
		a.ServerError(w, err)
	}
}

func (h *Handler) ViewTrack(w http.ResponseWriter, r *http.Request) {
	a := h.app

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		a.ClientError(w, err, http.StatusBadRequest)
		return
	}

	track, err := a.Repo(nil).GetTrack(r.Context(), int32(id))
	if err != nil {
		a.ServerError(w, err)
		return
	}

	if err := a.Render(w, "track/track.html", a.BaseTemplate(r).WithTrack(track)); err != nil {
		a.ServerError(w, err)
	}
}

func (h *Handler) ViewTrackPoints(w http.ResponseWriter, r *http.Request) {
	a := h.app

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		a.ClientError(w, err, http.StatusBadRequest)
		return
	}

	points, err := a.Repo(nil).GetTrackSegmentPoints(r.Context(), int32(id))
	if err != nil {
		a.ServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(points)
	if err != nil {
		panic(err)
	}
}
