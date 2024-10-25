package handlers

import "github.com/paulsonkoly/tracks/app"

type Handler struct {
	app *app.App
}

func New(app *app.App) *Handler {
	return &Handler{app: app}
}
