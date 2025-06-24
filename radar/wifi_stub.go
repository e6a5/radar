//go:build !darwin && !linux
// +build !darwin,!linux

package radar

import (
	"github.com/e6a5/radar/radar/scanner"
	"github.com/e6a5/radar/radar/wifi"
)

// createWiFiScanner creates a stub WiFi scanner for unsupported platforms
func createWiFiScanner(config *scanner.Config) scanner.Scanner {
	return wifi.NewStubWiFiScanner(config)
}
