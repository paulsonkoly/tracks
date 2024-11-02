package handler

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/paulsonkoly/tracks/app"
	"github.com/paulsonkoly/tracks/app/form"
	"github.com/paulsonkoly/tracks/repository"
)

func (h *Handler) ViewTrack(w http.ResponseWriter, r *http.Request) {
	a := h.app

	if err := a.Render(w, "track/track.html", a.BaseTemplate(r)); err != nil {
		a.ServerError(w, err)
	}
}

func (h *Handler) UploadTrack(w http.ResponseWriter, r *http.Request) {
	a := h.app

	form := form.GPXFile{}
	if err := a.Render(w, "track/upload.html", a.BaseTemplate(r).WithForm(form)); err != nil {
		a.ServerError(w, err)
	}
}

func (h *Handler) PostUploadTrack(w http.ResponseWriter, r *http.Request) {
	a := h.app
	committed := false

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

	// create a transaction before the uniqueness validation
	tx, err := a.DB.BeginTx(r.Context(), nil)
	if err != nil {
		a.ServerError(w, err)
		return
	}
	defer func() {
		if !committed {
			if err := tx.Rollback(); err != nil {
				a.ServerError(w, err)
			}
		}
	}()

	form := form.GPXFile{Filename: hdr.Filename}
	ok, err := form.Validate(r.Context(), a.Repo.WithTx(tx))
	if err != nil {
		a.ServerError(w, err)
		return
	}
	if !ok {
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

	size, err := io.Copy(outF, inF)
	if err != nil {
		a.ServerError(w, err)
		return
	}

	if err := a.Repo.WithTx(tx).InsertGPXFile(r.Context(),
		repository.InsertGPXFileParams{
			Filename: hdr.Filename,
			Filesize: size,
			Link:     "TODO link text"}); err != nil {
		os.Remove(uPath)
		a.ServerError(w, err)
		return
	}

	if err := tx.Commit(); err != nil {
		a.ServerError(w, err)
		return
	}

	committed = true

	a.FlashInfo(r.Context(), "File "+hdr.Filename+" uploaded.")
	a.LogAction(r.Context(), "file upload", "filename", hdr.Filename)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
