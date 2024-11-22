package app

import (
	"context"
	"database/sql"
	"errors"

	"github.com/paulsonkoly/tracks/repository"
	"golang.org/x/crypto/bcrypt"
)

// ErrAuthenticationFailed indicates that either the username or the password was wrong.
var ErrAuthenticationFailed = errors.New("authentication failed")

// SKCurrentUserID is the session key for the current user id.
const SKCurrentUserID = "currentUserID"

// AuthenticateUser attempts a user login. I returns ErrAuthenticationFailed in
// case of invalid credentials. It returns the logged in user in case of
// successful login.
//
// The login is stored in the session for following requests, until a user
// logout happens, or the session expires.
func (a *App) AuthenticateUser(ctx context.Context, name, password string) (*repository.User, error) {
	user, err := a.Q(ctx).GetUserByName(name)
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

	a.sm.Put(ctx, SKCurrentUserID, user.ID)

	if err := a.sm.RenewToken(ctx); err != nil {
		return nil, err
	}

	return &user, nil
}

// ClearCurrentUser logs out the current user by removing it from the session.
func (a *App) ClearCurrentUser(ctx context.Context) error {
	a.sm.Remove(ctx, SKCurrentUserID)

	if err := a.sm.RenewToken(ctx); err != nil {
		return err
	}

	return nil
}

// CurrentUser returns the current user from the request context. This relies
// on the [app.Dynamic] middleware to transfer the user from the session to the
// request context.
func (a *App) CurrentUser(ctx context.Context) *repository.User {
	currentUser, ok := ctx.Value(CurrentUser).(repository.User)
	if !ok {
		return nil
	}

	return &currentUser
}
