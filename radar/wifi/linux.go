//go:build linux
// +build linux

package wifi

import (
	"context"
	"math"
	"math/rand"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/e6a5/radar/radar/scanner"
	"github.com/gdamore/tcell/v2"
)

// LinuxWiFiScanner implements WiFi scanning for Linux systems
type LinuxWiFiScanner struct {
	lastScan time.Time
	config   *scanner.Config
}

// NewLinuxWiFiScanner creates a new Linux WiFi scanner
func NewLinuxWiFiScanner(config *scanner.Config) *LinuxWiFiScanner {
	return &LinuxWiFiScanner{
		config: config,
	}
}

// Name returns the scanner name
func (l *LinuxWiFiScanner) Name() string {
	return "Linux WiFi Scanner"
}

// IsAvailable checks if nmcli or iw is available
func (l *LinuxWiFiScanner) IsAvailable() bool {
	_, err1 := exec.LookPath("nmcli")
	_, err2 := exec.LookPath("iw")
	return err1 == nil || err2 == nil
}

// Scan scans for WiFi networks using available Linux tools
func (l *LinuxWiFiScanner) Scan(ctx context.Context) ([]scanner.Signal, error) {
	signals := make([]scanner.Signal, 0)
	now := time.Now()

	// Rate limiting
	if now.Sub(l.lastScan) < l.config.ScanInterval {
		return signals, nil
	}
	l.lastScan = now

	// Try nmcli first
	if nmcliSignals, err := l.scanWithNmcli(ctx, now); err == nil && len(nmcliSignals) > 0 {
		return nmcliSignals, nil
	}

	// Fallback to iw
	if iwSignals, err := l.scanWithIw(ctx, now); err == nil && len(iwSignals) > 0 {
		return iwSignals, nil
	}

	return signals, nil
}

// scanWithNmcli scans WiFi networks using NetworkManager
func (l *LinuxWiFiScanner) scanWithNmcli(ctx context.Context, now time.Time) ([]scanner.Signal, error) {
	signals := make([]scanner.Signal, 0)

	// Refresh scan
	exec.CommandContext(ctx, "nmcli", "dev", "wifi", "rescan").Run()

	// Get WiFi list
	cmd := exec.CommandContext(ctx, "nmcli", "dev", "wifi", "list")
	output, err := cmd.Output()
	if err != nil {
		cmd = exec.CommandContext(ctx, "nmcli", "dev", "wifi")
		output, err = cmd.Output()
		if err != nil {
			return signals, err
		}
	}

	lines := strings.Split(string(output), "\n")
	for i, line := range lines {
		if i == 0 || strings.TrimSpace(line) == "" {
			continue
		}

		if signal := l.parseNmcliLine(line, now); signal != nil {
			signals = append(signals, *signal)
			if len(signals) >= l.config.MaxSignals {
				break
			}
		}
	}

	return signals, nil
}

// scanWithIw scans WiFi networks using iw
func (l *LinuxWiFiScanner) scanWithIw(ctx context.Context, now time.Time) ([]scanner.Signal, error) {
	signals := make([]scanner.Signal, 0)

	// Find wireless interface
	interfacesCmd := exec.CommandContext(ctx, "iw", "dev")
	interfacesOutput, err := interfacesCmd.Output()
	if err != nil {
		return signals, err
	}

	var wifiInterface string
	lines := strings.Split(string(interfacesOutput), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Interface") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				wifiInterface = fields[1]
				break
			}
		}
	}

	if wifiInterface == "" {
		return signals, nil
	}

	// Scan for networks
	cmd := exec.CommandContext(ctx, "iw", wifiInterface, "scan")
	output, err := cmd.Output()
	if err != nil {
		return signals, err
	}

	lines = strings.Split(string(output), "\n")
	var currentSSID string
	var currentRSSI int

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "SSID:") {
			currentSSID = strings.TrimSpace(strings.TrimPrefix(line, "SSID:"))
		} else if strings.Contains(line, "signal:") {
			fields := strings.Fields(line)
			for i, field := range fields {
				if field == "signal:" && i+1 < len(fields) {
					rssiStr := strings.TrimSuffix(fields[i+1], ".00")
					if val, err := strconv.ParseFloat(rssiStr, 64); err == nil {
						currentRSSI = int(val)
					}
					break
				}
			}
		}

		// End of BSS entry
		if strings.HasPrefix(line, "BSS") && currentSSID != "" {
			strength := rssiToStrength(currentRSSI)
			distance := rssiToDistance(currentRSSI, l.config.MaxScanRange)

			// Get friendly display name
			displayName := GetFriendlyDisplayName(currentSSID, strength, false)

			signal := scanner.Signal{
				Type:        "WiFi",
				Icon:        "≋",
				Name:        displayName,
				Color:       tcell.ColorBlue,
				Strength:    strength,
				Distance:    distance,
				Angle:       rand.Float64() * 2 * math.Pi,
				Phase:       0,
				Lifetime:    now,
				LastSeen:    now,
				Persistence: 1.0,
				History:     make([]scanner.PositionHistory, 0, 20),
				MaxHistory:  20,
			}

			signal.AddToHistory(signal.Distance, signal.Angle, signal.Strength, true, now)
			signals = append(signals, signal)

			currentSSID = ""
			currentRSSI = -50

			if len(signals) >= l.config.MaxSignals {
				break
			}
		}
	}

	return signals, nil
}

// parseNmcliLine parses a line from nmcli output
func (l *LinuxWiFiScanner) parseNmcliLine(line string, now time.Time) *scanner.Signal {
	line = strings.TrimSpace(line)
	if line == "" || strings.HasPrefix(line, "IN-USE") || strings.HasPrefix(line, "SSID") {
		return nil
	}

	// Remove active connection indicator
	if strings.HasPrefix(line, "*") {
		line = strings.TrimSpace(line[1:])
	}

	parts := strings.Fields(line)
	if len(parts) < 4 {
		return nil
	}

	ssid := parts[0]
	if ssid == "" || ssid == "SSID" || ssid == "--" {
		return nil
	}

	// Find signal strength
	var strength int = 50
	for _, part := range parts {
		if val, err := strconv.Atoi(part); err == nil {
			if val >= 0 && val <= 100 {
				strength = val
				break
			} else if val < 0 && val >= -100 {
				strength = rssiToStrength(val)
				break
			}
		}
	}

	distance := float64(100-strength) / 15.0
	if distance < 0.5 {
		distance = 0.5
	}
	if distance > 10 {
		distance = 10
	}

	// Get friendly display name
	displayName := GetFriendlyDisplayName(ssid, strength, false)

	signal := scanner.Signal{
		Type:        "WiFi",
		Icon:        "≋",
		Name:        displayName,
		Color:       tcell.ColorBlue,
		Strength:    strength,
		Distance:    distance,
		Angle:       rand.Float64() * 2 * math.Pi,
		Phase:       0,
		Lifetime:    now,
		LastSeen:    now,
		Persistence: 1.0,
		History:     make([]scanner.PositionHistory, 0, 20),
		MaxHistory:  20,
	}

	signal.AddToHistory(signal.Distance, signal.Angle, signal.Strength, true, now)
	return &signal
}

// rssiToStrength converts RSSI to percentage
func rssiToStrength(rssi int) int {
	if rssi >= -30 {
		return 100
	}
	if rssi <= -90 {
		return 0
	}
	return int(100 * (float64(rssi+90) / 60.0))
}

// rssiToDistance converts RSSI to radar distance
func rssiToDistance(rssi int, maxRange float64) float64 {
	distance := math.Pow(10, float64(-rssi-30)/20.0) * 2.0
	if distance < 0.5 {
		distance = 0.5
	}
	if distance > maxRange {
		distance = maxRange
	}
	return distance
}
