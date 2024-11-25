package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"unicode/utf8"

	"github.com/paulsonkoly/tracks/repository"
)

func (h *Handler) ViewTracks(w http.ResponseWriter, r *http.Request) {
	a := h.app

	tracks, err := a.Q(r.Context()).GetTracks()
	if err != nil {
		a.ServerError(w, err)
		return
	}

	if err := a.Render(w, "track/tracks.html", a.BaseTemplate(r).WithTracks(tracks)); err != nil {
		a.ServerError(w, err)
	}
}

// ListTracks is a json end point to list track names matching a string fragment from ?name=<fragment>.
func (h *Handler) ListTracks(w http.ResponseWriter, r *http.Request) {
	a := h.app

	name := r.URL.Query().Get("name")
	if utf8.RuneCountInString(name) < 3 {
		a.ClientError(w, errors.New("Track name must be at least 3 characters"), http.StatusBadRequest)
		return
	}

	tracks, err := a.Q(r.Context()).GetMatchingTracks(name)
	if err != nil {
		a.ServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(map[string][]repository.Track{"tracks": tracks})
	if err != nil {
		a.ServerError(w, err)
		return
	}
}

func (h *Handler) ViewTrack(w http.ResponseWriter, r *http.Request) {
	a := h.app

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		a.ClientError(w, err, http.StatusBadRequest)
		return
	}

	track, err := a.Q(r.Context()).GetTrack(id)
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

	points, err := a.Q(r.Context()).GetTrackPoints(id)
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
