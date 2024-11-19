package repository

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/paulsonkoly/tracks/repository/sqlc"
)

// Tracks stores data about a GPX track. A GPX track belongs to a user and
// belongs to a gpx file and has many segments.
type Track struct {
	ID           int
	Name         string
	Type         sqlc.Tracktype // TODO remove sqlc
	CreatedAt    time.Time
	Time         *time.Time
	LengthMeters float64
	User         *User
}

// InsertTrack inserts a new track into the database.
func (q Queries) InsertTrack(gpxFileID int, t sqlc.Tracktype, name string, userID int) (int, error) {
	id, err := q.sqlc.InsertTrack(q.ctx,
		sqlc.InsertTrackParams{
			GpxfileID: int32(gpxFileID),
			Type:      t,
			Name:      name,
			UserID:    int32(userID),
		})
	return int(id), err
}

// GetTrack retrieves track data for id. It does not retrieve associated models.
func (q Queries) GetTrack(id int) (Track, error) {
	var r Track

	qr, err := q.sqlc.GetTrack(q.ctx, int32(id))
	if err != nil {
		return r, err
	}
	r.ID = int(qr.ID)
	r.Name = qr.Name
	r.Type = qr.Type
	r.CreatedAt = qr.CreatedAt
	if qr.Time.Valid {
		t := qr.Time.Time
		r.Time = &t
	}
	r.LengthMeters = qr.TrackLengthMeters
	return r, nil
}

// GetTracks retrieves all tracks and associated users.
func (q Queries) GetTracks() ([]Track, error) {
	r := []Track{}

	qrs, err := q.sqlc.GetTracks(q.ctx)
	if err != nil {
		return nil, err
	}
	for _, qr := range qrs {
		u := User{
			ID:             int(qr.User.ID),
			Username:       qr.User.Username,
			HashedPassword: qr.User.HashedPassword,
			CreatedAt:      qr.User.CreatedAt,
		}
		r = append(r,
			Track{
				ID:           int(qr.Track.ID),
				Name:         qr.Track.Name,
				Type:         qr.Track.Type,
				LengthMeters: qr.TrackLengthMeters,
				CreatedAt:    qr.Track.CreatedAt,
				User:         &u,
			})
	}
	return r, nil
}

// Point is a pair of longitude and latitude.
type Point struct {
	Longitude, Latitude float64
}

// InsertTrackSegment inserts a track segment into the database.
func (q Queries) InsertTrackSegment(id int, points []Point) error {
	gBuilder := strings.Builder{}
	gBuilder.WriteString("LINESTRING(")
	for ix, point := range points {
		if ix > 0 {
			gBuilder.WriteString(", ")
		}
		gBuilder.WriteString(fmt.Sprintf("%f %f", point.Longitude, point.Latitude))
	}
	gBuilder.WriteString(")")
	return q.sqlc.InsertSegment(q.ctx, sqlc.InsertSegmentParams{TrackID: int32(id), Geometry: gBuilder.String()})
}

// GetTrackPoints returns the points of a track.
func (q Queries) GetTrackPoints(id int) ([]Point, error) {
	qrs, err := q.sqlc.GetTrackSegmentPoints(q.ctx, int32(id))
	if err != nil {
		return nil, err
	}
	r := []Point{}
	for _, qr := range qrs {
		long, ok := qr.Longitude.(float64)
		if !ok {
			return nil, errors.New("invalid longitude")
		}
		lat, ok := qr.Latitude.(float64)
		if !ok {
			return nil, errors.New("invalid latitude")
		}
		r = append(r, Point{
			Latitude:  lat,
			Longitude: long,
		})
	}
	return r, nil
}
