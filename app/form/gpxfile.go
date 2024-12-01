package form

import (
	"regexp"
)

// GPXFile validates the data associated with a GPX file upload.
type GPXFile struct {
	Filename string `form:"filename"`
	errors   `form:"-"`
}

// GPXFileUniqueChecker checks if the file has already been uploaded.
type GPXFilePresenceChecker interface {
	GPXFilenameExists(filename string) (bool, error)
}

var filenameRexp = regexp.MustCompile(`^([-+%a-zA-Z0-9\[\]\(\)\{\}_ .]+).gpx$`)

// Validate validates the data associated with a GPX file upload.
func (f *GPXFile) Validate(check GPXFilePresenceChecker) (bool, error) {
	if !filenameRexp.MatchString(f.Filename) {
		f.AddFieldError("Filename", "Invalid filename. Allowed characters: -, +, %, a-Z, 0-9, [, ], (, ), {, }, _, ., \" \".")
	}

	ok, err := check.GPXFilenameExists(f.Filename)
	if err != nil {
		return false, err
	}
	if ok {
		f.AddFieldError("Filename", "File name already exist.")
	}

	return f.valid(), nil
}
