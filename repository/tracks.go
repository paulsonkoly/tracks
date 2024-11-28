package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/paulsonkoly/tracks/repository/sqlc"
)

// Tracks stores data about a GPX track. A GPX track belongs to a user and
// belongs to a gpx file and has many segments.
type Track struct {
	ID           int            `json:"id"`
	Name         string         `json:"name"`
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

// GetTracks retrieves tracks with matching names. Only name and id is set.
func (q Queries) GetMatchingTracks(name string) ([]Track, error) {
	r := []Track{}

	qrs, err := q.sqlc.GetMatchingTracks(q.ctx, fmt.Sprintf("%%%s%%", name))
	if err != nil {
		return nil, err
	}
	for _, qr := range qrs {
		r = append(r, Track{ID: int(qr.ID), Name: qr.Name})
	}
	return r, nil
}

func (q Queries) TrackIDsPresent(ids []int) (bool, error) {
	conv := make([]int32, len(ids))
	for i, id := range ids {
		conv[i] = int32(id)
	}

	badIDs, err := q.sqlc.GetNonExistentTrackIDs(q.ctx, conv)
	if err != nil {
		return false, err
	}
	return len(badIDs) == 0, nil
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
func (q Queries) GetTrackPoints(id int) ([]Segment, error) {
	result := []Segment{}

	sIDs, err := q.sqlc.GetTrackSegments(q.ctx, int32(id))
	if err != nil {
		return nil, err
	}

	for _, sID := range sIDs {
		pts, err := q.sqlc.GetSegmentPoints(q.ctx, sID)
		if err != nil {
			return nil, err
		}

		conv := make(Segment, len(pts))
		for i, p := range pts {
			conv[i] = Point{
				Latitude:  p.Latitude,
				Longitude: p.Longitude,
			}
		}

		result = append(result, conv)
	}

	return result, nil
}
