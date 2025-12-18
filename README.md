# Shago Lib (DevKit)

`shago-lib` is a standardized Go development kit providing essential, production-ready utilities for building robust services. It promotes code reuse and consistency across Shago microservices.

## 🚀 Modules

### 1. Currency (`shago-lib/currency`)
Advanced Rupiah implementation.
-   **Features**: Formatting (`Rp1.500.000`) and Parsing (`Rp 1.500` -> float).
-   **Usage**: E-commerce, Financial reports.

### 2. Logger (`shago-lib/logger`)
Structured logging wrapper around `uber-go/zap`.
-   **Features**: JSON/Console encoding, Context hooks, Standard levels.
-   **Config**: `LOG_LEVEL` env var.

### 3. Slow Query (`shago-lib/slowquery`)
Performance monitoring utility.
-   **Features**: Auto-detects functions exceeding a duration threshold.
-   **Config**: `SLOW_QUERY_THRESHOLD` env var.

### 4. Benchmark (`shago-lib/benchmark`)
Simple execution timer.
-   **Features**: `Start()` / `End()` blocks for quick profiling.

### 5. Crypto (`shago-lib/crypto`)
Security helpers.
-   **Features**: AES-GCM Encryption/Decryption.

### 6. Network (`shago-lib/network`)
System IP/Mac utilities.

### 7. Retry (`shago-lib/retry`)
Resaliency helpers.
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

## 📄 License
MIT License - see [LICENSE](LICENSE) file.
