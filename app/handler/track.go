package handler

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/paulsonkoly/tracks/app"
	"github.com/paulsonkoly/tracks/app/form"
)

func (h *Handler) ViewTrack(w http.ResponseWriter, r *http.Request) {
	a := h.app

	if err := a.Render(w, "track/track.html", a.BaseTemplate(r)); err != nil {
		a.ServerError(w, err)
	}
}

func (h *Handler) UploadTrack(w http.ResponseWriter, r *http.Request) {
	a := h.app

	form := form.File{}
	if err := a.Render(w, "track/upload.html", a.BaseTemplate(r).WithForm(form)); err != nil {
		a.ServerError(w, err)
	}
}

func (h *Handler) PostUploadTrack(w http.ResponseWriter, r *http.Request) {
	a := h.app

	inF, hdr, err := r.FormFile("file")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			a.ClientError(w, err, http.StatusBadRequest)
		} else {
			a.ServerError(w, err)
		}
		return
	}
	defer inF.Close()

	form := form.File{Filename: hdr.Filename}
	if !form.Validate() {
		if err := a.Render(w, "track/upload.html", a.BaseTemplate(r).WithForm(form)); err != nil {
			a.ServerError(w, err)
		}
		return
	}

	uPath := filepath.Join(app.TMPDir, hdr.Filename)
	outF, err := os.Create(uPath)
	if err != nil {
		a.ServerError(w, err)
		return
	}
	defer outF.Close()

	if _, err = io.Copy(outF, inF); err != nil {
		a.ServerError(w, err)
		return
	}

	a.FlashInfo(r.Context(), "File "+hdr.Filename+" uploaded.")
	a.LogAction(r.Context(), "file upload", "filename", hdr.Filename)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
