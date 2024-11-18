package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/paulsonkoly/tracks/repository"
	"github.com/tkrajina/gpxgo/gpx"
)

// ProcessGPXFile reads a GPX file pointed by path and loads the GPX data into
// the database. It assumes that the file is already created in the gpxfiles
// table, with status [repository.Filestatus]/uploaded. If the processing fails
// the record will be updated with [repository.Filestatus]/ProcessingFailed,
// otherwise [repository.Filestatus]/Processed.
func (a *App) ProcessGPXFile(path string, id int32, uid int32) {
	_ = a.WithTx(context.Background(), func(h TXHandle) error {
		gpxF, err := gpx.ParseFile(path)
		if err != nil {
			goto Failed
		}

		err = a.Repo(h).UpdateGPXFile(context.Background(),
			repository.UpdateGPXFileParams{
				ID:               id,
				Version:          nullString(gpxF.Version),
				Creator:          nullString(gpxF.Creator),
				Name:             nullString(gpxF.Name),
				Description:      nullString(gpxF.Description),
				AuthorName:       nullString(gpxF.AuthorName),
				AuthorEmail:      nullString(gpxF.AuthorEmail),
				AuthorLink:       nullString(gpxF.AuthorLink),
				AuthorLinkText:   nullString(gpxF.AuthorLinkText),
				AuthorLinkType:   nullString(gpxF.AuthorLinkType),
				Copyright:        nullString(gpxF.Copyright),
				CopyrightYear:    nullString(gpxF.CopyrightYear),
				CopyrightLicense: nullString(gpxF.CopyrightLicense),
				LinkText:         nullString(gpxF.LinkText),
				LinkType:         nullString(gpxF.LinkType),
				Time:             nullTime(gpxF.Time),
				Keywords:         nullString(gpxF.Keywords),
			})
		if err != nil {
			goto Failed
		}

		for _, track := range gpxF.Tracks {
			tid, err := a.Repo(h).InsertTrack(context.Background(),
				repository.InsertTrackParams{
					GpxfileID: id,
					Type:      repository.TracktypeTrack,
					Name:      track.Name,
					UserID:    uid,
				})
			if err != nil {
				goto Failed
			}

			for _, segment := range track.Segments {

				gBuilder := strings.Builder{}
				gBuilder.WriteString("LINESTRING(")
				for ix, point := range segment.Points {
					if ix > 0 {
						gBuilder.WriteString(", ")
					}
					gBuilder.WriteString(fmt.Sprintf("%f %f", point.Longitude, point.Latitude))
				}
				gBuilder.WriteString(")")

				err = a.Repo(h).InsertSegment(context.Background(), repository.InsertSegmentParams{TrackID: tid, Geometry: gBuilder.String()})
				if err != nil {
					goto Failed
				}
			}
		}

		for _, route := range gpxF.Routes {
			tid, err := a.Repo(h).InsertTrack(context.Background(), repository.InsertTrackParams{GpxfileID: id, Type: repository.TracktypeRoute, Name: route.Name})
			if err != nil {
				goto Failed
			}

			gBuilder := strings.Builder{}
			gBuilder.WriteString("LINESTRING(")
			for ix, point := range route.Points {
				if ix > 0 {
					gBuilder.WriteString(", ")
				}
				gBuilder.WriteString(fmt.Sprintf("%f %f", point.Longitude, point.Latitude))
			}
			gBuilder.WriteString(")")

			err = a.Repo(h).InsertSegment(context.Background(), repository.InsertSegmentParams{TrackID: tid, Geometry: gBuilder.String()})
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
}

func nullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: s, Valid: true}
}

func nullTime(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{}
	}
	return sql.NullTime{Time: *t, Valid: true}
}
