package radar

import (
	"fmt"
	"math"

	"github.com/gdamore/tcell/v2"
)

// Enhanced signal rendering that keeps original quality but adds subtle improvements
func (rd *Display) drawEnhancedSignalsSubtle(screen tcell.Screen) {
	maxRadius := float64(min(rd.width, rd.height)) / 2.1

	for i, s := range rd.signals {
		if !rd.angleWithinRadar(s.Angle) {
			continue
		}

		// Calculate signal position with enhanced smoothness
		phaseOffset := float64(s.Phase) * 0.3
		distance := s.Distance + phaseOffset
		scaleFactor := maxRadius / 10.0

		signalX := rd.centerX + int(math.Round(math.Cos(s.Angle)*distance*scaleFactor))
		signalY := rd.centerY + int(math.Round(math.Sin(s.Angle)*distance*scaleFactor*0.5))

		if signalX >= 0 && signalX < rd.width && signalY >= 3 && signalY < rd.height-3 {
			// Clear area around signal for visibility (CRITICAL for visibility over sweep)
			rd.clearSignalAreaSubtle(screen, signalX, signalY)

			// Enhanced signal icon based on strength (SUBTLE improvements)
			var signalChar rune
			var signalColor tcell.Color
			var bold bool

			// Improved signal visualization
			switch {
			case s.Strength >= 80:
				signalChar = '●' // Strong signal - filled circle
				signalColor = tcell.ColorRed
				bold = true
			case s.Strength >= 60:
				signalChar = '◉' // Good signal - dotted circle
				signalColor = tcell.ColorOrange
				bold = true
			case s.Strength >= 40:
				signalChar = '○' // Medium signal - empty circle
				signalColor = tcell.ColorYellow
				bold = false
			default:
				signalChar = '·' // Weak signal - dot
				signalColor = tcell.ColorGray
				bold = false
			}

			// Special styling for WiFi signals
			if s.Type == "WiFi" && s.Strength > 70 {
				signalColor = tcell.ColorLime
				signalChar = '◎'
			}

			// Draw selection indicator for selected signal
			if i == rd.selectedSignalIndex {
				// Subtle selection border
				positions := []struct{ dx, dy int }{
					{-1, 0}, {1, 0}, {0, -1}, {0, 1},
				}
				for _, pos := range positions {
					x := signalX + pos.dx
					y := signalY + pos.dy
					if x >= 0 && x < rd.width && y >= 3 && y < rd.height-3 {
						screen.SetContent(x, y, '·', nil,
							tcell.StyleDefault.Foreground(tcell.ColorBlue).Bold(true))
					}
				}
				bold = true
			}

			// Draw the signal with enhanced styling
			style := tcell.StyleDefault.Foreground(signalColor)
			if bold {
				style = style.Bold(true)
			}

			screen.SetContent(signalX, signalY, signalChar, nil, style)

			// Add subtle signal name label for strong signals
			if rd.config.ShowSignalNames && s.Strength > 60 && len(s.Name) > 0 {
				rd.drawSubtleSignalLabel(screen, signalX, signalY, s.Name)
			}
		}
	}
}

// Clear area around signal to ensure visibility over sweep
func (rd *Display) clearSignalAreaSubtle(screen tcell.Screen, centerX, centerY int) {
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

// Draw subtle signal labels that don't clutter the display
func (rd *Display) drawSubtleSignalLabel(screen tcell.Screen, x, y int, label string) {
	// Truncate long labels
	maxLen := 8
	if len(label) > maxLen {
		label = label[:maxLen-2] + ".."
	}

	// Position label to avoid center area
	labelX := x + 2
	labelY := y - 1

	// Check bounds and adjust position
	if labelX+len(label) >= rd.width {
		labelX = x - len(label) - 1
	}
	if labelX < 0 || labelY < 3 {
		return
	}

	// Draw label with subtle background
	for i, r := range label {
		if labelX+i < rd.width && labelY < rd.height-3 {
			screen.SetContent(labelX+i, labelY, r, nil,
				tcell.StyleDefault.Foreground(tcell.ColorDarkGray).
					Background(tcell.ColorBlack))
		}
	}
}

// Enhanced status bar that keeps the original style but adds useful info
func (rd *Display) drawEnhancedStatusSubtle(screen tcell.Screen) {
	if rd.height < 4 {
		return
	}

	// Keep original status bar but enhance the information
	statusY := rd.height - 1

	// Clear the status line
	for x := 0; x < rd.width; x++ {
		screen.SetContent(x, statusY, ' ', nil,
			tcell.StyleDefault.Background(tcell.ColorNavy))
	}

	// Build enhanced status components
	var statusParts []string

	// System status
	if rd.paused {
		statusParts = append(statusParts, "■ PAUSED")
	} else {
		statusParts = append(statusParts, "▶ ACTIVE")
	}

	// Data mode
	if rd.config.EnableRealData {
		statusParts = append(statusParts, "REAL-TIME")
	} else {
		statusParts = append(statusParts, "SIMULATION")
	}

	// Signal count
	visibleCount := rd.getVisibleSignalCount()
	totalCount := len(rd.signals)
	statusParts = append(statusParts, fmt.Sprintf("%d/%d SIGNALS", visibleCount, totalCount))

	// Speed
	speed := rd.config.RadarSpeed * 30 / 3.14159
	statusParts = append(statusParts, fmt.Sprintf("%.1fx", speed))

	// Join status parts
	statusText := ""
	for i, part := range statusParts {
		if i > 0 {
			statusText += " │ "
		}
		statusText += part
	}

	// Draw status text
	x := 2
	for _, r := range statusText {
		if x < rd.width-10 {
			screen.SetContent(x, statusY, r, nil,
				tcell.StyleDefault.Background(tcell.ColorNavy).
					Foreground(tcell.ColorWhite).Bold(true))
			x++
		}
	}

	// Add help hint
	helpText := "Press H for help"
	helpX := rd.width - len(helpText) - 2
	if helpX > x {
		for i, r := range helpText {
			screen.SetContent(helpX+i, statusY, r, nil,
				tcell.StyleDefault.Background(tcell.ColorNavy).
					Foreground(tcell.ColorGray))
		}
	}
}

// Enhanced range rings that improve the original without breaking it
func (rd *Display) drawEnhancedRangeRingsSubtle(screen tcell.Screen) {
	maxRadius := float64(min(rd.width, rd.height)) / 2.1

	for ring := 1; ring <= 4; ring++ {
		radius := maxRadius * float64(ring) / 4.0

		// Keep original ring style but with better colors
		ringChar := '○'
		ringColor := tcell.ColorGreen

		if ring%2 == 0 {
			ringChar = '●'
			ringColor = tcell.ColorDarkGreen
		}

		// Draw ring with original method but enhanced smoothness
		rd.drawCircleEnhanced(screen, rd.centerX, rd.centerY, radius, ringChar, ringColor)

		// Enhanced range labels with better positioning
		if ring*10 <= 40 {
			label := fmt.Sprintf("%dm", ring*10)

			// Try multiple label positions
			labelPositions := []struct{ x, y int }{
				{rd.centerX + int(radius) - len(label)/2, rd.centerY - 1},
				{rd.centerX + int(radius) - len(label)/2, rd.centerY + int(radius*0.5) + 1},
				{rd.centerX - int(radius) + 1, rd.centerY},
				{rd.centerX, rd.centerY - int(radius*0.5) - 1},
			}

			for _, pos := range labelPositions {
				if pos.x > 0 && pos.x < rd.width-len(label) &&
					pos.y > 2 && pos.y < rd.height-3 {

					// Draw subtle background
					for i := 0; i < len(label); i++ {
						screen.SetContent(pos.x+i, pos.y, ' ', nil,
							tcell.StyleDefault.Background(tcell.ColorDarkSlateGray))
					}

					// Draw label
					for i, r := range label {
						screen.SetContent(pos.x+i, pos.y, r, nil,
							tcell.StyleDefault.Foreground(tcell.ColorYellow).
								Background(tcell.ColorDarkSlateGray).Bold(true))
					}
					break
				}
			}
		}
	}
}

// Enhanced circle drawing that's smoother than original
func (rd *Display) drawCircleEnhanced(screen tcell.Screen, centerX, centerY int, radius float64, char rune, color tcell.Color) {
	// Use smaller step for smoother circles
	stepSize := 0.04 // Smaller than original 0.05
	for angle := 0.0; angle < 2*math.Pi; angle += stepSize {
		x := centerX + int(radius*math.Cos(angle))
		y := centerY + int(radius*math.Sin(angle)*0.5)
		if x >= 0 && x < rd.width && y >= 3 && y < rd.height-3 {
			screen.SetContent(x, y, char, nil, tcell.StyleDefault.Foreground(color).Bold(true))
		}
	}
}
