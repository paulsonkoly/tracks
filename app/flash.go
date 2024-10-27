package app

import "context"

type Flash = map[string][]string

func (a *App) flash(ctx context.Context, typ string, msg string) {
  sm := a.SM
  flash, ok := sm.Get(ctx, "flash").(Flash)
  if ! ok {
    flash = make(Flash)
  }
  msgs, ok := flash[typ]
  if !ok {
    msgs = []string{}
  }
  msgs = append(msgs, msg)

  flash[typ] = msgs
  sm.Put(ctx, "flash", flash)
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
