package app

import (
	"context"
	"errors"

	"github.com/paulsonkoly/tracks/repository"
	"github.com/paulsonkoly/tracks/repository/sqlc"
	"github.com/tkrajina/gpxgo/gpx"
)

// ProcessGPXFile reads a GPX file pointed by path and loads the GPX data into
// the database. It assumes that the file is already created in the gpxfiles
// table, with status [repository.Filestatus]/uploaded. If the processing fails
// the record will be updated with [repository.Filestatus]/ProcessingFailed,
// otherwise [repository.Filestatus]/Processed.
func (a *App) ProcessGPXFile(path string, id int, uid int) {
	_ = a.WithTx(context.Background(), func(ctx context.Context) error {
		gpxF, err := gpx.ParseFile(path)
		if err != nil {
			goto Failed
		}

		err = a.Q(ctx).UpdateGPXFile(
			repository.UpdateGPXFileParams{
				ID:               id,
				Version:          &gpxF.Version,
				Creator:          &gpxF.Creator,
				Name:             &gpxF.Name,
				Description:      &gpxF.Description,
				AuthorName:       &gpxF.AuthorName,
				AuthorEmail:      &gpxF.AuthorEmail,
				AuthorLink:       &gpxF.AuthorLink,
				AuthorLinkText:   &gpxF.AuthorLinkText,
				AuthorLinkType:   &gpxF.AuthorLinkType,
				Copyright:        &gpxF.Copyright,
				CopyrightYear:    &gpxF.CopyrightYear,
				CopyrightLicense: &gpxF.CopyrightLicense,
				LinkText:         &gpxF.LinkText,
				LinkType:         &gpxF.LinkType,
				Time:             gpxF.Time,
				Keywords:         &gpxF.Keywords,
			})
		if err != nil {
			goto Failed
		}

		for _, track := range gpxF.Tracks {
			tid, err := a.Q(ctx).InsertTrack(id, sqlc.TracktypeTrack, track.Name, uid)
			if err != nil {
				goto Failed
			}

			for _, segment := range track.Segments {
				points := []repository.Point{}

				for _, point := range segment.Points {
					points = append(points, repository.Point{Latitude: point.Latitude, Longitude: point.Longitude})
				}

				err = a.Q(ctx).InsertTrackSegment(tid, points)
				if err != nil {
					goto Failed
				}
			}
		}

		for _, route := range gpxF.Routes {
			tid, err := a.Q(ctx).InsertTrack(id, sqlc.TracktypeRoute, route.Name, uid)
			if err != nil {
				goto Failed
			}

			points := []repository.Point{}
			for _, point := range route.Points {
				points = append(points, repository.Point{Latitude: point.Latitude, Longitude: point.Longitude})
			}

			err = a.Q(ctx).InsertTrackSegment(tid, points)
			if err != nil {
				goto Failed
			}
		}

		err = a.Q(ctx).SetGPXFileStatus(id, sqlc.FilestatusProcessed)
		return err

	Failed:
		// this is outside of the transaction we roll back
		err2 := a.Q(context.Background()).SetGPXFileStatus(id, sqlc.FilestatusProcessingFailed)
		return errors.Join(err, err2)
	})
}
