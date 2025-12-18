package benchmark

import (
	"fmt"
	"time"
)

type Profiler struct {
	operation string
	start     time.Time
}

func New(operation string) *Profiler {
	return &Profiler{
		operation: operation,
		start:     time.Now(),
	}
}

func (p *Profiler) End() {
	duration := time.Since(p.start)
	fmt.Printf("[BENCHMARK] %s took %v\n", p.operation, duration)
}

func Track(operation string, fn func()) time.Duration {
	start := time.Now()
	fn()
	duration := time.Since(start)
	fmt.Printf("[BENCHMARK] %s took %v\n", operation, duration)
	return duration
}
