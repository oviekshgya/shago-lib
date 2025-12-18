package benchmark

import "sync"

type Stats struct {
	// Simple slice for now, could use a proper Histogram library
	durations []int64
	mu        sync.Mutex
}

func (s *Stats) Record(d int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.durations = append(s.durations, d)
}
