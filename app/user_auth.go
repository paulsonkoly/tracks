package app

import (
	"context"
	"database/sql"
	"errors"

	"github.com/paulsonkoly/tracks/repository"
	"golang.org/x/crypto/bcrypt"
)

var ErrAuthenticationFailed = errors.New("authentication failed")

const currentUserID = "currentUserID"

func (a *App) AuthenticateUser(ctx context.Context, name, password string) (*repository.User, error) {
	user, err := a.Repo.GetUserByName(ctx, name)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, ErrAuthenticationFailed
	} else if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil && errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return nil, ErrAuthenticationFailed
	} else if err != nil {
		return nil, err
	}

	a.sm.Put(ctx, currentUserID, user.ID)

	if err := a.sm.RenewToken(ctx); err != nil {
		return nil, err
	}

	return &user, nil
}

func (a *App) ClearCurrentUser(ctx context.Context) error {
	a.sm.Remove(ctx, currentUserID)

	if err := a.sm.RenewToken(ctx); err != nil {
		return err
	}

	return nil
}

func (a *App) CurrentUser(ctx context.Context) *repository.User {
	currentUser, ok := ctx.Value(CurrentUser).(repository.User)
	if !ok {
		return nil
	}

	return &currentUser
}
