package repository

import (
	"database/sql"
	"errors"

	"github.com/paulsonkoly/tracks/repository/sqlc"
)

// Collection is a collection of [Track]s.
type Collection struct {
	ID     int
	Name   string  // Name of the collection.
	User   User    // User is the owner of the collection.
	Tracks []Track // Tracks are the contained tracks.
}

// InsertCollection adds a new collection of tracks with the given name.
func (q Queries) InsertCollection(name string, user User, trackIDs []int) error {
	conv := make([]int32, len(trackIDs))
	for i, id := range trackIDs {
		conv[i] = int32(id)
	}
	return q.sqlc.InsertCollection(q.ctx, sqlc.InsertCollectionParams{
		Name:     name,
		UserID:   int32(user.ID),
		TrackIds: conv,
	})
}

// CollectionUnique checks if the collection with the given name exists. Returns true if it doesn't.
func (q Queries) CollectionUnique(name string) (bool, error) {
	_, err := q.sqlc.GetCollectionByName(q.ctx, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return true, nil
		}
		return false, err
	}
	return false, nil
}

// GetCollection fetches the collection with the given id.
func (q Queries) GetCollection(id int) (Collection, error) {
	var result Collection
	name, err := q.sqlc.GetCollection(q.ctx, int32(id))
	if err != nil {
		return result, err
	}

	trks, err := q.sqlc.GetCollectionTracks(q.ctx, int32(id))
	if err != nil {
		return result, nil
	}

	result.Name = name
	result.Tracks = make([]Track, len(trks))

	for i, trk := range trks {
		result.Tracks[i].ID = int(trk.ID)
		result.Tracks[i].Name = trk.Name
	}

	return result, nil
}

// GetCollectionPoints returns the points of a track.
func (q Queries) GetCollectionPoints(id int) ([]Segment, error) {
	result := []Segment{}

	sIDs, err := q.sqlc.GetCollectionSegments(q.ctx, int32(id))
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

// GetCollections returns all collections.
func (q Queries) GetCollections() ([]Collection, error) {
	cs, err := q.sqlc.GetCollections(q.ctx)
	if err != nil {
		return nil, err
	}

	result := make([]Collection, len(cs))

	for i, c := range cs {
		result[i].Name = c.Name
		result[i].ID = int(c.ID)
	}

	return result, nil
}
