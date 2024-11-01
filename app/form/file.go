package form

import "regexp"

type File struct {
	Filename string
	errors
}

var filenameRexp = regexp.MustCompile(`^([a-zA-Z0-9\[\]\(\)\{\}_ ]+).gpx$`)

func (f *File) Validate() bool {
	if !filenameRexp.MatchString(f.Filename) {
		f.AddFieldError("Filename", "Invalid filename")
	}
	return f.valid()
}
