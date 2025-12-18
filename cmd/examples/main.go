package main

import (
	"fmt"
	"os"
	"time"

	"shago-lib/currency"
	"shago-lib/logger"
	"shago-lib/slowquery"
)

func main() {
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}
	log := logger.New(logLevel)

	fmt.Println("--- Currency ---")
	rp := currency.FormatRupiah(1500000)
	fmt.Println(rp)

	fmt.Println("\n--- Logger ---")
	log.Info("Application started", logger.Field{Key: "version", Value: "1.0.0"})

	fmt.Println("\n--- Slow Query ---")
	thresholdStr := os.Getenv("SLOW_QUERY_THRESHOLD")
	threshold := 100 * time.Millisecond
	if thresholdStr != "" {
		if d, err := time.ParseDuration(thresholdStr); err == nil {
			threshold = d
		}
	}

	cfg := slowquery.Config{
		Threshold: threshold,
		OnSlow: func(op string, d time.Duration) {
			log.Warn("Slow operation detected", logger.Field{Key: "op", Value: op}, logger.Field{Key: "duration", Value: d})
		},
	}
	slowquery.TrackFunc(cfg, "heavy_computation", func() {
		time.Sleep(200 * time.Millisecond)
	})
}
