package form

import (
	"context"
	"regexp"
)

type GPXFile struct {
	Filename string
	errors
}

type GPXFileUnique interface {
	Unique(ctx context.Context, filename string) (bool, error)
}

var filenameRexp = regexp.MustCompile(`^([a-zA-Z0-9\[\]\(\)\{\}_ ]+).gpx$`)

func (f *GPXFile) Validate(ctx context.Context, uniq GPXFileUnique) (bool, error) {
	if !filenameRexp.MatchString(f.Filename) {
		f.AddFieldError("Filename", "Invalid filename")
	}

	ok, err := uniq.Unique(ctx, f.Filename)
	if err != nil {
		return false, err
	}
	if !ok {
		f.AddFieldError("Filename", "File name already exist.")
	}

	return f.valid(), nil
}
