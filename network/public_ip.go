package network

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// GetPublicIP fetches the public IP using multiple fallback services
func GetPublicIP() (string, error) {
	services := []string{
		"https://api.ipify.org?format=text",
		"https://ifconfig.me/ip",
		"https://icanhazip.com",
	}

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	for _, url := range services {
		resp, err := client.Get(url)
		if err == nil {
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)
			ip := strings.TrimSpace(string(body))
			if ip != "" {
				return ip, nil
			}
		}
	}

	return "", fmt.Errorf("failed to fetch public ip")
}
