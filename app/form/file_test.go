package form_test

import (
	"testing"

	"github.com/paulsonkoly/tracks/app/form"
)

type testDatum struct {
	name     string
	fileName string
	valid    bool
}

var testData = [...]testDatum{
	{"valid", "example.gpx", true},
	{"empty", "", false},
	{"special chars", "[UK] example (5).gpx", true},
	{"invalid chars", "!/what is it.gpx", false},
	{"invalid extension", "example.gpx.txt", false},
	{"no extension", "example", false},
}

func TestFileValidate(t *testing.T) {
	for _, test := range testData {
		t.Run(test.name, func(t *testing.T) {
			form := form.File{Filename: test.fileName}

			if form.Validate() != test.valid {
				t.Errorf("File{Filename: \"%s\"}.Validate() = %v, want %v", form.Filename, form.Validate(), test.valid)
			}
		})
	}
}
