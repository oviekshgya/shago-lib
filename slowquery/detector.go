package slowquery

import (
	"fmt"
	"time"
)

type Config struct {
	Threshold time.Duration
	OnSlow    func(operation string, duration time.Duration)
}

// TrackFunc executes a function and reports if it exceeds threshold
func TrackFunc(cfg Config, operation string, fn func()) {
	start := time.Now()
	fn()
	duration := time.Since(start)

	if duration >= cfg.Threshold {
		if cfg.OnSlow != nil {
			cfg.OnSlow(operation, duration)
		} else {
			// Default behavior: fmt.Println? Or use a separate logger
			fmt.Printf("[SLOW QUERY] operation=%s duration=%v\n", operation, duration)
		}
	}
}

// TODO: Middleware for SQL/NoSQL drivers usually requires wrapping specific interfaces
// which might be too heavy for a general lib. `TrackFunc` is generic enough.
