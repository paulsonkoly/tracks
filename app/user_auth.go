package app

import (
	"context"
	"database/sql"
	"errors"

	"github.com/paulsonkoly/tracks/repository"
)

func (a *App) AuthenticateUser(ctx context.Context, name, password string) (*repository.User, error) {
	user, err := a.Repo.GetUserByName(ctx, name)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	// TODO not hashed
	if err != nil || user.HashedPassword != password {
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
