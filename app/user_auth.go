package app

import (
	"context"
	"database/sql"
	"errors"

	"github.com/paulsonkoly/tracks/repository"
	"golang.org/x/crypto/bcrypt"
)

func (a *App) AuthenticateUser(ctx context.Context, name, password string) (*repository.User, error) {
	user, err := a.Repo.GetUserByName(ctx, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		return nil, err
	}

	a.SM.Put(ctx, currentUserID, user.ID)
	return &user, nil
}

func (a *App) ClearCurrentUser(ctx context.Context) {
	a.SM.Remove(ctx, currentUserID)
}

func (a *App) CurrentUser(ctx context.Context) *repository.User {
	currentUser, ok := ctx.Value(CurrentUser).(repository.User)
	if !ok {
		return nil
	}

	return &currentUser
}
