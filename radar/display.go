package radar

import (
	"math"
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v2"
)

type Display struct {
	width             int
	height            int
	centerX           int
	centerY           int
	radarAngle        float64
	signals           []Signal
	config            Config
	paused            bool
	lastUpdate        time.Time
	filters           FilterState // Signal type filtering state
	lastHistoryUpdate time.Time   // Last time we updated signal history
	// Signal selection and information panel
	selectedSignalIndex int                // Index of currently selected signal (-1 if none)
	showInfoPanel       bool               // Whether to show detailed info panel
	realDataCollector   *RealDataCollector // Add real data collector
}

func NewDisplay(width, height int) *Display {
	config := NewConfig()
	display := &Display{
		width:               width,
		height:              height,
		centerX:             width / 2,
		centerY:             height / 2,
		config:              config,
		lastUpdate:          time.Now(),
		filters:             NewFilterState(),
		lastHistoryUpdate:   time.Now(),
		selectedSignalIndex: -1, // No signal selected initially
		showInfoPanel:       false,
	}

	// Initialize real data collector with pointer to config
	display.realDataCollector = NewRealDataCollector(&display.config)

	// Generate initial signals based on configuration
	if config.EnableRealData {
		display.signals = display.realDataCollector.CollectRealSignals()
	}

	if len(display.signals) == 0 {
		display.signals = generateSignals()
	}

	return display
}

func (rd *Display) UpdatePhases() {
	if rd.paused {
		return
	}

	now := time.Now()

	// Update signal history if enough time has passed
	if rd.config.EnableHistory && now.Sub(rd.lastHistoryUpdate).Seconds() >= rd.config.HistoryUpdateRate {
		rd.updateSignalHistory(now)
		rd.lastHistoryUpdate = now
	}

	// Update signal phases and persistence
	for i := range rd.signals {
		rd.signals[i].Phase = (rd.signals[i].Phase + 1) % rd.config.MaxPhase

		// Check if signal is currently being swept by radar
		isBeingSwept := rd.angleWithinRadar(rd.signals[i].Angle)
		if isBeingSwept {
			// Signal is being swept - refresh it
			rd.signals[i].LastSeen = now
			rd.signals[i].Persistence = 1.0

			// Randomly change signal strength for realism when refreshed
			if rand.Float64() < 0.1 {
				rd.signals[i].Strength = max(10, min(100, rd.signals[i].Strength+rand.Intn(21)-10))
			}
		} else {
			// Signal is not being swept - apply persistence decay
			timeSinceLastSeen := now.Sub(rd.signals[i].LastSeen).Seconds()
			persistenceDecayRate := 1.0 / 8.0 // Takes ~8 seconds to fully fade

			rd.signals[i].Persistence = maxFloat(0.0, 1.0-timeSinceLastSeen*persistenceDecayRate)
		}
	}

	// Update radar angle
	rd.radarAngle += rd.config.RadarSpeed
	if rd.radarAngle > 2*math.Pi {
		rd.radarAngle -= 2 * math.Pi
	}

	// Remove old signals and add new ones occasionally
	if now.Sub(rd.lastUpdate) > time.Second*2 {
		rd.manageSignals(now)
		rd.lastUpdate = now
	}
}

func (rd *Display) manageSignals(now time.Time) {
	// Remove expired signals (both by lifetime and persistence)
	activeSignals := []Signal{}
	for _, s := range rd.signals {
		// Keep signal if it's within lifetime AND still has some persistence
		if now.Sub(s.Lifetime) < rd.config.SignalLifetime && s.IsVisible() {
			activeSignals = append(activeSignals, s)
		}
	}
	rd.signals = activeSignals

	// Collect real data if enabled and collector is available
	if rd.config.EnableRealData && rd.realDataCollector != nil {
		realSignals := rd.realDataCollector.CollectRealSignals()
		if len(realSignals) > 0 {
			// Replace or merge with real signals
			rd.signals = append(rd.signals, realSignals...)
			// Remove duplicates and limit to max signals
			if len(rd.signals) > rd.config.MaxSignals {
				rd.signals = rd.signals[len(rd.signals)-rd.config.MaxSignals:]
			}
		}
	}

	// Add new simulated signals occasionally if needed
	if len(rd.signals) < rd.config.MaxSignals && rand.Float64() < 0.3 {
		types := []struct {
			typeName string
			icon     string
			color    tcell.Color
		}{
			{"WiFi", "≋", tcell.ColorBlue},
			{"Bluetooth", "β", tcell.ColorNavy},
			{"Cellular", "▲", tcell.ColorGreen},
			{"Radio", "◈", tcell.ColorPurple},
			{"IoT", "◇", tcell.ColorOrange},
			{"Satellite", "★", tcell.ColorYellow},
		}

		t := types[rand.Intn(len(types))]
		distance := rand.Float64()*4 + 2
		angle := rand.Float64() * 2 * math.Pi
		strength := rand.Intn(51) + 50

		newSignal := Signal{
			Type:        t.typeName,
			Icon:        t.icon,
			Name:        "SIM-" + t.typeName, // Simple name for simulated signals
			Color:       t.color,
			Strength:    strength,
			Distance:    distance,
			Angle:       angle,
			Phase:       0,
			Lifetime:    now,
			LastSeen:    now,
			Persistence: 1.0,
			History:     make([]PositionHistory, 0, 20),
			MaxHistory:  20,
		}

		// Add initial position to history
		newSignal.addToHistory(distance, angle, strength, true, now)
		rd.signals = append(rd.signals, newSignal)
	}
}

func (rd *Display) angleWithinRadar(angle float64) bool {
	delta := math.Abs(angle - rd.radarAngle)
	if delta > math.Pi {
		delta = 2*math.Pi - delta
	}
	return delta < rd.config.BeamWidth
}

func (rd *Display) HandleInput(screen tcell.Screen) bool {
	if screen.HasPendingEvent() {
		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEscape, tcell.KeyCtrlC:
				return false // quit
			case tcell.KeyRune:
				switch ev.Rune() {
				case 'q', 'Q':
					return false
				case ' ':
					rd.paused = !rd.paused
				case '+', '=':
					rd.config.RadarSpeed = minFloat(rd.config.RadarSpeed*1.2, math.Pi/5)
				case '-', '_':
					rd.config.RadarSpeed = maxFloat(rd.config.RadarSpeed/1.2, math.Pi/120)
				case 'r', 'R':
					rd.signals = generateSignals()
					rd.radarAngle = 0
					rd.paused = false
				// Signal filtering controls
				case '1':
					rd.filters.WiFiVisible = !rd.filters.WiFiVisible
					rd.updateAllVisibleFilter()
				case '2':
					rd.filters.BluetoothVisible = !rd.filters.BluetoothVisible
					rd.updateAllVisibleFilter()
				case '3':
					rd.filters.CellularVisible = !rd.filters.CellularVisible
					rd.updateAllVisibleFilter()
				case '4':
					rd.filters.RadioVisible = !rd.filters.RadioVisible
					rd.updateAllVisibleFilter()
				case '5':
					rd.filters.IoTVisible = !rd.filters.IoTVisible
					rd.updateAllVisibleFilter()
				case '6':
					rd.filters.SatelliteVisible = !rd.filters.SatelliteVisible
					rd.updateAllVisibleFilter()
				case '0':
					// Toggle all signals
					rd.filters.AllVisible = !rd.filters.AllVisible
					rd.filters.WiFiVisible = rd.filters.AllVisible
					rd.filters.BluetoothVisible = rd.filters.AllVisible
					rd.filters.CellularVisible = rd.filters.AllVisible
					rd.filters.RadioVisible = rd.filters.AllVisible
					rd.filters.IoTVisible = rd.filters.AllVisible
					rd.filters.SatelliteVisible = rd.filters.AllVisible
				case 't', 'T':
					// Toggle signal trails
					rd.config.ShowTrails = !rd.config.ShowTrails
				case 'i', 'I':
					// Toggle information panel
					rd.showInfoPanel = !rd.showInfoPanel
				case 'n', 'N':
					// Select next signal
					rd.selectNextSignal()
				case 'p', 'P':
					// Select previous signal
					rd.selectPreviousSignal()
				case 'c', 'C':
					// Clear signal selection
					rd.selectedSignalIndex = -1
				case 's', 'S':
					// Toggle to simulation mode (temporary)
					rd.config.EnableRealData = !rd.config.EnableRealData
					if rd.config.EnableRealData && rd.realDataCollector != nil {
						// Switch back to real data
						realSignals := rd.realDataCollector.CollectRealSignals()
						if len(realSignals) > 0 {
							rd.signals = realSignals
						}
					} else {
						// Switch to simulated data temporarily
						rd.signals = generateSignals()
					}
				case 'l', 'L':
					// Toggle signal name labels
					rd.config.ShowSignalNames = !rd.config.ShowSignalNames
				}
			}
		case *tcell.EventResize:
			rd.width, rd.height = screen.Size()
			rd.centerX = rd.width / 2
			rd.centerY = rd.height / 2
		}
	}
	return true
}

// Update the AllVisible flag based on individual filter states
func (rd *Display) updateAllVisibleFilter() {
	rd.filters.AllVisible = rd.filters.WiFiVisible &&
		rd.filters.BluetoothVisible &&
		rd.filters.CellularVisible &&
		rd.filters.RadioVisible &&
		rd.filters.IoTVisible &&
		rd.filters.SatelliteVisible
}

// Check if a signal should be visible based on current filters
func (rd *Display) isSignalVisible(signal Signal) bool {
	if !rd.config.EnableFiltering {
		return true
	}

	switch signal.Type {
	case "WiFi":
		return rd.filters.WiFiVisible
	case "Bluetooth":
		return rd.filters.BluetoothVisible
	case "Cellular":
		return rd.filters.CellularVisible
	case "Radio":
		return rd.filters.RadioVisible
	case "IoT":
		return rd.filters.IoTVisible
	case "Satellite":
		return rd.filters.SatelliteVisible
	default:
		return true
	}
}

func (rd *Display) RefreshRate() time.Duration {
	return rd.config.RefreshRate
}

// Get count of visible signals
func (rd *Display) getVisibleSignalCount() int {
	count := 0
	for _, s := range rd.signals {
		if s.IsVisible() && rd.isSignalVisible(s) {
			count++
		}
	}
	return count
}

// Get signal counts by type (for display in legend)
func (rd *Display) getSignalCountsByType() map[string]int {
	counts := map[string]int{
		"WiFi":      0,
		"Bluetooth": 0,
		"Cellular":  0,
		"Radio":     0,
		"IoT":       0,
		"Satellite": 0,
	}

	for _, s := range rd.signals {
		if s.IsVisible() && rd.isSignalVisible(s) {
			counts[s.Type]++
		}
	}

	return counts
}

// Signal selection methods
func (rd *Display) selectNextSignal() {
	visibleSignals := rd.getVisibleSignalIndices()
	if len(visibleSignals) == 0 {
		rd.selectedSignalIndex = -1
		return
	}

	if rd.selectedSignalIndex == -1 {
		rd.selectedSignalIndex = visibleSignals[0]
	} else {
		// Find current signal in visible list
		currentPos := -1
		for i, idx := range visibleSignals {
			if idx == rd.selectedSignalIndex {
				currentPos = i
				break
			}
		}

		if currentPos == -1 {
			rd.selectedSignalIndex = visibleSignals[0]
		} else {
			rd.selectedSignalIndex = visibleSignals[(currentPos+1)%len(visibleSignals)]
		}
	}
}

func (rd *Display) selectPreviousSignal() {
	visibleSignals := rd.getVisibleSignalIndices()
	if len(visibleSignals) == 0 {
		rd.selectedSignalIndex = -1
		return
	}

	if rd.selectedSignalIndex == -1 {
		rd.selectedSignalIndex = visibleSignals[len(visibleSignals)-1]
	} else {
		// Find current signal in visible list
		currentPos := -1
		for i, idx := range visibleSignals {
			if idx == rd.selectedSignalIndex {
				currentPos = i
				break
			}
		}

		if currentPos == -1 {
			rd.selectedSignalIndex = visibleSignals[len(visibleSignals)-1]
		} else {
			rd.selectedSignalIndex = visibleSignals[(currentPos-1+len(visibleSignals))%len(visibleSignals)]
		}
	}
}

func (rd *Display) getVisibleSignalIndices() []int {
	var indices []int
	for i, s := range rd.signals {
		if s.IsVisible() && rd.isSignalVisible(s) {
			indices = append(indices, i)
		}
	}
	return indices
}

func (rd *Display) getSelectedSignal() *Signal {
	if rd.selectedSignalIndex >= 0 && rd.selectedSignalIndex < len(rd.signals) {
		return &rd.signals[rd.selectedSignalIndex]
	}
	return nil
}

// Update signal positions and track in history
func (rd *Display) updateSignalHistory(now time.Time) {
	for i := range rd.signals {
		// Update signal position (simulate movement)
		rd.signals[i].updatePosition(now)

		// Add current position to history
		isBeingSwept := rd.angleWithinRadar(rd.signals[i].Angle)
		rd.signals[i].addToHistory(
			rd.signals[i].Distance,
			rd.signals[i].Angle,
			rd.signals[i].Strength,
			isBeingSwept,
			now,
		)
	}
}
