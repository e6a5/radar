package radar

import (
	"math"
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v2"
)

// Historical position point for signal trails
type PositionHistory struct {
	Distance    float64
	Angle       float64
	Timestamp   time.Time
	Strength    int
	WasDetected bool // Was the signal detected at this position
}

type Signal struct {
	Type        string
	Icon        string
	Name        string    // Signal identifier/name (e.g., WiFi SSID, device name)
	Strength    int       // 0–100
	Distance    float64   // 0–10 meters
	Angle       float64   // radians
	Phase       int       // for animation (wave ring phase)
	Lifetime    time.Time // when signal was created
	Color       tcell.Color // Signal's base color
	LastSeen    time.Time // when signal was last detected by radar sweep
	Persistence float64   // how long signal stays visible after last sweep (0.0-1.0)
	// New history tracking
	History     []PositionHistory // Track signal positions over time
	MaxHistory  int               // Maximum number of history points to keep
}

func generateSignals() []Signal {
	types := []struct {
		typeName string
		icon     string
		color    tcell.Color
		names    []string
	}{
		{"WiFi", "≋", tcell.ColorBlue, []string{"MyWiFi_5G", "NETGEAR_2.4G", "Linksys_AC", "TP-Link_Guest"}},
		{"Bluetooth", "β", tcell.ColorNavy, []string{"iPhone-12", "AirPods-Pro", "MacBook", "Xbox-Controller"}},
		{"Cellular", "▲", tcell.ColorGreen, []string{"Verizon-LTE", "AT&T-5G", "T-Mobile", "Cell-Tower-1"}},
		{"Radio", "◈", tcell.ColorPurple, []string{"FM-101.5", "AM-680", "HAM-Radio", "Emergency-Freq"}},
		{"IoT", "◇", tcell.ColorOrange, []string{"Smart-TV", "Nest-Cam", "Ring-Door", "Alexa-Echo"}},
		{"Satellite", "★", tcell.ColorYellow, []string{"GPS-III", "Starlink", "ISS", "Weather-Sat"}},
	}

	signals := []Signal{}
	now := time.Now()
	
	// Generate initial set of diverse signals
	for i, t := range types {
		if i < 4 || rand.Float64() < 0.7 { // Always include first 4, 70% chance for others
			distance := rand.Float64()*4 + 2
			angle := rand.Float64() * 2 * math.Pi
			strength := rand.Intn(51) + 50
			
			// Pick a random name from the type's name list
			signalName := t.names[rand.Intn(len(t.names))]
			
			s := Signal{
				Type:        t.typeName,
				Icon:        t.icon,
				Name:        signalName,
				Color:       t.color,
				Strength:    strength,
				Distance:    distance,
				Angle:       angle,
				Phase:       rand.Intn(4),
				Lifetime:    now,
				LastSeen:    now, // Initially "seen"
				Persistence: 1.0, // Full brightness initially
				History:     make([]PositionHistory, 0, 20), // Pre-allocate for 20 positions
				MaxHistory:  20,  // Keep last 20 positions (about 40 seconds of history)
			}
			
			// Add initial position to history
			s.addToHistory(distance, angle, strength, true, now)
			
			signals = append(signals, s)
		}
	}
	
	return signals
}

// Add current position to signal history
func (s *Signal) addToHistory(distance, angle float64, strength int, wasDetected bool, timestamp time.Time) {
	newPos := PositionHistory{
		Distance:    distance,
		Angle:       angle,
		Timestamp:   timestamp,
		Strength:    strength,
		WasDetected: wasDetected,
	}
	
	s.History = append(s.History, newPos)
	
	// Keep only recent history
	if len(s.History) > s.MaxHistory {
		s.History = s.History[1:] // Remove oldest entry
	}
}

// Update signal position and track in history
func (s *Signal) updatePosition(now time.Time) {
	// Simulate realistic signal movement
	switch s.Type {
	case "WiFi", "Radio", "IoT":
		// These are typically stationary with minor fluctuations
		if rand.Float64() < 0.05 { // 5% chance to move slightly
			s.Distance += (rand.Float64() - 0.5) * 0.2 // Small distance change
			s.Angle += (rand.Float64() - 0.5) * 0.1   // Small angle change
		}
	case "Bluetooth":
		// Mobile devices - moderate movement
		if rand.Float64() < 0.15 { // 15% chance to move
			s.Distance += (rand.Float64() - 0.5) * 0.5
			s.Angle += (rand.Float64() - 0.5) * 0.2
		}
	case "Cellular":
		// Mobile phones - more movement
		if rand.Float64() < 0.25 { // 25% chance to move
			s.Distance += (rand.Float64() - 0.5) * 0.8
			s.Angle += (rand.Float64() - 0.5) * 0.3
		}
	case "Satellite":
		// Satellites move in predictable patterns
		if rand.Float64() < 0.20 { // 20% chance to move
			s.Angle += 0.05 // Consistent orbital movement
			s.Distance += (rand.Float64() - 0.5) * 0.3
		}
	}
	
	// Keep signals within reasonable bounds
	s.Distance = math.Max(1.0, math.Min(9.0, s.Distance))
	
	// Normalize angle
	for s.Angle < 0 {
		s.Angle += 2 * math.Pi
	}
	for s.Angle >= 2*math.Pi {
		s.Angle -= 2 * math.Pi
	}
}

func getColorByStrength(strength int) tcell.Color {
	switch {
	case strength > 80:
		return tcell.ColorRed
	case strength > 60:
		return tcell.ColorOrange
	case strength > 40:
		return tcell.ColorYellow
	default:
		return tcell.ColorGray
	}
}

// Get enhanced color that combines base signal color with strength and persistence
func (s *Signal) GetEnhancedColor() tcell.Color {
	baseColor := s.Color
	
	// Apply persistence fading
	if s.Persistence < 0.3 {
		return tcell.ColorDarkSlateGray // Very faded
	} else if s.Persistence < 0.6 {
		return tcell.ColorGray // Moderately faded
	}
	
	// Modify based on strength for fresh signals
	switch {
	case s.Strength > 80:
		return baseColor // Keep original color for strong signals
	case s.Strength > 60:
		return baseColor // Still original but might be styled differently
	case s.Strength > 40:
		return tcell.ColorYellow // Weaker signals become yellowish
	default:
		return tcell.ColorGray // Very weak signals are gray
	}
}

// Check if signal should be visible based on persistence
func (s *Signal) IsVisible() bool {
	return s.Persistence > 0.1
}

// Get the visual style based on persistence level
func (s *Signal) GetVisualStyle(baseStyle tcell.Style) tcell.Style {
	if s.Persistence > 0.8 {
		return baseStyle.Bold(true) // Fresh signal - bold
	} else if s.Persistence > 0.5 {
		return baseStyle // Normal intensity
	} else {
		return baseStyle.Dim(true) // Fading signal - dim
	}
} 