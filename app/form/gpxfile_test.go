package form_test

import (
	"context"
	"testing"

	"github.com/paulsonkoly/tracks/app/form"
)

type gpxFileTestDatum struct {
	name     string
	fileName string
	valid    bool
}

var gpxFileTestData = [...]gpxFileTestDatum{
	{"valid", "example.gpx", true},
	{"empty", "", false},
	{"special chars", "[UK] example (5).gpx", true},
	{"invalid chars", "!/what is it.gpx", false},
	{"invalid extension", "example.gpx.txt", false},
	{"no extension", "example", false},
}

type gpxFileTestUnique struct{}

func (u gpxFileTestUnique) GPXFileUnique(_ context.Context, _ string) (bool, error) { return true, nil }

func TestFileValidate(t *testing.T) {
	for _, test := range gpxFileTestData {
		t.Run(test.name, func(t *testing.T) {
			form := form.GPXFile{Filename: test.fileName}

			result, err := form.Validate(context.Background(), gpxFileTestUnique{})
			if err != nil {
				t.Errorf("File{Filename: %q}.Validate() returned error: %v", form.Filename, err)
			}

			if result != test.valid {
				t.Errorf("File{Filename: %q}.Validate() = %v, want %v", form.Filename, result, test.valid)
			}
		})
	}
}
