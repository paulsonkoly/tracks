package app

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/paulsonkoly/tracks/repository"
	"github.com/tkrajina/gpxgo/gpx"
)

func (a *App) ProcessGPXFile(path string, id int32) {
	_ = a.WithTx(context.Background(), func(h TXHandle) error {
		gpxF, err := gpx.ParseFile(path)
		if err != nil {
			goto Failed
		}

		for _, track := range gpxF.Tracks {

			tid, err := a.Repo(h).InsertTrack(context.Background(), repository.InsertTrackParams{GpxfileID: id, Type: repository.TracktypeTrack, Name: track.Name})
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
