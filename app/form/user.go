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

func (f *User) Validate(ctx context.Context, repo *repository.Queries) (bool, error) {
	if err := f.validateUsername(ctx, repo); err != nil {
		return false, err
	}
	f.validatePassword()

	return f.valid(), nil
}

func (f *User) ValidateEdit(ctx context.Context, repo *repository.Queries) (bool, error) {
	if f.Username != "" {
		if err := f.validateUsername(ctx, repo); err != nil {
			return false, err
		}
	}

	if f.Password != "" || f.PasswordConfirm != "" {
		f.validatePassword()
	}

	return f.valid(), nil
}

func (f *User) validateUsername(ctx context.Context, repo *repository.Queries) error {
	if utf8.RuneCountInString(f.Username) < 3 {
		f.AddFieldError("Username", "Username too short. Must be at least 3 characters long.")
	}
	_, err := repo.GetUserByName(ctx, f.Username)
	if err != nil {
		if !errs.Is(err, sql.ErrNoRows) {
			return err
		}
	} else {
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
