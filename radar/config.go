package radar

import (
	"math"
	"time"
)

type Config struct {
	RefreshRate       time.Duration
	RadarSpeed        float64
	MaxSignals        int
	SignalLifetime    time.Duration
	BeamWidth         float64
	MaxPhase          int
	SweepTrails       int     // Number of sweep trail segments
	GridSpacing       int     // Grid pattern spacing
	EnableRipples     bool    // Enable signal ripple effects
	PersistenceTime   float64 // How long signals persist after sweep (seconds)
	EnablePersistence bool    // Enable signal persistence feature
	// Signal filtering configuration
	EnableFiltering  bool // Enable signal type filtering
	ShowFilterStatus bool // Show filter status in UI
	// Signal history configuration
	EnableHistory     bool    // Enable signal history tracking
	ShowTrails        bool    // Show signal movement trails
	MaxTrailLength    int     // Maximum number of trail points to show
	HistoryUpdateRate float64 // How often to update history (seconds)
	// Signal name display configuration
	ShowSignalNames  bool // Show signal names/identifiers on radar
	ShowNamesOnHover bool // Show names only for strong signals or selected signals
	// Real data collection configuration
	EnableRealData   bool    // Enable real device data collection
	ScanInterval     float64 // How often to scan for real devices (seconds)
	UseSimulatedData bool    // Fallback to simulated data if real data fails
	MaxScanRange     float64 // Maximum simulated distance for real devices
}

// Signal type filter state
type FilterState struct {
	WiFiVisible      bool
	BluetoothVisible bool
	CellularVisible  bool
	RadioVisible     bool
	IoTVisible       bool
	SatelliteVisible bool
	AllVisible       bool // Quick toggle for all types
}

func NewConfig() Config {
	return Config{
		RefreshRate:       80 * time.Millisecond, // Faster for smoother animation
		RadarSpeed:        math.Pi / 30,
		MaxSignals:        8,
		SignalLifetime:    30 * time.Second,
		BeamWidth:         math.Pi / 60,
		MaxPhase:          6,  // More phases for smoother pulsing
		SweepTrails:       15, // Length of sweep trail
		GridSpacing:       8,  // Grid dot spacing
		EnableRipples:     true,
		PersistenceTime:   8.0, // 8 seconds to fully fade
		EnablePersistence: true,
		EnableFiltering:   true,
		ShowFilterStatus:  true,
		EnableHistory:     true,
		ShowTrails:        true,
		MaxTrailLength:    100,
		HistoryUpdateRate: 1.0,
		ShowSignalNames:   false, // Names off by default to avoid clutter
		ShowNamesOnHover:  true,  // Show names for selected/strong signals
		EnableRealData:    true,  // Real data mode by default
		ScanInterval:      2.0,   // Scan interval for real data
		UseSimulatedData:  true,  // Fallback to simulation if real data fails
		MaxScanRange:      1000.0,
	}
}

func NewFilterState() FilterState {
	return FilterState{
		WiFiVisible:      true,
		BluetoothVisible: true,
		CellularVisible:  true,
		RadioVisible:     true,
		IoTVisible:       true,
		SatelliteVisible: true,
		AllVisible:       true,
	}
}
