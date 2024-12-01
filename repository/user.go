package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/paulsonkoly/tracks/repository/sqlc"
)

// User contains user data.
type User struct {
	ID             int
	Username       string
	HashedPassword string
	CreatedAt      time.Time
}

// InsertUser inserts user into database.
func (q Queries) InsertUser(username, hashedPassword string) (User, error) {
	var user User
	qr, err := q.sqlc.InsertUser(q.ctx,
		sqlc.InsertUserParams{
			Username:       username,
			HashedPassword: hashedPassword,
		})
	if err != nil {
		return user, err
	}
	user.ID = int(qr.ID)
	user.Username = qr.Username
	user.HashedPassword = qr.HashedPassword
	user.CreatedAt = qr.CreatedAt
	return user, err
}

// UpdateUser updates user with given id with username and hashed password.
func (q Queries) UpdateUser(id int, username, hashedPassword string) error {
	err := q.sqlc.UpdateUser(q.ctx,
		sqlc.UpdateUserParams{
			ID:             int32(id),
			Username:       username,
			HashedPassword: hashedPassword,
		})
	return err
}

// GetUser returns user with given id.
func (q Queries) GetUser(id int) (User, error) {
	var user User
	qr, err := q.sqlc.GetUser(q.ctx, int32(id))
	if err != nil {
		return user, err
	}

	user.ID = int(qr.ID)
	user.Username = qr.Username
	user.HashedPassword = qr.HashedPassword
	user.CreatedAt = qr.CreatedAt
	return user, err
}

// GetUserByName returns user with given username.
func (q Queries) GetUserByName(username string) (User, error) {
	var user User
	qr, err := q.sqlc.GetUserByName(q.ctx, username)
	if err != nil {
		return user, err
	}

	user.ID = int(qr.ID)
	user.Username = qr.Username
	user.HashedPassword = qr.HashedPassword
	user.CreatedAt = qr.CreatedAt
	return user, err
}

// GetUserByNameNotID returns user with given username excluding users with matching id.
func (q Queries) GetUserByNameNotID(username string, id int) (User, error) {
	var user User
	qr, err := q.sqlc.GetUserByNameNotID(q.ctx,
		sqlc.GetUserByNameNotIDParams{
			Username: username,
			ID:       int32(id),
		})
	if err != nil {
		return user, err
	}

	user.ID = int(qr.ID)
	user.Username = qr.Username
	user.HashedPassword = qr.HashedPassword
	user.CreatedAt = qr.CreatedAt
	return user, err
}

// GetUsers returns all users.
func (q Queries) GetUsers() ([]User, error) {
	qrs, err := q.sqlc.GetUsers(q.ctx)
	if err != nil {
		return []User{}, err
	}
	users := make([]User, 0, len(qrs))
	for _, qr := range qrs {
		users = append(users,
			User{

				ID:             int(qr.ID),
				Username:       qr.Username,
				HashedPassword: qr.HashedPassword,
				CreatedAt:      qr.CreatedAt,
			})
	}
	return users, err
}

// DeleteUser deletes user with given id.
func (q Queries) DeleteUser(id int) error {
	return q.sqlc.DeleteUser(q.ctx, int32(id))
}

// UsernameExists checks if the username exists in the database.
func (q Queries) UsernameExists(username string) (bool, error) {
	_, err := q.sqlc.GetUserByName(q.ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// UsernameExistsNotID checks if user with given username exists except given id.
func (q Queries) UsernameExistsNotID(id int, username string) (bool, error) {
	_, err := q.sqlc.GetUserByNameNotID(q.ctx, sqlc.GetUserByNameNotIDParams{ID: int32(id), Username: username})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
