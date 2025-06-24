//go:build !darwin && !linux
// +build !darwin,!linux

package wifi

import (
	"context"

	"github.com/e6a5/radar/radar/scanner"
)

// StubWiFiScanner is a no-op scanner for unsupported platforms
type StubWiFiScanner struct {
	config *scanner.Config
}

// NewStubWiFiScanner creates a new stub WiFi scanner
func NewStubWiFiScanner(config *scanner.Config) *StubWiFiScanner {
	return &StubWiFiScanner{
		config: config,
	}
}

// Name returns the scanner name
func (s *StubWiFiScanner) Name() string {
	return "Stub WiFi Scanner"
}

// IsAvailable returns false for unsupported platforms
func (s *StubWiFiScanner) IsAvailable() bool {
	return false
}

// Scan returns empty signals for unsupported platforms
func (s *StubWiFiScanner) Scan(ctx context.Context) ([]scanner.Signal, error) {
	return []scanner.Signal{}, nil
}
