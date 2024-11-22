package app

import "context"

// Flash is a message that appears on the site for the next request.
//
// Flash messages are stored in the session and thus they require the
// [app.Dynamic] middleware.
type Flash = map[string][]string

func (a *App) flash(ctx context.Context, typ string, msg string) {
	flash, ok := a.sm.Get(ctx, SKFlash).(Flash)
	if !ok {
		flash = make(Flash)
	}
	msgs, ok := flash[typ]
	if !ok {
		msgs = []string{}
	}
	msgs = append(msgs, msg)

	flash[typ] = msgs
	a.sm.Put(ctx, SKFlash, flash)
}

// FlashInfo is a flash message with general information.
func (a *App) FlashInfo(ctx context.Context, msg string) {
	a.flash(ctx, "info", msg)
}

// FlashWarning is a flash message that indicates a warning.
func (a *App) FlashWarning(ctx context.Context, msg string) {
	a.flash(ctx, "warning", msg)
}

// FlashError is a flash message that indicates an error.
func (a *App) FlashError(ctx context.Context, msg string) {
	a.flash(ctx, "error", msg)
}
