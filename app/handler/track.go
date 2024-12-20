package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"unicode/utf8"

	"github.com/paulsonkoly/tracks/repository"
	"github.com/timewasted/go-accept-headers"
)

// Tracks renders the tracks list page.
func (h *Handler) Tracks(w http.ResponseWriter, r *http.Request) {
	a := h.app

	name := r.URL.Query().Get("name")
	if name != "" && utf8.RuneCountInString(name) < 3 {
		a.ClientError(w, errors.New("Track name must be at least 3 characters"), http.StatusBadRequest)
		return
	}

	var (
		tracks []repository.Track
		err    error
	)
	if name == "" {
		tracks, err = a.Q(r.Context()).GetTracks()
	} else {
		tracks, err = a.Q(r.Context()).GetMatchingTracks(name)
	}

	if err != nil {
		a.ServerError(w, err)
		return
	}

	hdr := accept.Parse(r.Header.Get("Accept"))
	switch {
	case hdr.Accepts("text/html"):
		if err := a.Render(w, "track/tracks.html", a.BaseTemplate(r).WithTracks(tracks)); err != nil {
			a.ServerError(w, err)
		}
	case hdr.Accepts("application/json"):
		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(map[string][]repository.Track{"tracks": tracks}); err != nil {
			a.ServerError(w, err)
		}
	default:
		a := h.app
		a.ClientError(w, errors.New("cannot list tracks"), http.StatusUnsupportedMediaType)
	}
}

// ViewTrack renders the track page for a specific id.
func (h *Handler) ViewTrack(w http.ResponseWriter, r *http.Request) {
	a := h.app

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		a.ClientError(w, err, http.StatusBadRequest)
		return
	}

	track, err := a.Q(r.Context()).GetTrack(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			a.ClientError(w, err, http.StatusNotFound)
			return
		}
		a.ServerError(w, err)
		return
	}

	if err := a.Render(w, "track/track.html", a.BaseTemplate(r).WithTrack(track)); err != nil {
		a.ServerError(w, err)
	}
}

// ListTrackPoints returns a json array of segments of points for a track.
func (h *Handler) ListTrackPoints(w http.ResponseWriter, r *http.Request) {
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
