// Package handler defines the collection of request handlers.
//
// For any new request a new handler needs to be defined that produces the
// correct http response. A handler has access to the opaque app structure that
// provides application wide APIs for handlers.
package handler

import "github.com/paulsonkoly/tracks/app"

// Handler is an application request handler.
type Handler struct {
	app *app.App
}

// New creates a new handler.
func New(app *app.App) *Handler {
	return &Handler{app: app}
}
