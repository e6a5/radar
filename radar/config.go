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
	// Performance optimization settings
	EnableVSync          bool    // Enable vertical sync for smoother rendering
	ReducedMotion        bool    // Reduce animations for better performance
	AdaptiveRefreshRate  bool    // Automatically adjust refresh rate based on performance
	MaxRenderRadius      float64 // Maximum radius to render (optimization)
	EnableSpatialCaching bool    // Cache spatial calculations
	// Visual enhancement settings
	EnableZoom bool    // Enable zoom functionality
	ZoomLevel  float64 // Current zoom level (1.0 = normal)
	MinZoom    float64 // Minimum zoom level
	MaxZoom    float64 // Maximum zoom level
	EnablePan  bool    // Enable pan functionality
	PanX       float64 // Pan offset X
	PanY       float64 // Pan offset Y
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
		RefreshRate:       80 * time.Millisecond, // Optimized: was 150ms
		RadarSpeed:        math.Pi / 30,
		MaxSignals:        8, // Increased from 5 for more activity
		SignalLifetime:    30 * time.Second,
		BeamWidth:         math.Pi / 60,
		MaxPhase:          8,  // Increased for smoother animations
		SweepTrails:       6,  // Reduced from 8 for better performance
		GridSpacing:       12, // Increased spacing to reduce grid density
		EnableRipples:     true,
		PersistenceTime:   8.0,
		EnablePersistence: true,
		EnableFiltering:   true,
		ShowFilterStatus:  true,
		EnableHistory:     true,
		ShowTrails:        true,
		MaxTrailLength:    50,  // Reduced from 100 for performance
		HistoryUpdateRate: 0.5, // Faster updates for smoother trails
		ShowSignalNames:   false,
		ShowNamesOnHover:  true,
		EnableRealData:    true,
		ScanInterval:      8.0, // Faster scanning for more responsive updates
		UseSimulatedData:  true,
		MaxScanRange:      1000.0,
		// Performance optimizations
		EnableVSync:          true,
		ReducedMotion:        false,
		AdaptiveRefreshRate:  true,
		MaxRenderRadius:      500.0,
		EnableSpatialCaching: true,
		// Visual enhancements
		EnableZoom: true,
		ZoomLevel:  1.0,
		MinZoom:    0.5,
		MaxZoom:    3.0,
		EnablePan:  true,
		PanX:       0.0,
		PanY:       0.0,
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
