//go:build darwin
// +build darwin

package radar

import (
	"github.com/e6a5/radar/radar/scanner"
	"github.com/e6a5/radar/radar/wifi"
)

// createWiFiScanner creates a WiFi scanner for Darwin/macOS
func createWiFiScanner(config *scanner.Config) scanner.Scanner {
	return wifi.NewCoreWLANScanner(config)
}
