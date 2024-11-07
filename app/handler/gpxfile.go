package handler

import (
	"context"
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
	"github.com/tkrajina/gpxgo/gpx"
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

	err = a.WithTx(r.Context(), func(h app.TXHandle) error {

		form := form.GPXFile{Filename: hdr.Filename}
		ok, err := form.Validate(r.Context(), a.Repo(h))
		if err != nil {
			return err
		}
		if !ok {
			if err := a.Render(w, "gpxfile/gpxfiles.html", a.BaseTemplate(r).WithForm(form)); err != nil {
				a.ServerError(w, err)
			}
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
				Link:     "TODO link text"})
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
	}

	// process uploaded file
	go func() {
		// we see the status from the gpxfiles.status column, otherwise error
		// indication from this transaction is not really possible as the request
		// is gone at this point.
		_ = a.WithTx(context.Background(), func(h app.TXHandle) error {
			gpxF, err := gpx.ParseFile(uPath)
			if err != nil {
				goto Failed
			}

			for _, track := range gpxF.Tracks {
				err = a.Repo(h).InsertTrack(context.Background(), repository.InsertTrackParams{GpxfileID: id, Type: repository.TracktypeTrack, Name: track.Name})
				if err != nil {
					goto Failed
				}
			}

			for _, route := range gpxF.Routes {
				err = a.Repo(h).InsertTrack(context.Background(), repository.InsertTrackParams{GpxfileID: id, Type: repository.TracktypeRoute, Name: route.Name})
				if err != nil {
					goto Failed
				}
			}

			err = a.Repo(h).SetGPXFileStatus(context.Background(), repository.SetGPXFileStatusParams{ID: id, Status: repository.FilestatusProcessed})
			return err

		Failed:
			err2 := a.Repo(h).SetGPXFileStatus(context.Background(), repository.SetGPXFileStatusParams{ID: id, Status: repository.FilestatusProcessingFailed})
			return errors.Join(err, err2)
		})
	}()
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
