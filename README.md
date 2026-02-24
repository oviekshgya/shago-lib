# Shago Lib (DevKit)

`github.com/oviekshgya/shago-lib` is a standardized Go development kit providing essential, production-ready utilities for building robust services. It promotes code reuse and consistency across Shago microservices.

## 🚀 Modules

### 1. Currency (`github.com/oviekshgya/shago-lib/currency`)
Advanced Rupiah implementation.
-   **Features**: Formatting (`Rp1.500.000`) and Parsing (`Rp 1.500` -> float).
-   **Usage**: E-commerce, Financial reports.

### 2. Logger (`github.com/oviekshgya/shago-lib/logger`)
Structured logging wrapper around `uber-go/zap`.
-   **Features**: JSON/Console encoding, Context hooks, Standard levels.
-   **Config**: `LOG_LEVEL` env var.

### 3. Slow Query (`github.com/oviekshgya/shago-lib/slowquery`)
Performance monitoring utility.
-   **Features**: Auto-detects functions exceeding a duration threshold.
-   **Config**: `SLOW_QUERY_THRESHOLD` env var.

### 4. Benchmark (`github.com/oviekshgya/shago-lib/benchmark`)
Simple execution timer.
-   **Features**: `New()`/`End()`, `Track()`, and lightweight stats recorder.

### 5. Crypto (`github.com/oviekshgya/shago-lib/crypto`)
Security helpers.
-   **Features**: AES-GCM Encryption/Decryption.

### 6. Network (`github.com/oviekshgya/shago-lib/network`)
System IP/Mac utilities.

### 7. Retry (`github.com/oviekshgya/shago-lib/retry`)
Resiliency helpers.
-   **Features**: Exponential backoff.

## 📦 Installation

```bash
go get github.com/oviekshgya/shago-lib
```

## 🚀 Usage

Check `cmd/examples/main.go` for a full runner.

```bash
# Set Env config
export LOG_LEVEL=debug
export SLOW_QUERY_THRESHOLD=50ms

# Run example
go run cmd/examples/main.go
```

## 🧪 Package Examples

### Currency
```go
package main

import (
	"fmt"
	"github.com/oviekshgya/shago-lib/currency"
)

func main() {
	fmt.Println(currency.FormatRupiah(1500000)) // Rp1.500.000
	fmt.Println(currency.FormatRupiahWithOption(
		1500000.5,
		currency.Option{Prefix: "Rp ", Decimal: 2},
	)) // Rp 1.500.000,50

	amount, err := currency.ParseRupiah("Rp 1.500.000,50")
	fmt.Println(amount, err) // 1500000.5 <nil>
}
```

### Logger
```go
package main

import (
	"context"
	"github.com/oviekshgya/shago-lib/logger"
)

func main() {
	log := logger.New("info").With(
		logger.Field{Key: "service", Value: "payment"},
		logger.Field{Key: "version", Value: "1.0.0"},
	)

	ctx := context.WithValue(context.Background(), "trace_id", "trace-123")
	log.WithContext(ctx).Info("application started")
}
```

### Slow Query
```go
package main

import (
	"time"
	"github.com/oviekshgya/shago-lib/logger"
	"github.com/oviekshgya/shago-lib/slowquery"
)

func main() {
	log := logger.New("info")
	cfg := slowquery.Config{
		Threshold: 100 * time.Millisecond,
		OnSlow: func(op string, d time.Duration) {
			log.Warn("slow operation", logger.Field{Key: "op", Value: op}, logger.Field{Key: "duration", Value: d})
		},
	}

	slowquery.TrackFunc(cfg, "heavy_computation", func() {
		time.Sleep(200 * time.Millisecond)
	})
}
```

### Benchmark
```go
package main

import (
	"time"
	"github.com/oviekshgya/shago-lib/benchmark"
)

func main() {
	p := benchmark.New("load-data")
	time.Sleep(20 * time.Millisecond)
	p.End()

	d := benchmark.Track("transform-data", func() {
		time.Sleep(15 * time.Millisecond)
	})

	stats := &benchmark.Stats{}
	stats.Record(d.Microseconds())
}
```

### Crypto
```go
package main

import (
	"fmt"
	"github.com/oviekshgya/shago-lib/crypto"
)

func main() {
	c, err := crypto.New(crypto.Config{
		Secret: "0123456789abcdef0123456789abcdef", // 32 bytes
	})
	if err != nil {
		panic(err)
	}

	encrypted, _ := c.Encrypt("hello shago")
	decrypted, _ := c.Decrypt(encrypted)
	fmt.Println(encrypted)
	fmt.Println(decrypted) // hello shago
}
```

### Network
```go
package main

import (
	"fmt"
	"github.com/oviekshgya/shago-lib/network"
)

func main() {
	ip, err := network.GetPublicIP()
	fmt.Println("ip:", ip, "err:", err)

	macs, err := network.GetMacAddresses()
	fmt.Println("macs:", macs, "err:", err)
}
```

### Retry
```go
package main

import (
	"fmt"
	"time"
	"github.com/oviekshgya/shago-lib/retry"
)

func main() {
	attempt := 0
	backoff := retry.ExponentialBackoff(50*time.Millisecond, 2, 300*time.Millisecond)

	err := retry.Do(func() error {
		attempt++
		if attempt < 3 {
			return fmt.Errorf("failed at attempt %d", attempt)
		}
		return nil
	}, backoff, 5)

	fmt.Println("attempt:", attempt, "err:", err) // attempt: 3 err: <nil>
}
```

## 📄 License
MIT License - see [LICENSE](LICENSE) file.
