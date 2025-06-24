package network

import (
	"context"
	"fmt"
	"math/rand"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/e6a5/radar/radar/scanner"
	"github.com/gdamore/tcell/v2"
)

// InterfaceScanner monitors network interfaces and active connections
type InterfaceScanner struct {
	lastScan time.Time
	config   *scanner.Config
}

// NewInterfaceScanner creates a new network interface scanner
func NewInterfaceScanner(config *scanner.Config) *InterfaceScanner {
	return &InterfaceScanner{
		config: config,
	}
}

// Name returns the scanner name
func (n *InterfaceScanner) Name() string {
	return "Network Interface Scanner"
}

// IsAvailable checks if netstat is available
func (n *InterfaceScanner) IsAvailable() bool {
	_, err := exec.LookPath("netstat")
	return err == nil
}

// Scan scans for active network connections and interfaces
func (n *InterfaceScanner) Scan(ctx context.Context) ([]scanner.Signal, error) {
	signals := make([]scanner.Signal, 0)
	now := time.Now()

	// Rate limiting
	if now.Sub(n.lastScan) < n.config.ScanInterval {
		return signals, nil
	}
	n.lastScan = now

	// Scan active connections
	connectionSignals, err := n.scanConnections(ctx, now)
	if err == nil {
		signals = append(signals, connectionSignals...)
	}

	// Scan network interfaces
	interfaceSignals, err := n.scanInterfaces(ctx, now)
	if err == nil {
		signals = append(signals, interfaceSignals...)
	}

	// Limit results
	if len(signals) > n.config.MaxSignals {
		signals = signals[:n.config.MaxSignals]
	}

	return signals, nil
}

// scanConnections scans for active network connections
func (n *InterfaceScanner) scanConnections(ctx context.Context, now time.Time) ([]scanner.Signal, error) {
	signals := make([]scanner.Signal, 0)

	cmd := exec.CommandContext(ctx, "netstat", "-n")
	output, err := cmd.Output()
	if err != nil {
		return signals, err
	}

	connectionCounts := make(map[string]int)
	lines := strings.Split(string(output), "\n")

	for _, line := range lines {
		if strings.Contains(line, "ESTABLISHED") {
			if strings.Contains(line, ":80") || strings.Contains(line, ":443") {
				connectionCounts["HTTP"]++
			} else if strings.Contains(line, ":22") {
				connectionCounts["SSH"]++
			} else if strings.Contains(line, ":53") {
				connectionCounts["DNS"]++
			} else {
				connectionCounts["Other"]++
			}
		}
	}

	// Create signals for different connection types
	connectionTypes := []struct {
		name  string
		icon  string
		color tcell.Color
	}{
		{"HTTP", "⚡", tcell.ColorGreen},
		{"SSH", "🔐", tcell.ColorYellow},
		{"DNS", "🌐", tcell.ColorBlue},
		{"Other", "▲", tcell.ColorWhite},
	}

	for _, connType := range connectionTypes {
		count := connectionCounts[connType.name]
		if count > 0 {
			signal := scanner.Signal{
				Type:        "Network",
				Icon:        connType.icon,
				Name:        fmt.Sprintf("%s (%d)", connType.name, count),
				Color:       connType.color,
				Strength:    min(100, count*20),
				Distance:    rand.Float64()*3 + 1,
				Angle:       rand.Float64() * 2 * 3.14159,
				Phase:       0,
				Lifetime:    now,
				LastSeen:    now,
				Persistence: 1.0,
				History:     make([]scanner.PositionHistory, 0, 20),
				MaxHistory:  20,
			}

			signal.AddToHistory(signal.Distance, signal.Angle, signal.Strength, true, now)
			signals = append(signals, signal)
		}
	}

	return signals, nil
}

// scanInterfaces scans network interfaces for activity
func (n *InterfaceScanner) scanInterfaces(ctx context.Context, now time.Time) ([]scanner.Signal, error) {
	signals := make([]scanner.Signal, 0)

	cmd := exec.CommandContext(ctx, "netstat", "-i")
	output, err := cmd.Output()
	if err != nil {
		return signals, err
	}

	lines := strings.Split(string(output), "\n")
	for i, line := range lines {
		if i == 0 || strings.TrimSpace(line) == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) >= 4 {
			interfaceName := fields[0]

			// Skip loopback and down interfaces
			if strings.HasPrefix(interfaceName, "lo") ||
				strings.Contains(interfaceName, "*") {
				continue
			}

			// Try to get packet counts
			var rxPackets, txPackets int
			if len(fields) >= 7 {
				if rx, err := strconv.Atoi(fields[3]); err == nil {
					rxPackets = rx
				}
				if tx, err := strconv.Atoi(fields[7]); err == nil {
					txPackets = tx
				}
			}

			// Calculate activity level
			totalPackets := rxPackets + txPackets
			if totalPackets > 0 {
				// Determine interface type
				var icon string
				var color tcell.Color
				signalType := "Network"

				if strings.HasPrefix(interfaceName, "en") ||
					strings.HasPrefix(interfaceName, "eth") {
					icon = "≋"
					color = tcell.ColorBlue
					signalType = "Ethernet"
				} else if strings.HasPrefix(interfaceName, "wl") ||
					strings.HasPrefix(interfaceName, "wifi") {
					icon = "≋"
					color = tcell.ColorGreen
					signalType = "WiFi"
				} else {
					icon = "▲"
					color = tcell.ColorWhite
				}

				// Normalize activity to strength percentage
				strength := min(100, totalPackets/1000)
				if strength < 10 {
					strength = 10
				}

				signal := scanner.Signal{
					Type:        signalType,
					Icon:        icon,
					Name:        fmt.Sprintf("%s Interface", interfaceName),
					Color:       color,
					Strength:    strength,
					Distance:    rand.Float64()*2 + 0.5,
					Angle:       rand.Float64() * 2 * 3.14159,
					Phase:       0,
					Lifetime:    now,
					LastSeen:    now,
					Persistence: 1.0,
					History:     make([]scanner.PositionHistory, 0, 20),
					MaxHistory:  20,
				}

				signal.AddToHistory(signal.Distance, signal.Angle, signal.Strength, true, now)
				signals = append(signals, signal)
			}
		}
	}

	return signals, nil
}

// min returns the smaller of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
