package radar

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
)

func (rd *Display) Render(screen tcell.Screen) {
	screen.Clear()
	
	// Draw background grid pattern
	rd.drawBackground(screen)
	
	// Draw range rings
	rd.drawRangeRings(screen)
	
	// Draw subtle radar sweep trail (more transparent)
	rd.drawRadarSweep(screen)
	
	// Draw center point with crosshairs
	rd.drawCenter(screen)
	
	// Draw signals with enhanced visuals (ON TOP of sweep)
	rd.drawSignals(screen)
	
	// Draw enhanced UI panels
	rd.drawUI(screen)
	
	// Draw information panel if enabled and signal selected
	if rd.showInfoPanel && rd.getSelectedSignal() != nil {
		rd.drawInfoPanel(screen)
	}
	
	screen.Show()
}

func (rd *Display) drawBackground(screen tcell.Screen) {
	// Create a subtle grid pattern
	for y := 0; y < rd.height; y++ {
		for x := 0; x < rd.width; x++ {
			if (x-rd.centerX)%8 == 0 || (y-rd.centerY)%4 == 0 {
				if x < rd.width-1 && y > 2 && y < rd.height-3 {
					screen.SetContent(x, y, '·', nil, tcell.StyleDefault.Foreground(tcell.ColorDarkSlateGray))
				}
			}
		}
	}
}

func (rd *Display) drawRangeRings(screen tcell.Screen) {
	// Make range rings much larger - use most of the available space
	maxRadius := float64(min(rd.width, rd.height)) / 2.1
	
	// Draw concentric range rings with better visibility
	for ring := 1; ring <= 4; ring++ {
		radius := maxRadius * float64(ring) / 4.0
		// Use brighter characters and colors for better visibility
		ringChar := '○'
		ringColor := tcell.ColorGreen
		
		if ring%2 == 0 {
			ringChar = '●' // Alternate between filled and empty circles
			ringColor = tcell.ColorDarkGreen
		}
		
		rd.drawCircle(screen, rd.centerX, rd.centerY, radius, ringChar, ringColor)
		
		// Add range labels with better positioning
		if ring*10 <= 40 {
			label := fmt.Sprintf("%dm", ring*10)
			labelX := rd.centerX + int(radius) - len(label)/2
			labelY := rd.centerY - 1
			
			// Try multiple label positions for better visibility
			if labelX > 0 && labelX < rd.width-len(label) && labelY > 2 {
				for i, r := range label {
					screen.SetContent(labelX+i, labelY, r, nil, 
						tcell.StyleDefault.Foreground(tcell.ColorGreen).Bold(true))
				}
			} else {
				// Alternative position: bottom of ring
				labelY = rd.centerY + int(radius*0.5) + 1
				if labelY < rd.height-3 && labelX > 0 && labelX < rd.width-len(label) {
					for i, r := range label {
						screen.SetContent(labelX+i, labelY, r, nil, 
							tcell.StyleDefault.Foreground(tcell.ColorGreen).Bold(true))
					}
				}
			}
		}
	}
}

func (rd *Display) drawCircle(screen tcell.Screen, centerX, centerY int, radius float64, char rune, color tcell.Color) {
	// Use smaller angle increment for smoother, more visible circles
	for angle := 0.0; angle < 2*math.Pi; angle += 0.05 {
		x := centerX + int(radius*math.Cos(angle))
		y := centerY + int(radius*math.Sin(angle)*0.5) // Compress vertically for terminal aspect ratio
		if x >= 0 && x < rd.width && y >= 3 && y < rd.height-3 {
			screen.SetContent(x, y, char, nil, tcell.StyleDefault.Foreground(color).Bold(true))
		}
	}
}

func (rd *Display) drawRadarSweep(screen tcell.Screen) {
	maxRadius := float64(min(rd.width, rd.height)) / 2.1
	
	// Draw more subtle fading sweep trail
	for i := 0; i < 12; i++ { // Reduced from 15 to 12 for less clutter
		sweepAngle := rd.radarAngle - float64(i)*0.08 // Slightly tighter spacing
		intensity := 12 - i
		
		var color tcell.Color
		var char rune
		
		switch {
		case intensity > 10:
			color = tcell.ColorLime
			char = '│' // Thinner main beam
		case intensity > 7:
			color = tcell.ColorGreen
			char = '│'
		case intensity > 4:
			color = tcell.ColorDarkGreen
			char = '│'
		case intensity > 2:
			color = tcell.ColorDarkSlateGray
			char = '·' // Much lighter trail
		default:
			continue // Skip the faintest trails to reduce clutter
		}
		
		// Draw sweep line with reduced density
		for r := 8.0; r < maxRadius; r += 2.0 { // Increased step size for less density
			dx := int(math.Round(math.Cos(sweepAngle) * r))
			dy := int(math.Round(math.Sin(sweepAngle) * r * 0.5))
			x := rd.centerX + dx
			y := rd.centerY + dy
			
			if x >= 0 && x < rd.width && y >= 3 && y < rd.height-3 {
				// Check if there's a signal nearby - if so, make sweep even more subtle
				hasNearbySignal := rd.hasSignalNear(x, y, 2)
				
				if hasNearbySignal {
					if intensity <= 4 {
						continue // Skip faint sweep near signals
					}
					if intensity > 7 {
						char = '·' // Make bright sweep more subtle near signals
					}
				}
				
				screen.SetContent(x, y, char, nil, tcell.StyleDefault.Foreground(color))
			}
		}
	}
}

// Helper function to check if there's a signal nearby
func (rd *Display) hasSignalNear(x, y, radius int) bool {
	maxRadiusDisplay := float64(min(rd.width, rd.height)) / 2.1
	
	for _, s := range rd.signals {
		if !rd.angleWithinRadar(s.Angle) {
			continue
		}
		
		phaseOffset := float64(s.Phase) * 0.3
		distance := s.Distance + phaseOffset
		scaleFactor := maxRadiusDisplay / 10.0
		
		signalX := rd.centerX + int(math.Round(math.Cos(s.Angle) * distance * scaleFactor))
		signalY := rd.centerY + int(math.Round(math.Sin(s.Angle) * distance * scaleFactor * 0.5))
		
		// Check distance
		dx := x - signalX
		dy := y - signalY
		if dx*dx + dy*dy <= radius*radius {
			return true
		}
	}
	return false
}

func (rd *Display) drawCenter(screen tcell.Screen) {
	// Draw center crosshairs
	screen.SetContent(rd.centerX, rd.centerY, '⊕', nil, 
		tcell.StyleDefault.Foreground(tcell.ColorYellow).Bold(true))
	
	// Draw crosshair lines
	for i := -3; i <= 3; i++ {
		if i != 0 {
			// Horizontal line
			if rd.centerX+i >= 0 && rd.centerX+i < rd.width {
				screen.SetContent(rd.centerX+i, rd.centerY, '─', nil, 
					tcell.StyleDefault.Foreground(tcell.ColorYellow))
			}
			// Vertical line
			if rd.centerY+i >= 3 && rd.centerY+i < rd.height-3 {
				screen.SetContent(rd.centerX, rd.centerY+i, '│', nil, 
					tcell.StyleDefault.Foreground(tcell.ColorYellow))
			}
		}
	}
}

func (rd *Display) drawSignals(screen tcell.Screen) {
	maxRadius := float64(min(rd.width, rd.height)) / 2.1
	
	// First draw signal trails if enabled
	if rd.config.ShowTrails {
		rd.drawSignalTrails(screen, maxRadius)
	}
	
	// Then draw all visible signals (including persistent ones)
	for i, s := range rd.signals {
		if !s.IsVisible() {
			continue // Skip completely faded signals
		}
		
		// Apply signal type filtering
		if !rd.isSignalVisible(s) {
			continue // Skip filtered out signals
		}
		
		phaseOffset := float64(s.Phase) * 0.3
		distance := s.Distance + phaseOffset
		
		// Scale distance based on terminal size
		scaleFactor := maxRadius / 10.0
		dx := int(math.Round(math.Cos(s.Angle) * distance * scaleFactor))
		dy := int(math.Round(math.Sin(s.Angle) * distance * scaleFactor * 0.5))
		x := rd.centerX + dx
		y := rd.centerY + dy
		
		if x >= 0 && x < rd.width && y >= 3 && y < rd.height-3 {
			// Use enhanced signal color that combines type, strength and persistence
			color := s.GetEnhancedColor()
			
			// Use the signal's predefined icon
			icon := []rune(s.Icon)[0]
			
			// Create base style with persistence-based styling
			baseStyle := tcell.StyleDefault.Foreground(color)
			style := s.GetVisualStyle(baseStyle)
			
			// Check if this signal is selected
			isSelected := (i == rd.selectedSignalIndex)
			
			// Additional effects for signals currently being swept
			isBeingSwept := rd.angleWithinRadar(s.Angle)
			if isBeingSwept {
				// Add pulsing effect for currently swept signals
				if s.Phase%2 == 0 {
					style = style.Reverse(true)
				}
				// Clear area around active signals for better visibility
				rd.clearSignalArea(screen, x, y)
			}
			
			// Highlight selected signal
			if isSelected {
				// Draw selection indicator around signal
				rd.drawSelectionIndicator(screen, x, y)
				// Make selected signal more prominent
				style = style.Bold(true).Background(tcell.ColorDarkBlue)
			}
			
			// Draw the main signal
			screen.SetContent(x, y, icon, nil, style)
			
			// Draw signal ripples only for strong, currently swept signals
			if s.Strength > 70 && rd.config.EnableRipples && isBeingSwept && s.Persistence > 0.8 {
				rd.drawSignalRipples(screen, x, y, s.Phase, color)
			}
			
			// Add signal info for very strong signals, selected signals, or when names are enabled
			if (s.Strength > 85 && isBeingSwept) || isSelected || rd.config.ShowSignalNames {
				rd.drawSignalInfo(screen, x, y, s)
			}
		}
	}
}

// Draw signal trails showing movement history
func (rd *Display) drawSignalTrails(screen tcell.Screen, maxRadius float64) {
	now := time.Now()
	scaleFactor := maxRadius / 10.0
	
	for _, s := range rd.signals {
		// Skip if signal is filtered out or not visible
		if !s.IsVisible() || !rd.isSignalVisible(s) {
			continue
		}
		
		// Only show trails for signals with movement history
		if len(s.History) < 2 {
			continue
		}
		
		// Draw trail points (older to newer)
		maxTrailPoints := min(len(s.History)-1, rd.config.MaxTrailLength)
		for i := len(s.History) - maxTrailPoints; i < len(s.History)-1; i++ {
			if i < 0 {
				continue
			}
			
			pos := s.History[i]
			age := now.Sub(pos.Timestamp).Seconds()
			
			// Skip very old positions
			if age > 30.0 { // 30 seconds max trail age
				continue
			}
			
			// Calculate screen position
			dx := int(math.Round(math.Cos(pos.Angle) * pos.Distance * scaleFactor))
			dy := int(math.Round(math.Sin(pos.Angle) * pos.Distance * scaleFactor * 0.5))
			x := rd.centerX + dx
			y := rd.centerY + dy
			
			if x >= 0 && x < rd.width && y >= 3 && y < rd.height-3 {
				// Calculate trail intensity based on age and whether it was detected
				intensity := 1.0 - (age / 30.0) // Fade over 30 seconds
				
				var color tcell.Color
				var char rune
				
				if pos.WasDetected {
					// Detected positions use signal type color but faded
					color = s.Color
					char = '•'
				} else {
					// Undetected/estimated positions are gray
					color = tcell.ColorGray
					char = '·'
				}
				
				// Apply intensity fading
				if intensity > 0.7 {
					// Recent trail points - bright
					style := tcell.StyleDefault.Foreground(color)
					screen.SetContent(x, y, char, nil, style)
				} else if intensity > 0.4 {
					// Medium age - dim
					style := tcell.StyleDefault.Foreground(color).Dim(true)
					screen.SetContent(x, y, char, nil, style)
				} else if intensity > 0.2 {
					// Old trail points - very dim
					style := tcell.StyleDefault.Foreground(tcell.ColorDarkGray)
					screen.SetContent(x, y, '·', nil, style)
				}
				// Very old points (intensity <= 0.2) are not drawn
			}
		}
	}
}

// Clear area around signal to ensure visibility
func (rd *Display) clearSignalArea(screen tcell.Screen, centerX, centerY int) {
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			x := centerX + dx
			y := centerY + dy
			if x >= 0 && x < rd.width && y >= 3 && y < rd.height-3 {
				// Only clear if it's not the center position
				if dx != 0 || dy != 0 {
					screen.SetContent(x, y, ' ', nil, tcell.StyleDefault)
				}
			}
		}
	}
}

// Draw signal information for very strong signals
func (rd *Display) drawSignalInfo(screen tcell.Screen, x, y int, signal Signal) {
	if rd.width < 80 { // Skip on narrow screens
		return
	}
	
	// Show signal name if enabled, otherwise show signal type
	var label string
	if rd.config.ShowSignalNames {
		// Truncate long names to fit on screen
		if len(signal.Name) > 12 {
			label = signal.Name[:12] + "..."
		} else {
			label = signal.Name
		}
	} else {
		label = signal.Type[:1] // First letter of type
	}
	
	if y > 4 {
		// Draw name/type above the signal
		for i, r := range label {
			if x-len(label)/2+i >= 0 && x-len(label)/2+i < rd.width {
				screen.SetContent(x-len(label)/2+i, y-1, r, nil, 
					tcell.StyleDefault.Foreground(tcell.ColorWhite).Dim(true))
			}
		}
	}
}

func (rd *Display) drawSignalRipples(screen tcell.Screen, centerX, centerY, phase int, color tcell.Color) {
	rippleRadius := float64(phase%3 + 1)
	
	for angle := 0.0; angle < 2*math.Pi; angle += math.Pi/4 {
		x := centerX + int(rippleRadius*math.Cos(angle))
		y := centerY + int(rippleRadius*math.Sin(angle)*0.5)
		
		if x >= 0 && x < rd.width && y >= 3 && y < rd.height-3 {
			screen.SetContent(x, y, '·', nil, tcell.StyleDefault.Foreground(color))
		}
	}
}

func (rd *Display) drawUI(screen tcell.Screen) {
	rd.drawTopPanel(screen)
	rd.drawBottomPanel(screen)
	rd.drawSidePanel(screen)
}

func (rd *Display) drawTopPanel(screen tcell.Screen) {
	// Top border
	for x := 0; x < rd.width; x++ {
		screen.SetContent(x, 0, '═', nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
		screen.SetContent(x, 2, '═', nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
	}
	
	// Title and status
	title := "🌊 RADAR TERMINAL v2.0"
	if rd.paused {
		title += " [PAUSED]"
	}
	
	for i, r := range title {
		if i < rd.width {
			style := tcell.StyleDefault.Foreground(tcell.ColorAqua).Bold(true)
			if rd.paused && i >= len("🌊 RADAR TERMINAL v2.0") {
				style = tcell.StyleDefault.Foreground(tcell.ColorRed).Bold(true)
			}
			screen.SetContent(i, 1, r, nil, style)
		}
	}
	
	// Signal count, radar speed, trail status, real data status, and selection info
	trailStatus := ""
	if rd.config.ShowTrails {
		trailStatus = " | TRAILS"
	}
	
	labelStatus := ""
	if rd.config.ShowSignalNames {
		labelStatus = " | LABELS"
	}
	
	dataStatus := ""
	if rd.config.EnableRealData {
		dataStatus = " | REAL"
	} else {
		dataStatus = " | SIM"
	}
	
	selectionStatus := ""
	if rd.selectedSignalIndex >= 0 {
		selectionStatus = fmt.Sprintf(" | SEL:%d", rd.selectedSignalIndex+1)
	}
	
	info := fmt.Sprintf("Signals: %d | Speed: %.1fx%s%s%s%s", rd.getVisibleSignalCount(), rd.config.RadarSpeed/(math.Pi/30), trailStatus, labelStatus, dataStatus, selectionStatus)
	startX := rd.width - len(info)
	if startX > len(title)+2 {
		for i, r := range info {
			color := tcell.ColorWhite
			if rd.config.ShowTrails && strings.Contains(string(r), "TRAILS") {
				color = tcell.ColorGreen
			}
			screen.SetContent(startX+i, 1, r, nil, tcell.StyleDefault.Foreground(color))
		}
	}
}

func (rd *Display) drawBottomPanel(screen tcell.Screen) {
	bottomY := rd.height - 3
	
	// Bottom border
	for x := 0; x < rd.width; x++ {
		screen.SetContent(x, bottomY, '═', nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
		screen.SetContent(x, rd.height-1, '═', nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
	}
	
	// Controls
	controls := []string{
		"ESC/Q:Quit",
		"SPACE:Pause",
		"+/-:Speed", 
		"R:Reset",
		"1-6:Filter",
		"T:Trails",
		"L:Labels",
		"D:Data",
		"I:Info",
		"N/P:Select",
	}
	
	controlsLine := ""
	for i, ctrl := range controls {
		if i > 0 {
			controlsLine += " │ "
		}
		controlsLine += ctrl
	}
	
	startX := (rd.width - len(controlsLine)) / 2
	if startX > 0 {
		for i, r := range controlsLine {
			color := tcell.ColorWhite
			if r == '│' {
				color = tcell.ColorDarkGray
			}
			screen.SetContent(startX+i, rd.height-2, r, nil, tcell.StyleDefault.Foreground(color))
		}
	}
}

func (rd *Display) drawSidePanel(screen tcell.Screen) {
	if rd.width < 80 {
		return // Skip side panel on narrow screens
	}
	
	panelX := rd.width - 25
	
	// Side panel border
	for y := 3; y < rd.height-3; y++ {
		screen.SetContent(panelX-1, y, '│', nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
	}
	
	// Signal legend with filtering
	legendY := 4
	legend := "SIGNAL TYPES:"
	for i, r := range legend {
		screen.SetContent(panelX+i, legendY, r, nil, tcell.StyleDefault.Foreground(tcell.ColorAqua).Bold(true))
	}
	
	signalTypes := []struct {
		name string
		icon rune
		color tcell.Color
		key string
		visible bool
	}{
		{"WiFi", '≋', tcell.ColorBlue, "1", rd.filters.WiFiVisible},
		{"Bluetooth", 'β', tcell.ColorNavy, "2", rd.filters.BluetoothVisible},
		{"Cellular", '▲', tcell.ColorGreen, "3", rd.filters.CellularVisible},
		{"Radio", '◈', tcell.ColorPurple, "4", rd.filters.RadioVisible},
		{"IoT", '◇', tcell.ColorOrange, "5", rd.filters.IoTVisible},
		{"Satellite", '★', tcell.ColorYellow, "6", rd.filters.SatelliteVisible},
	}
	
	counts := rd.getSignalCountsByType()
	
	for i, sig := range signalTypes {
		y := legendY + 2 + i
		if y < rd.height-4 {
			// Show filter key
			keyStyle := tcell.StyleDefault.Foreground(tcell.ColorGray)
			if sig.visible {
				keyStyle = tcell.StyleDefault.Foreground(tcell.ColorWhite).Bold(true)
			}
			screen.SetContent(panelX, y, []rune(sig.key)[0], nil, keyStyle)
			
			// Show signal icon (dimmed if filtered out)
			iconStyle := tcell.StyleDefault.Foreground(sig.color).Bold(true)
			if !sig.visible {
				iconStyle = tcell.StyleDefault.Foreground(tcell.ColorDarkGray)
			}
			screen.SetContent(panelX+1, y, sig.icon, nil, iconStyle)
			
			// Show signal name
			nameStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite)
			if !sig.visible {
				nameStyle = tcell.StyleDefault.Foreground(tcell.ColorGray)
			}
			for j, r := range sig.name {
				screen.SetContent(panelX+3+j, y, r, nil, nameStyle)
			}
			
			// Show count
			count := counts[sig.name]
			countStr := fmt.Sprintf("(%d)", count)
			for j, r := range countStr {
				screen.SetContent(panelX+13+j, y, r, nil, nameStyle)
			}
		}
	}
	
	// Signal strength legend
	strengthY := legendY + 10
	if strengthY < rd.height-8 {
		strengthTitle := "STRENGTH:"
		for i, r := range strengthTitle {
			screen.SetContent(panelX+i, strengthY, r, nil, tcell.StyleDefault.Foreground(tcell.ColorAqua).Bold(true))
		}
		
		strengths := []struct {
			label string
			color tcell.Color
		}{
			{"Strong", tcell.ColorRed},
			{"Good", tcell.ColorOrange},
			{"Medium", tcell.ColorYellow},
			{"Weak", tcell.ColorGray},
		}
		
		for i, str := range strengths {
			y := strengthY + 2 + i
			if y < rd.height-4 {
				screen.SetContent(panelX, y, '●', nil, tcell.StyleDefault.Foreground(str.color).Bold(true))
				for j, r := range str.label {
					screen.SetContent(panelX+2+j, y, r, nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
				}
			}
		}
	}
}

// Draw selection indicator around selected signal
func (rd *Display) drawSelectionIndicator(screen tcell.Screen, centerX, centerY int) {
	// Draw a selection box around the signal
	positions := []struct{ dx, dy int }{
		{-1, -1}, {0, -1}, {1, -1},
		{-1, 0},           {1, 0},
		{-1, 1},  {0, 1},  {1, 1},
	}
	
	for _, pos := range positions {
		x, y := centerX+pos.dx, centerY+pos.dy
		if x >= 0 && x < rd.width && y >= 3 && y < rd.height-3 {
			screen.SetContent(x, y, '□', nil, 
				tcell.StyleDefault.Foreground(tcell.ColorYellow).Bold(true))
		}
	}
}

// Draw detailed information panel for selected signal
func (rd *Display) drawInfoPanel(screen tcell.Screen) {
	signal := rd.getSelectedSignal()
	if signal == nil {
		return
	}
	
	// Panel dimensions and position
	panelWidth := 40
	panelHeight := 15
	startX := rd.width - panelWidth - 2
	startY := 4
	
	// Skip if not enough space
	if startX < 0 || startY+panelHeight >= rd.height-3 {
		return
	}
	
	// Draw panel background and border
	for y := startY; y < startY+panelHeight; y++ {
		for x := startX; x < startX+panelWidth; x++ {
			if y == startY || y == startY+panelHeight-1 || 
			   x == startX || x == startX+panelWidth-1 {
				screen.SetContent(x, y, '═', nil, 
					tcell.StyleDefault.Foreground(tcell.ColorAqua))
			} else {
				screen.SetContent(x, y, ' ', nil, 
					tcell.StyleDefault.Background(tcell.ColorDarkSlateGray))
			}
		}
	}
	
	// Panel title
	title := "SIGNAL INFORMATION"
	titleX := startX + (panelWidth-len(title))/2
	for i, r := range title {
		screen.SetContent(titleX+i, startY, r, nil, 
			tcell.StyleDefault.Foreground(tcell.ColorAqua).Bold(true))
	}
	
	// Signal details
	details := []string{
		fmt.Sprintf("Name:     %s", signal.Name),
		fmt.Sprintf("Type:     %s %s", signal.Type, signal.Icon),
		fmt.Sprintf("Strength: %d%% (%s)", signal.Strength, rd.getStrengthLabel(signal.Strength)),
		fmt.Sprintf("Distance: %.1fm", signal.Distance),
		fmt.Sprintf("Bearing:  %.0f°", signal.Angle*180/math.Pi),
		fmt.Sprintf("Age:      %.0fs", time.Since(signal.Lifetime).Seconds()),
		fmt.Sprintf("Last Seen: %.1fs ago", time.Since(signal.LastSeen).Seconds()),
		fmt.Sprintf("Persist:  %.0f%%", signal.Persistence*100),
		"",
		"HISTORY:",
		fmt.Sprintf("Positions: %d/%d", len(signal.History), signal.MaxHistory),
	}
	
	// Add movement analysis
	if len(signal.History) >= 2 {
		recent := signal.History[len(signal.History)-1]
		older := signal.History[max(0, len(signal.History)-5)]
		distanceMoved := math.Sqrt(math.Pow(recent.Distance-older.Distance, 2) + 
			math.Pow(recent.Angle-older.Angle, 2))
		details = append(details, fmt.Sprintf("Movement: %.2f units", distanceMoved))
	}
	
	// Draw details
	for i, detail := range details {
		if i+2 < panelHeight-1 {
			y := startY + 2 + i
			for j, r := range detail {
				if j+2 < panelWidth-2 {
					color := tcell.ColorWhite
					if strings.Contains(detail, "HISTORY:") {
						color = tcell.ColorYellow
					}
					screen.SetContent(startX+2+j, y, r, nil, 
						tcell.StyleDefault.Foreground(color))
				}
			}
		}
	}
}

func (rd *Display) getStrengthLabel(strength int) string {
	switch {
	case strength > 80:
		return "Strong"
	case strength > 60:
		return "Good"
	case strength > 40:
		return "Medium"
	default:
		return "Weak"
	}
}

 