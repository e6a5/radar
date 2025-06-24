package radar

import (
	"fmt"
	"math"

	"github.com/gdamore/tcell/v2"
)

// Get current theme for display components
func (rd *Display) getCurrentTheme() RadarTheme {
	return GetModernDarkTheme() // Default theme for now
}

// Enhanced visual constants for modern radar UI
const (
	// Radar circle styles
	RangeRingPrimary   = '●'
	RangeRingSecondary = '○'
	RangeRingDotted    = '·'
	RangeRingDashed    = '-'

	// Sweeper styles
	SweeperPrimary   = '│'
	SweeperSecondary = '┃'
	SweeperFade      = '╎'
	SweeperDot       = '·'

	// Signal styles
	SignalStrong = '●'
	SignalMedium = '◉'
	SignalWeak   = '○'
	SignalPulse  = '◎'
	SignalRipple = '◯'

	// Grid patterns
	GridCross = '+'
	GridDot   = '·'
	GridTick  = '╎'
)

// Color schemes for enhanced visuals
type RadarTheme struct {
	Background    tcell.Color
	GridPrimary   tcell.Color
	GridSecondary tcell.Color

	// Range rings
	RingPrimary   tcell.Color
	RingSecondary tcell.Color
	RingLabels    tcell.Color

	// Sweeper colors
	SweepPrimary   tcell.Color
	SweepSecondary tcell.Color
	SweepFade      tcell.Color
	SweepTrail     tcell.Color

	// Signal colors by strength
	SignalExcellent tcell.Color
	SignalGood      tcell.Color
	SignalFair      tcell.Color
	SignalPoor      tcell.Color
	SignalConnected tcell.Color

	// UI accents
	AccentPrimary   tcell.Color
	AccentSecondary tcell.Color
	TextPrimary     tcell.Color
	TextSecondary   tcell.Color
}

// Modern dark theme
func GetModernDarkTheme() RadarTheme {
	return RadarTheme{
		Background:    tcell.ColorBlack,
		GridPrimary:   tcell.ColorDarkSlateGray,
		GridSecondary: tcell.Color16,

		RingPrimary:   tcell.ColorLime,
		RingSecondary: tcell.ColorGreen,
		RingLabels:    tcell.ColorYellow,

		SweepPrimary:   tcell.ColorLime,
		SweepSecondary: tcell.ColorGreen,
		SweepFade:      tcell.ColorDarkGreen,
		SweepTrail:     tcell.ColorDarkSlateGray,

		SignalExcellent: tcell.ColorRed,
		SignalGood:      tcell.ColorOrange,
		SignalFair:      tcell.ColorYellow,
		SignalPoor:      tcell.ColorGray,
		SignalConnected: tcell.ColorBlue,

		AccentPrimary:   tcell.ColorBlue,
		AccentSecondary: tcell.ColorBlue,
		TextPrimary:     tcell.ColorWhite,
		TextSecondary:   tcell.ColorGray,
	}
}

// Enhanced background grid with crosshairs and coordinate system
func (rd *Display) drawEnhancedBackground(screen tcell.Screen) {
	theme := rd.getCurrentTheme()
	maxRadius := float64(min(rd.width, rd.height)) / 2.1
	centerRadius := int(maxRadius)

	// Draw coordinate grid with modern styling
	spacing := rd.config.GridSpacing

	// Primary grid lines (every 4th line)
	for y := rd.centerY - centerRadius; y <= rd.centerY+centerRadius; y += spacing * 4 {
		if y < 3 || y >= rd.height-3 {
			continue
		}
		for x := rd.centerX - centerRadius; x <= rd.centerX+centerRadius; x++ {
			if x < 1 || x >= rd.width-1 {
				continue
			}

			dx := float64(x - rd.centerX)
			dy := float64(y-rd.centerY) * 2.0
			if dx*dx+dy*dy <= maxRadius*maxRadius {
				screen.SetContent(x, y, GridTick, nil,
					tcell.StyleDefault.Foreground(theme.GridPrimary))
			}
		}
	}

	// Secondary grid dots
	for y := rd.centerY - centerRadius; y <= rd.centerY+centerRadius; y += spacing {
		if y < 3 || y >= rd.height-3 {
			continue
		}
		for x := rd.centerX - centerRadius; x <= rd.centerX+centerRadius; x += spacing {
			if x < 1 || x >= rd.width-1 {
				continue
			}

			dx := float64(x - rd.centerX)
			dy := float64(y-rd.centerY) * 2.0
			if dx*dx+dy*dy <= maxRadius*maxRadius {
				screen.SetContent(x, y, GridDot, nil,
					tcell.StyleDefault.Foreground(theme.GridSecondary))
			}
		}
	}

	// Draw center crosshairs with enhanced styling
	rd.drawEnhancedCrosshairs(screen, theme)
}

// Enhanced crosshairs at radar center
func (rd *Display) drawEnhancedCrosshairs(screen tcell.Screen, theme RadarTheme) {
	crosshairLength := 15

	// Horizontal crosshair
	for i := -crosshairLength; i <= crosshairLength; i++ {
		x := rd.centerX + i
		if x >= 0 && x < rd.width {
			char := GridTick
			color := theme.AccentPrimary
			if i == 0 {
				char = GridCross
				color = theme.AccentSecondary
			} else if i%5 == 0 {
				char = '|'
			}
			screen.SetContent(x, rd.centerY, char, nil,
				tcell.StyleDefault.Foreground(color).Bold(true))
		}
	}

	// Vertical crosshair
	for i := -crosshairLength / 2; i <= crosshairLength/2; i++ {
		y := rd.centerY + i
		if y >= 3 && y < rd.height-3 {
			char := GridTick
			color := theme.AccentPrimary
			if i == 0 {
				char = GridCross
				color = theme.AccentSecondary
			} else if i%3 == 0 {
				char = '-'
			}
			screen.SetContent(rd.centerX, y, char, nil,
				tcell.StyleDefault.Foreground(color).Bold(true))
		}
	}
}

// Enhanced range rings with modern styling and better labels
func (rd *Display) drawEnhancedRangeRings(screen tcell.Screen) {
	theme := rd.getCurrentTheme()
	maxRadius := float64(min(rd.width, rd.height)) / 2.1

	for ring := 1; ring <= 4; ring++ {
		radius := maxRadius * float64(ring) / 4.0

		// Alternate ring styles for better visual hierarchy
		var ringChar rune
		var ringColor tcell.Color
		var bold bool

		if ring%2 == 1 {
			// Primary rings (1st, 3rd)
			ringChar = RangeRingPrimary
			ringColor = theme.RingPrimary
			bold = true
		} else {
			// Secondary rings (2nd, 4th)
			ringChar = RangeRingSecondary
			ringColor = theme.RingSecondary
			bold = false
		}

		// Draw ring with enhanced density for smoother appearance
		stepSize := 0.05 // Smaller steps for smoother circles
		for angle := 0.0; angle < 2*math.Pi; angle += stepSize {
			x := rd.centerX + int(math.Round(radius*math.Cos(angle)))
			y := rd.centerY + int(math.Round(radius*math.Sin(angle)*0.5))

			if x >= 0 && x < rd.width && y >= 3 && y < rd.height-3 {
				style := tcell.StyleDefault.Foreground(ringColor)
				if bold {
					style = style.Bold(true)
				}
				screen.SetContent(x, y, ringChar, nil, style)
			}
		}

		// Enhanced range labels with better positioning
		rd.drawEnhancedRangeLabel(screen, ring, radius, theme)
	}
}

// Enhanced range labels with modern styling
func (rd *Display) drawEnhancedRangeLabel(screen tcell.Screen, ring int, radius float64, theme RadarTheme) {
	distance := ring * 10
	label := fmt.Sprintf("%dm", distance)

	// Multiple label positions for better visibility
	positions := []struct {
		angle  float64
		offset int
	}{
		{0, 2},               // Right
		{math.Pi / 2, 0},     // Top
		{math.Pi, -2},        // Left
		{3 * math.Pi / 2, 1}, // Bottom
	}

	for _, pos := range positions {
		labelX := rd.centerX + int(radius*math.Cos(pos.angle)) + pos.offset
		labelY := rd.centerY + int(radius*math.Sin(pos.angle)*0.5)

		// Only draw if position is clear and visible
		if labelX > len(label) && labelX < rd.width-len(label) &&
			labelY > 3 && labelY < rd.height-4 {

			// Draw label background for better readability
			for i := -1; i <= len(label); i++ {
				bgX := labelX + i
				if bgX >= 0 && bgX < rd.width {
					screen.SetContent(bgX, labelY, ' ', nil,
						tcell.StyleDefault.Background(tcell.ColorDarkSlateGray))
				}
			}

			// Draw label text
			for i, r := range label {
				screen.SetContent(labelX+i, labelY, r, nil,
					tcell.StyleDefault.Foreground(theme.RingLabels).
						Background(tcell.ColorDarkSlateGray).Bold(true))
			}
			break // Only draw one label per ring to avoid clutter
		}
	}
}

// Enhanced radar sweep with sophisticated trail effect
func (rd *Display) drawEnhancedRadarSweep(screen tcell.Screen) {
	theme := rd.getCurrentTheme()
	maxRadius := float64(min(rd.width, rd.height)) / 2.1

	// Enhanced sweep trail with multiple intensity levels
	trailCount := rd.config.SweepTrails

	for i := 0; i < trailCount; i++ {
		sweepAngle := rd.radarAngle - float64(i)*0.06 // Tighter trail spacing
		intensity := float64(trailCount-i) / float64(trailCount)

		// Advanced color and character selection based on intensity
		var color tcell.Color
		var char rune
		var bold bool

		switch {
		case intensity > 0.8:
			color = theme.SweepPrimary
			char = SweeperPrimary
			bold = true
		case intensity > 0.6:
			color = theme.SweepSecondary
			char = SweeperSecondary
			bold = true
		case intensity > 0.4:
			color = theme.SweepFade
			char = SweeperFade
			bold = false
		case intensity > 0.2:
			color = theme.SweepTrail
			char = SweeperDot
			bold = false
		default:
			color = theme.GridSecondary
			char = GridDot
			bold = false
		}

		// Draw sweep line with variable density
		stepSize := 1.5
		if intensity > 0.5 {
			stepSize = 1.0 // Denser for brighter parts
		}

		for r := 8.0; r < maxRadius; r += stepSize {
			dx := int(math.Round(math.Cos(sweepAngle) * r))
			dy := int(math.Round(math.Sin(sweepAngle) * r * 0.5))

			x := rd.centerX + dx
			y := rd.centerY + dy

			if x >= 0 && x < rd.width && y >= 3 && y < rd.height-3 {
				// Add sweep glow effect for primary beam
				if intensity > 0.7 && i < 2 {
					rd.addSweepGlow(screen, x, y, theme, intensity)
				}

				style := tcell.StyleDefault.Foreground(color)
				if bold {
					style = style.Bold(true)
				}
				screen.SetContent(x, y, char, nil, style)
			}
		}
	}
}

// Add glow effect around primary sweep beam
func (rd *Display) addSweepGlow(screen tcell.Screen, centerX, centerY int, theme RadarTheme, intensity float64) {
	if intensity < 0.7 {
		return
	}

	// Add subtle glow around the main beam
	glowPositions := []struct{ dx, dy int }{
		{-1, 0}, {1, 0}, {0, -1}, {0, 1},
		{-1, -1}, {1, -1}, {-1, 1}, {1, 1},
	}

	for _, pos := range glowPositions {
		x := centerX + pos.dx
		y := centerY + pos.dy

		if x >= 0 && x < rd.width && y >= 3 && y < rd.height-3 {
			// Only add glow if the position is empty or has low-intensity content
			existing, _, _, _ := screen.GetContent(x, y)
			if existing == ' ' || existing == GridDot {
				screen.SetContent(x, y, SweeperDot, nil,
					tcell.StyleDefault.Foreground(theme.SweepFade))
			}
		}
	}
}

// Enhanced signal rendering with dynamic visual effects
func (rd *Display) drawEnhancedSignals(screen tcell.Screen) {
	theme := rd.getCurrentTheme()
	maxRadiusDisplay := float64(min(rd.width, rd.height)) / 2.1

	for i, s := range rd.signals {
		if !s.IsVisible() || !rd.isSignalVisible(s) {
			continue
		}

		// Enhanced signal positioning with zoom support
		scaleFactor := maxRadiusDisplay / rd.config.MaxScanRange * 10.0
		distance := s.Distance * scaleFactor

		if distance > maxRadiusDisplay {
			distance = maxRadiusDisplay
		}

		signalX := rd.centerX + int(math.Round(math.Cos(s.Angle)*distance))
		signalY := rd.centerY + int(math.Round(math.Sin(s.Angle)*distance*0.5))

		// Enhanced signal visualization based on strength and type
		rd.drawEnhancedSignal(screen, signalX, signalY, s, i, theme)

		// Add signal trails if enabled
		if rd.config.ShowTrails {
			rd.drawSignalTrail(screen, s, scaleFactor, theme)
		}

		// Add ripple effects for strong signals
		if s.Strength > 70 {
			rd.drawEnhancedSignalRipples(screen, signalX, signalY, s, theme)
		}
	}
}

// Enhanced individual signal rendering
func (rd *Display) drawEnhancedSignal(screen tcell.Screen, x, y int, signal Signal, index int, theme RadarTheme) {
	if x < 0 || x >= rd.width || y < 3 || y >= rd.height-3 {
		return
	}

	// Select signal character and color based on strength and type
	var char rune
	var color tcell.Color
	var bold bool

	// Strength-based visualization
	switch {
	case signal.Strength >= 80:
		char = SignalStrong
		color = theme.SignalExcellent
		bold = true
	case signal.Strength >= 60:
		char = SignalMedium
		color = theme.SignalGood
		bold = true
	case signal.Strength >= 40:
		char = SignalWeak
		color = theme.SignalFair
		bold = false
	default:
		char = SignalRipple
		color = theme.SignalPoor
		bold = false
	}

	// Special styling for connected signals
	if signal.Type == "WiFi" && signal.Strength > 70 {
		color = theme.SignalConnected
		char = SignalPulse
	}

	// Pulsing effect based on phase
	if signal.Phase > rd.config.MaxPhase/2 {
		char = SignalPulse
		bold = true
	}

	// Selection highlighting
	if index == rd.selectedSignalIndex {
		// Draw selection border
		rd.drawSelectionBorder(screen, x, y, theme)
		bold = true
	}

	style := tcell.StyleDefault.Foreground(color)
	if bold {
		style = style.Bold(true)
	}

	screen.SetContent(x, y, char, nil, style)

	// Add signal name if labels are enabled
	if rd.config.ShowSignalNames && signal.Strength > 50 {
		rd.drawSignalLabel(screen, x, y, signal.Name, theme)
	}
}

// Draw selection border around selected signal
func (rd *Display) drawSelectionBorder(screen tcell.Screen, centerX, centerY int, theme RadarTheme) {
	borderPositions := []struct{ dx, dy int }{
		{-1, -1}, {0, -1}, {1, -1},
		{-1, 0}, {1, 0},
		{-1, 1}, {0, 1}, {1, 1},
	}

	for _, pos := range borderPositions {
		x := centerX + pos.dx
		y := centerY + pos.dy

		if x >= 0 && x < rd.width && y >= 3 && y < rd.height-3 {
			screen.SetContent(x, y, '·', nil,
				tcell.StyleDefault.Foreground(theme.AccentPrimary).Bold(true))
		}
	}
}

// Draw signal trails showing movement history
func (rd *Display) drawSignalTrail(screen tcell.Screen, signal Signal, scaleFactor float64, theme RadarTheme) {
	if len(signal.History) < 2 {
		return
	}

	maxHistory := min(len(signal.History), rd.config.MaxTrailLength)

	for i := len(signal.History) - maxHistory; i < len(signal.History)-1; i++ {
		if i < 0 {
			continue
		}

		h := signal.History[i]
		intensity := float64(i) / float64(len(signal.History)-1)

		distance := h.Distance * scaleFactor
		if distance > float64(min(rd.width, rd.height))/2.1 {
			continue
		}

		trailX := rd.centerX + int(math.Round(math.Cos(h.Angle)*distance))
		trailY := rd.centerY + int(math.Round(math.Sin(h.Angle)*distance*0.5))

		if trailX >= 0 && trailX < rd.width && trailY >= 3 && trailY < rd.height-3 {
			var char rune = '·'
			var color tcell.Color = theme.SweepTrail

			if intensity > 0.7 {
				char = '○'
				color = theme.SweepFade
			}

			screen.SetContent(trailX, trailY, char, nil,
				tcell.StyleDefault.Foreground(color))
		}
	}
}

// Draw enhanced ripple effects for strong signals
func (rd *Display) drawEnhancedSignalRipples(screen tcell.Screen, centerX, centerY int, signal Signal, theme RadarTheme) {
	rippleRadius := float64((signal.Phase % 4) + 2)

	// Draw ripple circle
	for angle := 0.0; angle < 2*math.Pi; angle += math.Pi / 8 {
		x := centerX + int(rippleRadius*math.Cos(angle))
		y := centerY + int(rippleRadius*math.Sin(angle)*0.5)

		if x >= 0 && x < rd.width && y >= 3 && y < rd.height-3 {
			screen.SetContent(x, y, SignalRipple, nil,
				tcell.StyleDefault.Foreground(theme.SweepFade))
		}
	}
}

// Draw signal labels with background
func (rd *Display) drawSignalLabel(screen tcell.Screen, x, y int, label string, theme RadarTheme) {
	if len(label) == 0 {
		return
	}

	// Truncate long labels
	maxLen := 12
	if len(label) > maxLen {
		label = label[:maxLen-3] + "..."
	}

	labelX := x + 2
	labelY := y - 1

	if labelX+len(label) >= rd.width || labelY < 3 {
		labelX = x - len(label) - 2
	}
	if labelX < 0 {
		return
	}

	// Draw label background
	for i := 0; i < len(label); i++ {
		screen.SetContent(labelX+i, labelY, ' ', nil,
			tcell.StyleDefault.Background(tcell.ColorDarkGray))
	}

	// Draw label text
	for i, r := range label {
		screen.SetContent(labelX+i, labelY, r, nil,
			tcell.StyleDefault.Foreground(theme.TextPrimary).
				Background(tcell.ColorDarkGray))
	}
}
