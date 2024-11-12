// Package template serves html from ui/html. The base template is
// ui/html/base.html.
package template

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"math"
	"os"
	"path/filepath"
	"slices"
	"strings"

	humanize "github.com/dustin/go-humanize"
)

// ErrTemplateNotFound indicates that the template file was not found.
var ErrTemplateNotFound = errors.New("template not found")

// Template stores the precomplied page templates.
type Template struct {
	cache map[string]*template.Template
}

var funcMap = template.FuncMap{
	"bytes":    bytes,
	"time":     humanize.Time,
	"distance": distance,
}

func bytes(u int64) string { return humanize.Bytes(uint64(u)) }

func distance(l float64) string {
	if l > 1000 {
		l = math.Round(l/100) / 10
		return fmt.Sprintf("%v km", l)
	}
	return fmt.Sprintf("%v m", math.Round(l))
}

// New loads and precompiles all page templates.
func New() Template {
	cache := make(map[string]*template.Template)

	const (
		htmlPath     = "ui/html/"
		basePath     = "ui/html/base.html"
		partialsPath = "ui/html/partials/"
	)

	partials := []string{}

	// load partials for all templates
	for _, partial := range must(os.ReadDir(partialsPath)) {
		name := partial.Name()
		filename := filepath.Join(partialsPath, name)

		partials = append(partials, filename)
	}

	// per resource directory in html path
	for _, resource := range must(os.ReadDir(htmlPath)) {
		if !resource.IsDir() {
			continue
		}

		// per page html in resources
		resourcePath := filepath.Join(htmlPath, resource.Name())
		for _, page := range must(os.ReadDir(resourcePath)) {
			if !page.Type().IsRegular() || !strings.HasSuffix(page.Name(), ".html") {
				continue
			}

			name := page.Name()
			pagePath := filepath.Join(resourcePath, name)
			// key is not using os dependent dir separator so the handlers don't have to
			key := resource.Name() + "/" + name

			paths := slices.Concat([]string{basePath, pagePath}, partials)
			t := template.New("base.html").Funcs(funcMap)
			cache[key] = must(t.ParseFiles(paths...))
		}
	}

	return Template{cache: cache}
}

// Render renders the template from ui/html/<resource>/page.html. name is the
// path name with ui/html/ removed. data is template specific dynamic data for
// template content.
func (t Template) Render(w io.Writer, name string, data any) error {
	tmpl, ok := t.cache[name]
	if !ok {
		return ErrTemplateNotFound
	}

	return tmpl.ExecuteTemplate(w, "base.html", data)
}

func must[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}

	return value
}
