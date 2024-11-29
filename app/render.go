package app

import (
	"io"
	"net/http"

	"github.com/justinas/nosurf"
	"github.com/paulsonkoly/tracks/repository"
)

// SKFlash is the session key for storing flash.
const SKFlash = "flash"

type renderData struct {
	Users       []repository.User
	GPXFiles    []repository.GPXFile
	Track       repository.Track
	Tracks      []repository.Track
	Collection  repository.Collection
	Collections []repository.Collection
	Form        any
	Flash       Flash
	CurrentUser *repository.User
	CSRFToken   string
}

func (a *App) BaseTemplate(r *http.Request) renderData { // nolint:revive
	// exporting function returning struct with non-exported fields. This is
	// intentional here so the handlers can only construct renderData with
	// CurrentUser & CSRFToken etc.
	td := renderData{}

	user := a.CurrentUser(r.Context())

	td.CurrentUser = user
	td.CSRFToken = nosurf.Token(r)
	if flash, ok := a.sm.Pop(r.Context(), SKFlash).(Flash); ok {
		td.Flash = flash
	}

	cs := r.Context().Value(Collections)
	if cs, ok := cs.([]repository.Collection); ok {
		td.Collections = cs
	}

	return td
}

func (r renderData) WithUsers(users []repository.User) renderData {
	r.Users = users
	return r
}

func (r renderData) WithGPXFiles(files []repository.GPXFile) renderData {
	r.GPXFiles = files
	return r
}

func (r renderData) WithTrack(track repository.Track) renderData {
	r.Track = track
	return r
}

func (r renderData) WithTracks(tracks []repository.Track) renderData {
	r.Tracks = tracks
	return r
}

func (r renderData) WithCollection(collection repository.Collection) renderData {
	r.Collection = collection
	return r
}

func (r renderData) WithForm(form any) renderData {
	r.Form = form
	return r
}

// Render renders the template from ui/html/<resource>/page.html. name is the
// path name with ui/html/ removed. renderData can be obtained by calling
// BaseTemplate().
func (a *App) Render(w io.Writer, name string, data renderData) error {
	return a.template.Render(w, name, data)
}
