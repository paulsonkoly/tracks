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
func (q Queries) InsertCollection(name string, user User, tracks []Track) error {
	trackIDs := make([]int32, len(tracks))
	for i, track := range tracks {
		trackIDs[i] = int32(track.ID)
	}
	return q.sqlc.InsertCollection(q.ctx, sqlc.InsertCollectionParams{
		Name:     name,
		UserID:   int32(user.ID),
		TrackIds: trackIDs,
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
