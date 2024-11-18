package handler

import (
	"database/sql"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/paulsonkoly/tracks/app"
	"github.com/paulsonkoly/tracks/app/form"
	"github.com/paulsonkoly/tracks/repository"
)

func (h *Handler) GPXFiles(w http.ResponseWriter, r *http.Request) {
	a := h.app

	form := form.GPXFile{}
	files, err := a.Repo(nil).GetGPXFiles(r.Context())
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		a.ServerError(w, err)
		return
	}
	if err := a.Render(w, "gpxfile/gpxfiles.html", a.BaseTemplate(r).WithForm(form).WithGPXFiles(files)); err != nil {
		a.ServerError(w, err)
	}
}

func (h *Handler) PostUploadGPXFile(w http.ResponseWriter, r *http.Request) {
	a := h.app
	uid := a.CurrentUser(r.Context()).ID

	var (
		id    int32
		uPath string
	)

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

	// flag to indicate if a file upload was successful and we need to background process it
	process := true
	err = a.WithTx(r.Context(), func(h app.TXHandle) error {

		files, err := a.Repo(nil).GetGPXFiles(r.Context())
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			a.ServerError(w, err)
			return err
		}

		form := form.GPXFile{Filename: hdr.Filename}
		ok, err := form.Validate(r.Context(), a.Repo(h))
		if err != nil {
			return err
		}
		// if not valid but no erros.
		if !ok {
			if err := a.Render(w, "gpxfile/gpxfiles.html", a.BaseTemplate(r).WithGPXFiles(files).WithForm(form)); err != nil {
				return err
			}
			process = false
			return nil
		}

		uPath = filepath.Join(app.TMPDir, hdr.Filename)
		outF, err := os.Create(uPath)
		if err != nil {
			return err
		}
		defer outF.Close()

		size, err := io.Copy(outF, inF)
		if err != nil {
			return err
		}

		id, err = a.Repo(h).InsertGPXFile(r.Context(),
			repository.InsertGPXFileParams{
				Filename: hdr.Filename,
				Filesize: size,
				UserID:   uid,
			})
		if err != nil {
			os.Remove(uPath)
			return err
		}

		a.FlashInfo(r.Context(), "File "+hdr.Filename+" uploaded.")
		a.LogAction(r.Context(), "file upload", "filename", hdr.Filename)
		http.Redirect(w, r, "/gpxfiles", http.StatusSeeOther)

		return nil
	})
	if err != nil {
		a.ServerError(w, err)
		return
	}

	// process uploaded file
	if process {
		go a.ProcessGPXFile(uPath, id, uid)
	}
}

// DeleteTrack deletes a GPX file.
func (h *Handler) DeleteGPXFile(w http.ResponseWriter, r *http.Request) {
	a := h.app

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		a.ClientError(w, err, http.StatusBadRequest)
		return
	}

	filename, err := a.Repo(nil).DeleteGPXFile(r.Context(), int32(id))
	if err != nil {
		a.ServerError(w, err)
		return
	}
	uPath := filepath.Join(app.TMPDir, filename)
	os.Remove(uPath)

	a.FlashInfo(r.Context(), "GPX file deleted.")
	a.LogAction(r.Context(), "gpx file deleted", slog.Int("id", id))
	http.Redirect(w, r, "/gpxfiles", http.StatusSeeOther)
}
