package form

import (
	"unicode/utf8"
)

// User provides user data validation.
type User struct {
	ID              int    `form:"-"`                // ID points to the user to edit
	Username        string `form:"username"`         // Username is the login account
	Password        string `form:"password"`         // Password is the login password
	PasswordConfirm string `form:"password_confirm"` // PasswordConfirm is the confirmation password
	errors          `form:"-"`
}

// UserPresenceChecker checks if the user exists in the database.
type UserPresenceChecker interface {
	// UserExists returns wether the username already exists in the database.
	UsernameExists(username string) (bool, error)

	// UsernameExistsNotID returns wether the username already exists in the database apart from checking the user with id.
	UsernameExistsNotID(id int, username string) (bool, error)
}

// Validate validates the user data.
func (f *User) Validate(check UserPresenceChecker) (bool, error) {
	f.validateUsername()

	exists, err := check.UsernameExists(f.Username)
	if err != nil {
		return false, err
	}
	if exists {
		f.AddError("Username taken.")
	}

	f.validatePassword()

	return f.valid(), nil
}

// ValidateEdit validates the user data for editing. Empty data fields are not
// updated, so they are valid. Username uniqueness is not validated *if* it's the same userid.
func (f *User) ValidateEdit(check UserPresenceChecker) (bool, error) {
	if f.Username != "" {
		f.validateUsername()

		exists, err := check.UsernameExistsNotID(f.ID, f.Username)
		if err != nil {
			return false, err
		}
		if exists {
			f.AddError("Username taken.")
		}
	}

	if f.Password != "" || f.PasswordConfirm != "" {
		f.validatePassword()
	}

	return f.valid(), nil
}

func (f *User) validateUsername() {
	if utf8.RuneCountInString(f.Username) < 3 {
		f.AddFieldError("Username", "Username too short. Must be at least 3 characters long.")
	}
}

func (f *User) validatePassword() {
	if utf8.RuneCountInString(f.Password) < 6 {
		f.AddFieldError("Password", "Password too short. Must be at least 6 characters long.")
	}

	if f.Password != f.PasswordConfirm {
		f.AddFieldError("PasswordConfirm", "Passwords do not match.")
	}
}
