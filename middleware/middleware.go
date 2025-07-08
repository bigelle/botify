package middleware

import (
	"log"
	"time"

	"github.com/bigelle/botify"
)

func LoggingMiddleware(next botify.HandlerFunc) botify.HandlerFunc {
	return  func(ctx botify.Context) {
		start := time.Now()
		
		next(ctx)
		
		end := time.Since(start)
	
		log.Printf("%s ID=%d %v", ctx.UpdateType(), ctx.UpdateID(), end)
	}
}

func RecoveryMiddleware(next botify.HandlerFunc) botify.HandlerFunc {
	return  func(ctx botify.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("PANIC in handler for update_type=%s with ID=%d", ctx.UpdateType(), ctx.UpdateID())
			}
		}()

		next(ctx)
	}
}
