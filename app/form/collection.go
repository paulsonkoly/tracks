package form

// Collection is a form object for track collections.
type Collection struct {
	ID       int    `form:"-"`    // ID is the collection id.
	Name     string `form:"name"` // Name is the collection name.
	TrackIDs []int  `form:"track_ids[]"`

	errors `form:"-"`
}

type Track struct {
	ID   int    // ID is the track id.
	Name string // Name is the track name.
}

// CollectionUniqueChecker checks if a collection name is unqiue.
type CollectionUniqueChecker interface {
	CollectionUnique(name string) (bool, error)
}

type TrackIDsPresentChecker interface {
	TrackIDsPresent(ids []int) (bool, error)
}

type formChecker interface {
	CollectionUniqueChecker
	TrackIDsPresentChecker
}

// Validate validates the collection data.
func (c *Collection) Validate(chk formChecker) (bool, error) {
	c.validateName()

	ok, err := chk.CollectionUnique(c.Name)
	if err != nil {
		return false, err
	}
	if !ok {
		c.AddError("Name already taken.")
	}

	ok, err = chk.TrackIDsPresent(c.TrackIDs)
	if err != nil {
		return false, err
	}
	if !ok {
		c.AddError("Tracks are not valid.")
	}

	return c.valid(), nil
}

func (c *Collection) validateName() {
	if c.Name == "" {
		c.AddFieldError("Name", "Name is required.")
	}
}
