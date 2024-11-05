package form

import (
	"context"
	"unicode/utf8"
)

// User provides user data validation.
type User struct {
	ID              int    // ID points to the user to edit
	Username        string // Username is the login account
	Password        string // Password is the login password
	PasswordConfirm string // PasswordConfirm is the confirmation password
	errors
}

// UserUniqueChecker checks if the user does not exist in the database.
type UserUniqueChecker interface {
	UserUnique(ctx context.Context, username string) (bool, error)
}

// Validate validates the user data.
func (f *User) Validate(ctx context.Context, uniq UserUniqueChecker) (bool, error) {
	if err := f.validateUsername(ctx, uniq); err != nil {
		return false, err
	}
	f.validatePassword()

	return f.valid(), nil
}

// ValidateEdit validates the user data for editing. Empty data fields are not
// updated, so they are valid.
func (f *User) ValidateEdit(ctx context.Context, uniq UserUniqueChecker) (bool, error) {
	if f.Username != "" {
		if err := f.validateUsername(ctx, uniq); err != nil {
			return false, err
		}
	}

	if f.Password != "" || f.PasswordConfirm != "" {
		f.validatePassword()
	}

	return f.valid(), nil
}

func (f *User) validateUsername(ctx context.Context, uniq UserUniqueChecker) error {
	if utf8.RuneCountInString(f.Username) < 3 {
		f.AddFieldError("Username", "Username too short. Must be at least 3 characters long.")
	}

	ok, err := uniq.UserUnique(ctx, f.Username)
	if err != nil {
		return err
	}
	if !ok {
		f.AddError("User already exist.")
	}

	return nil
}

func (f *User) validatePassword() {
	if utf8.RuneCountInString(f.Password) < 6 {
		f.AddFieldError("Password", "Password too short. Must be at least 6 characters long.")
	}

	if f.Password != f.PasswordConfirm {
		f.AddFieldError("PasswordConfirm", "Passwords do not match.")
	}
}
