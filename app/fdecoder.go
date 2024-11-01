package app

import (
	"net/http"

	"github.com/go-playground/form/v4"
)

type fDecoder struct {
	decoder *form.Decoder
}

func newDecoder() fDecoder {
	return fDecoder{decoder: form.NewDecoder()}
}

func (f fDecoder) decode(values any, r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = f.decoder.Decode(values, r.PostForm)
	if err != nil {
		return err
	}

	return nil
}

// DecodeForm decodes the POST form from r into the form structure pointed by values.
func (a *App) DecodeForm(values any, r *http.Request) error {
	return a.decoder.decode(values, r)
}
