package radar

import (
	"context"
	"time"

	"github.com/e6a5/radar/radar/network"
	"github.com/e6a5/radar/radar/scanner"
	"github.com/gdamore/tcell/v2"
)

// RealDataCollector coordinates multiple focused scanners
type RealDataCollector struct {
	coordinator *scanner.Coordinator
	config      *Config
}

// NewRealDataCollector creates a new real data collector using modular scanners
func NewRealDataCollector(config *Config) *RealDataCollector {
	// Convert radar config to scanner config
	scannerConfig := &scanner.Config{
		ScanInterval:  time.Duration(config.ScanInterval * float64(time.Second)),
		MaxSignals:    config.MaxSignals,
		MaxScanRange:  config.MaxScanRange,
		UseRealData:   config.EnableRealData,
		EnableConsent: true,
	}

	coordinator := scanner.NewCoordinator(scannerConfig)

	// Add WiFi scanner (platform-specific implementation will be selected at compile time)
	if wifiScanner := createWiFiScanner(scannerConfig); wifiScanner != nil {
		coordinator.AddScanner(wifiScanner)
	}

	// Add network interface scanner (cross-platform)
	coordinator.AddScanner(network.NewInterfaceScanner(scannerConfig))

	return &RealDataCollector{
		coordinator: coordinator,
		config:      config,
	}
}

// CollectRealSignals gathers signals from all available scanners
func (rdc *RealDataCollector) CollectRealSignals() []Signal {
	ctx := context.Background()

	// Get signals from coordinator
	scannerSignals, err := rdc.coordinator.Scan(ctx)
	if err != nil || len(scannerSignals) == 0 {
		// Return basic fallback signals if real scanning fails
		return rdc.generateBasicSignals()
	}

	// Convert scanner.Signal to radar.Signal
	signals := make([]Signal, len(scannerSignals))
	for i, s := range scannerSignals {
		// Handle color conversion safely
		var color tcell.Color
		if c, ok := s.Color.(tcell.Color); ok {
			color = c
		} else {
			color = tcell.ColorWhite // default
		}

		signals[i] = Signal{
			Type:        s.Type,
			Icon:        s.Icon,
			Name:        s.Name,
			Color:       color,
			Strength:    s.Strength,
			Distance:    s.Distance,
			Angle:       s.Angle,
			Phase:       s.Phase,
			Lifetime:    s.Lifetime,
			LastSeen:    s.LastSeen,
			Persistence: s.Persistence,
			History:     convertHistory(s.History),
			MaxHistory:  s.MaxHistory,
		}
	}

	return signals
}

// GetAvailableScanners returns the names of available scanners
func (rdc *RealDataCollector) GetAvailableScanners() []string {
	return rdc.coordinator.GetScanners()
}

// convertHistory converts scanner position history to radar position history
func convertHistory(scannerHistory []scanner.PositionHistory) []PositionHistory {
	history := make([]PositionHistory, len(scannerHistory))
	for i, h := range scannerHistory {
		history[i] = PositionHistory{
			Distance:  h.Distance,
			Angle:     h.Angle,
			Strength:  h.Strength,
			Timestamp: h.Timestamp,
		}
	}
	return history
}

// generateBasicSignals creates fallback signals when real scanning fails
func (rdc *RealDataCollector) generateBasicSignals() []Signal {
	signals := make([]Signal, 0)
	now := time.Now()

	// Generate a few basic signals to show that the system is working
	basicSignals := []struct {
		name       string
		icon       string
		signalType string
		color      tcell.Color
	}{
		{"WiFi-Network", "≋", "WiFi", tcell.ColorBlue},
		{"Network-Activity", "▲", "Network", tcell.ColorGreen},
	}

	for i, basic := range basicSignals {
		signal := Signal{
			Type:        basic.signalType,
			Icon:        basic.icon,
			Name:        basic.name,
			Color:       basic.color,
			Strength:    50 + i*10,
			Distance:    float64(i + 2),
			Angle:       float64(i) * 1.57, // 90 degrees apart
			Phase:       0,
			Lifetime:    now,
			LastSeen:    now,
			Persistence: 1.0,
			History:     make([]PositionHistory, 0, 20),
			MaxHistory:  20,
		}

		signal.addToHistory(signal.Distance, signal.Angle, signal.Strength, false, now)
		signals = append(signals, signal)
	}

	return signals
}
