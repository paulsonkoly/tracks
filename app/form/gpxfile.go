package form

import (
	"context"
	"regexp"
)

// GPXFile validates the data associated with a GPX file upload.
type GPXFile struct {
	Filename string
	errors
}

// GPXFileUniqueChecker checks if the file has already been uploaded.
type GPXFileUniqueChecker interface {
	GPXFileUnique(ctx context.Context, filename string) (bool, error)
}

var filenameRexp = regexp.MustCompile(`^([a-zA-Z0-9\[\]\(\)\{\}_ ]+).gpx$`)

// Validate validates the data associated with a GPX file upload.
func (f *GPXFile) Validate(ctx context.Context, uniq GPXFileUniqueChecker) (bool, error) {
	if !filenameRexp.MatchString(f.Filename) {
		f.AddFieldError("Filename", "Invalid filename")
	}

	ok, err := uniq.GPXFileUnique(ctx, f.Filename)
	if err != nil {
		return false, err
	}
	if !ok {
		f.AddFieldError("Filename", "File name already exist.")
	}

	return f.valid(), nil
}
