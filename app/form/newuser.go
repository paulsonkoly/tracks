package form

import (
	"context"
	"database/sql"
	errs "errors"
	"unicode/utf8"

	"github.com/paulsonkoly/tracks/repository"
)

type NewUserForm struct {
	Username        string
	Password        string
	PasswordConfirm string
	errors
}

func (f *NewUserForm) Validate(ctx context.Context, repo * repository.Queries) {
  if utf8.RuneCountInString(f.Username) < 3 {
    f.AddFieldError("Username", "Username too short. Must be at least 3 characters long.")
  }

  if utf8.RuneCountInString(f.Password) < 6 {
    f.AddFieldError("Password", "Password too short. Must be at least 6 characters long.")
  }

  if f.Password != f.PasswordConfirm {
    f.AddFieldError("PasswordConfirm", "Password confirmation does not match password.")
  }

  _, err := repo.GetUserByName(ctx, f.Username)
  if !errs.Is(err, sql.ErrNoRows) {
    f.AddError("User already exists.")
  }
}
