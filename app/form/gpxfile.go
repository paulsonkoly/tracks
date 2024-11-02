package form

import (
	"context"
	"database/sql"
	errs "errors"
	"regexp"

	"github.com/paulsonkoly/tracks/repository"
)

type GPXFile struct {
	Filename string
	errors
}

var filenameRexp = regexp.MustCompile(`^([a-zA-Z0-9\[\]\(\)\{\}_ ]+).gpx$`)

func (f *GPXFile) Validate(ctx context.Context, r *repository.Queries) (bool, error) {
	if !filenameRexp.MatchString(f.Filename) {
		f.AddFieldError("Filename", "Invalid filename")
	}
	_, err := r.GetGPXFileByFilename(ctx, f.Filename)
	if err != nil {
		if !errs.Is(err, sql.ErrNoRows) {
			return false, err
		}
	} else {
		f.AddFieldError("Filename", "File name already exist.")
	}

	return f.valid(), nil
}
