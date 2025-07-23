package middleware

import (
	"fmt"

	"github.com/bigelle/botify"
)

// RecoveryMiddleware is a middleware that wraps next [botify.HandlerFunc] with a built-in Go recovery function.
func RecoveryMiddleware(next botify.HandlerFunc) botify.HandlerFunc {
	return func(ctx *botify.Context) error {
		defer func() {
			if r := recover(); r != nil {
				ctx.Bot().Logger.Error(fmt.Errorf("%+v", r), "PANIC in handler for update", "type", ctx.UpdateType(), "ID", ctx.UpdateID())
			}
		}()
		next(ctx)
		return nil
	}
}
