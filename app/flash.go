package app

import "context"

type Flash = map[string][]string

func (a *App) flash(ctx context.Context, typ string, msg string) {
	flash, ok := a.sm.Get(ctx, flashKey).(Flash)
	if !ok {
		flash = make(Flash)
	}
	msgs, ok := flash[typ]
	if !ok {
		msgs = []string{}
	}
	msgs = append(msgs, msg)

	flash[typ] = msgs
	a.sm.Put(ctx, flashKey, flash)
}

func (a *App) FlashInfo(ctx context.Context, msg string) {
	a.flash(ctx, "info", msg)
}

func (a *App) FlashWarning(ctx context.Context, msg string) {
	a.flash(ctx, "warning", msg)
}

func (a *App) FlashError(ctx context.Context, msg string) {
	a.flash(ctx, "error", msg)
}
