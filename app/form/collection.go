package form

// Collection is a form object for track collections.
type Collection struct {
	ID     int     // ID is the collection id.
	Name   string  // Name is the collection name.
	Tracks []Track // Tracks is the collection tracks.

	errors
}

type Track struct {
	ID   int    // ID is the track id.
	Name string // Name is the track name.
}

// CollectionUniqueChecker checks if a collection name is unqiue.
type CollectionUniqueChecker interface {
	CollectionUnique(name string) (bool, error)
}

// Validate validates the collection data.
func (c *Collection) Validate(uniq CollectionUniqueChecker) (bool, error) {
	c.validateName()

	ok, err := uniq.CollectionUnique(c.Name)
	if err != nil {
		return false, err
	}
	if !ok {
		c.AddError("Name already taken.")
	}

	return c.valid(), nil
}

func (c *Collection) validateName() {
	if c.Name == "" {
		c.AddFieldError("Name", "Name is required.")
	}
}
