package template

import (
	"html/template"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type Cache struct {
	cache map[string]*template.Template
}

func NewCache() *Cache {
	cache := make(map[string]*template.Template)

	var (
		htmlPath     = "ui/html/"
		basePath     = "ui/html/base.html"
		partialsPath = "ui/html/partials/"

		partials = []string{}
	)

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
			cache[key] = must(template.ParseFiles(paths...))
		}
	}

	return &Cache{
		cache: cache,
	}
}

func (t Cache) Get(name string) (*template.Template, bool) {
	tmpl, ok := t.cache[name]
	return tmpl, ok
}

func must[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}

	return value
}
