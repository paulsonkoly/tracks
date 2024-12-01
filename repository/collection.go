package repository

import (
	"database/sql"
	"errors"

	"github.com/paulsonkoly/tracks/repository/sqlc"
)

// Collection is a collection of [Track]s.
type Collection struct {
	ID     int     `json:"-"`
	Name   string  `json:"-"`      // Name of the collection.
	User   User    `json:"-"`      // User is the owner of the collection.
	Tracks []Track `json:"tracks"` // Tracks are the contained tracks.
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

// CollectionNameExists checks if the collection with the given name exists.
func (q Queries) CollectionNameExists(name string) (bool, error) {
	_, err := q.sqlc.GetCollectionByName(q.ctx, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// GetCollection fetches the collection with the given id.
func (q Queries) GetCollection(id int) (Collection, error) {
	name, err := q.sqlc.GetCollectionName(q.ctx, int32(id))
	if err != nil {
		return Collection{}, err
	}

	return Collection{ID: id, Name: name}, nil
}

// GetCollectionTracks fetches the tracks in a collection.
func (q Queries) GetCollectionTracks(id int) (result Collection, err error) {
	trks, err := q.sqlc.GetCollectionTracks(q.ctx, int32(id))
	if err != nil {
		return result, nil
	}

	result.Tracks = make([]Track, len(trks))

	for i, trk := range trks {
		result.Tracks[i].ID = int(trk.ID)
		result.Tracks[i].Name = trk.Name
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
