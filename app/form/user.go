package form

import (
	"context"
	"database/sql"
	errs "errors"
	"unicode/utf8"

	"github.com/paulsonkoly/tracks/repository"
)

type User struct {
	ID              int
	Username        string
	Password        string
	PasswordConfirm string
	errors
}

func (f *User) Validate(ctx context.Context, repo *repository.Queries) {
	f.validateUsername()
	f.validatePassword()

	_, err := repo.GetUserByName(ctx, f.Username)
	if !errs.Is(err, sql.ErrNoRows) {
		f.AddError("User already exists.")
	}
}

func (f *User) ValidateEdit() {
	if f.Username != "" {
		f.validateUsername()
	}

	if f.Password != "" || f.PasswordConfirm != "" {
		f.validatePassword()
	}
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
