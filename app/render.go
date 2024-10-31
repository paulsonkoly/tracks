package app

import (
	"errors"
	"io"
	"net/http"

	"github.com/justinas/nosurf"
	"github.com/paulsonkoly/tracks/repository"
)

var ErrTemplateNotFound = errors.New("template not found")

var flashKey = "flash"

type renderData struct {
	Users       []repository.User
	Form        any
	Flash       Flash
	CurrentUser *repository.User
	CSRFToken   string
}

// nolint:revive
// exporting function returning struct with non-exported fields. This is
// intentional here so the handlers can only construct renderData with
// CurrentUser & CSRFToken etc.
func (a *App) BaseTemplate(r *http.Request) renderData {
	td := renderData{}

	user := a.CurrentUser(r.Context())

	td.CurrentUser = user
	td.CSRFToken = nosurf.Token(r)
	if flash, ok := a.sm.Pop(r.Context(), flashKey).(Flash); ok {
		td.Flash = flash
	}

	return td
}

func (r renderData) WithUsers(users []repository.User) renderData {
	r.Users = users
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
	tmpl, ok := a.Template.Get(name)
	if !ok {
		return ErrTemplateNotFound
	}

	return tmpl.Execute(w, data)
}
