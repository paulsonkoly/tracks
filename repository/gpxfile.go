package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/paulsonkoly/tracks/repository/sqlc"
)

// GPXFile contains data related to a gpx file uploaded to the system.
type GPXFile struct {
	ID               int
	Filename         string
	Filesize         int64
	Status           sqlc.Filestatus // TODO: remove sqlc type
	Link             string          // TODO : nullable
	CreatedAt        time.Time
	UserID           int
	Version          *string
	Creator          *string
	Name             *string
	Description      *string
	AuthorName       *string
	AuthorEmail      *string
	AuthorLink       *string
	AuthorLinkText   *string
	AuthorLinkType   *string
	Copyright        *string
	CopyrightYear    *string
	CopyrightLicense *string
	LinkText         *string
	LinkType         *string
	Time             *time.Time
	Keywords         *string
	User             *User
}

// InsertGPXFile inserts a new GPX file and returns the new id.
func (q Queries) InsertGPXFile(filename string, filesize int64, userID int) (int, error) {
	id, err := q.sqlc.InsertGPXFile(q.ctx,
		sqlc.InsertGPXFileParams{
			Filename: filename,
			Filesize: filesize,
			UserID:   int32(userID),
		})
	return int(id), err
}

// SetGPXFileStatus sets the status of a GPX file.
func (q Queries) SetGPXFileStatus(id int, status sqlc.Filestatus) error {
	return q.sqlc.SetGPXFileStatus(q.ctx,
		sqlc.SetGPXFileStatusParams{
			ID:     int32(id),
			Status: status,
		})
}

// UpdateGPXFileParams are updateable properties of a GPX file.
type UpdateGPXFileParams struct {
	ID               int // ID identifies the file to update.
	Version          *string
	Creator          *string
	Name             *string
	Description      *string
	AuthorName       *string
	AuthorEmail      *string
	AuthorLink       *string
	AuthorLinkText   *string
	AuthorLinkType   *string
	Copyright        *string
	CopyrightYear    *string
	CopyrightLicense *string
	Link             string
	LinkText         *string
	LinkType         *string
	Time             *time.Time
	Keywords         *string
}

// UpdateGPXFile updates a GPX file.
func (q Queries) UpdateGPXFile(args UpdateGPXFileParams) error {
	var qArgs sqlc.UpdateGPXFileParams

	qArgs.ID = int32(args.ID)

	qArgs.Version = nullString(args.Version)
	qArgs.Creator = nullString(args.Creator)
	qArgs.Name = nullString(args.Name)
	qArgs.Description = nullString(args.Description)
	qArgs.AuthorName = nullString(args.AuthorName)
	qArgs.AuthorEmail = nullString(args.AuthorEmail)
	qArgs.AuthorLink = nullString(args.AuthorLink)
	qArgs.AuthorLinkText = nullString(args.AuthorLinkText)
	qArgs.AuthorLinkType = nullString(args.AuthorLinkType)
	qArgs.Copyright = nullString(args.Copyright)
	qArgs.CopyrightYear = nullString(args.CopyrightYear)
	qArgs.CopyrightLicense = nullString(args.CopyrightLicense)
	qArgs.Link = args.Link
	qArgs.LinkText = nullString(args.LinkText)
	qArgs.LinkType = nullString(args.LinkType)
	qArgs.Time = nullTime(args.Time)
	qArgs.Keywords = nullString(args.Keywords)

	return q.sqlc.UpdateGPXFile(q.ctx, qArgs)
}

// GetGPXFile returns a GPX file matching the id. It does not retrieve associated models.
func (q Queries) GetGPXFile(id int) (GPXFile, error) {
	var r GPXFile
	qr, err := q.sqlc.GetGPXFile(q.ctx, int32(id))
	if err != nil {
		return r, err
	}

	r.ID = int(qr.ID)
	r.Filename = qr.Filename
	r.Filesize = qr.Filesize
	r.Status = qr.Status
	r.Link = qr.Link
	r.CreatedAt = qr.CreatedAt
	r.Version = strConv(qr.Version)
	r.Creator = strConv(qr.Creator)
	r.Name = strConv(qr.Name)
	r.Description = strConv(qr.Description)
	r.AuthorName = strConv(qr.AuthorName)
	r.AuthorEmail = strConv(qr.AuthorEmail)
	r.AuthorLink = strConv(qr.AuthorLink)
	r.AuthorLinkText = strConv(qr.AuthorLinkText)
	r.AuthorLinkType = strConv(qr.AuthorLinkType)
	r.Copyright = strConv(qr.Copyright)
	r.Copyright = strConv(qr.Copyright)
	r.CopyrightYear = strConv(qr.CopyrightYear)
	r.CopyrightLicense = strConv(qr.CopyrightLicense)
	r.LinkText = strConv(qr.LinkText)
	r.LinkType = strConv(qr.LinkType)
	r.Time = timeConv(qr.Time)
	r.Keywords = strConv(qr.Keywords)

	return r, nil
}

// DeleteGPXFile deletes a GPX file and returns the filename.
func (q Queries) DeleteGPXFile(id int) (string, error) {
	return q.sqlc.DeleteGPXFile(q.ctx, int32(id))
}

// GetGPXFiles gets all GPX files with associated users.
func (q Queries) GetGPXFiles() ([]GPXFile, error) {
	qrs, err := q.sqlc.GetGPXFiles(q.ctx)
	if err != nil {
		return nil, err
	}

	result := []GPXFile{}
	for _, qr := range qrs {
		u := User{
			ID:             int(qr.User.ID),
			Username:       qr.User.Username,
			HashedPassword: qr.User.HashedPassword,
			CreatedAt:      qr.User.CreatedAt,
		}
		result = append(result, GPXFile{
			ID:        int(qr.ID),
			Filename:  qr.Filename,
			Filesize:  qr.Filesize,
			Status:    qr.Status,
			CreatedAt: qr.CreatedAt,
			User:      &u,
		})
	}

	return result, nil
}

// GPXFileUnique returns wether a GPX file will be unique with the given file
// name.
func (q Queries) GPXFileUnique(filename string) (bool, error) {
	_, err := q.sqlc.GetGPXFileByFilename(q.ctx, filename)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return true, nil
		}
		return false, err
	}
	return false, nil
}

func strConv(s sql.NullString) *string {
	if s.Valid {
		return &s.String
	}
	return nil
}

func timeConv(t sql.NullTime) *time.Time {
	if t.Valid {
		return &t.Time
	}
	return nil
}

func nullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *s, Valid: true}
}

func nullTime(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{}
	}
	return sql.NullTime{Time: *t, Valid: true}
}
