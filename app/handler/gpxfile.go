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
	"golang.org/x/net/context"
)

// GPXFiles handler renders the file list page.
func (h *Handler) GPXFiles(w http.ResponseWriter, r *http.Request) {
	a := h.app

	form := form.GPXFile{}
	files, err := a.Q(r.Context()).GetGPXFiles()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		a.ServerError(w, err)
		return
	}
	if err := a.Render(w, "gpxfile/gpxfiles.html", a.BaseTemplate(r).WithForm(form).WithGPXFiles(files)); err != nil {
		a.ServerError(w, err)
	}
}

// PostUploadGPXFile handles the post request for file upload.
func (h *Handler) PostUploadGPXFile(w http.ResponseWriter, r *http.Request) {
	a := h.app
	uid := a.CurrentUser(r.Context()).ID

	var (
		id    int
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
	err = a.WithTx(r.Context(), func(ctx context.Context) error {

		files, err := a.Q(ctx).GetGPXFiles()
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			a.ServerError(w, err)
			return err
		}

		form := form.GPXFile{Filename: hdr.Filename}
		ok, err := form.Validate(a.Q(ctx))
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

		id, err = a.Q(ctx).InsertGPXFile(hdr.Filename, size, uid)
		if err != nil {
			os.Remove(uPath)
			return err
		}

		a.FlashInfo(ctx, "File "+hdr.Filename+" uploaded.")
		a.LogAction(ctx, "file upload", "filename", hdr.Filename)
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

// DeleteGPXFile deletes a GPX file.
func (h *Handler) DeleteGPXFile(w http.ResponseWriter, r *http.Request) {
	a := h.app

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		a.ClientError(w, err, http.StatusBadRequest)
		return
	}

	filename, err := a.Q(r.Context()).DeleteGPXFile(id)
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
