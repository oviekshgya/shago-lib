package retry

import (
	"time"
)

type BackoffFunc func(attempt int) time.Duration

func ExponentialBackoff(initial time.Duration, multiplier float64, max time.Duration) BackoffFunc {
	return func(attempt int) time.Duration {
		d := float64(initial) * float64(mathPow(multiplier, attempt))
		if d > float64(max) {
			return max
		}
		return time.Duration(d)
	}
}

// Simple Pow helper since generic Pow requires type casting
func mathPow(x float64, n int) float64 {
	res := 1.0
	for i := 0; i < n; i++ {
		res *= x
	}
	return res
}

func Do(fn func() error, backoff BackoffFunc, maxAttempts int) error {
	var err error
	for i := 0; i < maxAttempts; i++ {
		if i > 0 {
			time.Sleep(backoff(i))
		}
		if err = fn(); err == nil {
			return nil
		}
	}
	return err
}
