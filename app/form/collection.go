package form

// Collection is a form object for track collections.
type Collection struct {
	ID       int    `form:"-"`    // ID is the collection id.
	Name     string `form:"name"` // Name is the collection name.
	TrackIDs []int  `form:"track_ids[]"`

	errors `form:"-"`
}

// CollectionUniqueChecker checks if a collection exists in the database.
type CollectionPresenceChecker interface {
	// CollectionNameExists checks if the collection name exists in the database.
	CollectionNameExists(name string) (bool, error)
}

// TrackPresenceChecker checks if the track is present in the database.
type TrackPresenceChecker interface {
	// TrackIDsExist checks if the track ids all exist in the database.
	TrackIDsExist(ids []int) (bool, error)
}

type formChecker interface {
	CollectionPresenceChecker
	TrackPresenceChecker
}

// Validate validates the collection data.
func (c *Collection) Validate(check formChecker) (bool, error) {
	c.validateName()

	exists, err := check.CollectionNameExists(c.Name)
	if err != nil {
		return false, err
	}
	if exists {
		c.AddError("Name already taken.")
	}

	exist, err := check.TrackIDsExist(c.TrackIDs)
	if err != nil {
		return false, err
	}
	if !exist {
		c.AddError("Tracks are not valid.")
	}

	return c.valid(), nil
}

func (c *Collection) validateName() {
	if c.Name == "" {
		c.AddFieldError("Name", "Name is required.")
	}
}
