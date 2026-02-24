package main

import (
	"context"
	"fmt"
	"os"

	"time"

	"github.com/oviekshgya/shago-lib/benchmark"
	"github.com/oviekshgya/shago-lib/crypto"
	"github.com/oviekshgya/shago-lib/currency"
	"github.com/oviekshgya/shago-lib/logger"
	"github.com/oviekshgya/shago-lib/network"
	"github.com/oviekshgya/shago-lib/retry"
	"github.com/oviekshgya/shago-lib/slowquery"
)

func main() {
	//log := exampleLogger()
	//exampleCurrency()
	//exampleSlowQuery(log)
	//exampleBenchmark()
	//exampleCrypto()
	exampleNetwork()
	//exampleRetry(log)
}

func exampleCurrency() {
	fmt.Println("--- Currency ---")
	fmt.Printf("FormatRupiah: %s\n", currency.FormatRupiah(1500000))
	fmt.Printf(
		"FormatRupiahWithOption: %s\n",
		currency.FormatRupiahWithOption(1500000.5, currency.Option{Prefix: "Rp ", Decimal: 2}),
	)

	parsed, err := currency.ParseRupiah("Rp 1.500.000,50")
	if err != nil {
		fmt.Printf("ParseRupiah error: %v\n", err)
		return
	}
	fmt.Printf("ParseRupiah: %.2f\n", parsed)
}

func exampleLogger() logger.Logger {
	fmt.Println("\n--- Logger ---")
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	baseLog := logger.New(logLevel)
	log := baseLog.With(
		logger.Field{Key: "service", Value: "examples"},
		logger.Field{Key: "version", Value: "1.0.0"},
	)

	ctx := context.WithValue(context.Background(), "trace_id", "trace-example-001")
	log = log.WithContext(ctx)
	log.Info("Application started")
	log.Debug("Debug log enabled when LOG_LEVEL=debug")
	return log
}

func exampleSlowQuery(log logger.Logger) {
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
	fmt.Printf("Threshold: %v\n", threshold)
}

func exampleBenchmark() {
	fmt.Printf("\n--- Benchmark ---\n")
	bm := benchmark.New("simple-operation")
	time.Sleep(20 * time.Millisecond)
	bm.End()

	duration := benchmark.Track("tracked-operation", func() {
		time.Sleep(15 * time.Millisecond)
	})
	fmt.Printf("Track returned duration: %v\n", duration)

	stats := &benchmark.Stats{}
	stats.Record(duration.Microseconds())
	fmt.Println("Stats.Record called with latest duration")
}

func exampleCrypto() {
	fmt.Println("\n--- Crypto ---")
	const secret = "0123456789abcdef0123456789abcdef"

	c, err := crypto.New(crypto.Config{Secret: secret})
	if err != nil {
		fmt.Printf("crypto.New error: %v\n", err)
		return
	}

	plaintext := "hello shago"
	encrypted, err := c.Encrypt(plaintext)
	if err != nil {
		fmt.Printf("Encrypt error: %v\n", err)
		return
	}

	decrypted, err := c.Decrypt(encrypted)
	if err != nil {
		fmt.Printf("Decrypt error: %v\n", err)
		return
	}

	fmt.Printf("Plaintext: %s\n", plaintext)
	fmt.Printf("Ciphertext: %s\n", encrypted)
	fmt.Printf("Decrypted: %s\n", decrypted)
}

func exampleNetwork() {
	fmt.Println("\n--- Network ---")
	ip, err := network.GetPublicIP()
	if err != nil {
		fmt.Printf("GetPublicIP error: %v\n", err)
	} else {
		fmt.Printf("Public IP: %s\n", ip)
	}

	macs, err := network.GetMacAddresses()
	if err != nil {
		fmt.Printf("GetMacAddresses error: %v\n", err)
		return
	}
	fmt.Printf("MAC Addresses: %v\n", macs)
}

func exampleRetry(log logger.Logger) {
	fmt.Println("\n--- Retry ---")
	attempt := 0
	backoff := retry.ExponentialBackoff(50*time.Millisecond, 2, 300*time.Millisecond)
	err := retry.Do(func() error {
		attempt++
		if attempt < 3 {
			return fmt.Errorf("simulated failure at attempt %d", attempt)
		}
		return nil
	}, backoff, 5)

	if err != nil {
		log.Error("Retry failed", err, logger.Field{Key: "attempts", Value: attempt})
		return
	}

	log.Info("Retry succeeded", logger.Field{Key: "attempts", Value: attempt})
}
