//go:build linux
// +build linux

package radar

import (
	"github.com/e6a5/radar/radar/scanner"
	"github.com/e6a5/radar/radar/wifi"
)

// createWiFiScanner creates a WiFi scanner for Linux
func createWiFiScanner(config *scanner.Config) scanner.Scanner {
	return wifi.NewLinuxWiFiScanner(config)
}
