package middleware

import (
	"log"
	"time"

	"github.com/bigelle/botify"
)

// LoggingMiddleware is a middleware that logs update information including the update ID and processing duration,
// excluding network latencies, serializations and deserializations.
func LoggingMiddleware(next botify.HandlerFunc) botify.HandlerFunc {
	return func(ctx *botify.Context) {
		start := time.Now()

		next(ctx)

		end := time.Since(start)

		for _, req := range ctx.SendedRequests() {
			end = end - req.Duration
		}

		log.Printf("%s ID=%d %v", ctx.UpdateType(), ctx.UpdateID(), end)
	}
}

// RecoveryMiddleware is a middleware that wraps next [botify.HandlerFunc] with a built-in Go recovery function.
func RecoveryMiddleware(next botify.HandlerFunc) botify.HandlerFunc {
	return func(ctx *botify.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("PANIC in handler for update_type=%s with ID=%d", ctx.UpdateType(), ctx.UpdateID())
			}
		}()

		next(ctx)
	}
}
