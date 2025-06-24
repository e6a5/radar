package radar

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
)

// Enhanced modern status bar with sleek design
func (rd *Display) drawModernStatusBar(screen tcell.Screen) {
	if rd.height < 4 {
		return
	}

	statusY := rd.height - 1

	// Clear status line with modern background
	for x := 0; x < rd.width; x++ {
		screen.SetContent(x, statusY, ' ', nil,
			tcell.StyleDefault.Background(tcell.ColorDarkBlue).Foreground(tcell.ColorWhite))
	}

	// Status components with modern icons
	var statusComponents []StatusComponent

	// System status
	if rd.paused {
		statusComponents = append(statusComponents, StatusComponent{
			Icon: "■", Text: "PAUSED", Color: tcell.ColorRed, Bold: true,
		})
	} else {
		statusComponents = append(statusComponents, StatusComponent{
			Icon: "▶", Text: "ACTIVE", Color: tcell.ColorGreen, Bold: true,
		})
	}

	// Data mode
	if rd.config.EnableRealData {
		statusComponents = append(statusComponents, StatusComponent{
			Icon: "📡", Text: "REAL-TIME", Color: tcell.ColorYellow, Bold: true,
		})
	} else {
		statusComponents = append(statusComponents, StatusComponent{
			Icon: "🎲", Text: "SIMULATION", Color: tcell.ColorGray, Bold: false,
		})
	}

	// Signal statistics
	visibleCount := rd.getVisibleSignalCount()
	totalCount := len(rd.signals)
	statusComponents = append(statusComponents, StatusComponent{
		Icon: "📶", Text: fmt.Sprintf("%d/%d SIGNALS", visibleCount, totalCount),
		Color: tcell.ColorWhite, Bold: false,
	})

	// Performance indicators
	if rd.showPerformanceStats {
		statusComponents = append(statusComponents, StatusComponent{
			Icon: "⚡", Text: "ENHANCED MODE", Color: tcell.ColorBlue, Bold: true,
		})
	}

	// Speed indicator
	speed := rd.config.RadarSpeed * 30 / 3.14159
	statusComponents = append(statusComponents, StatusComponent{
		Icon: "🔄", Text: fmt.Sprintf("%.1fx", speed),
		Color: tcell.ColorWhite, Bold: false,
	})

	// Draw status components with modern styling
	x := 2
	for i, comp := range statusComponents {
		if i > 0 {
			// Draw separator
			screen.SetContent(x, statusY, '│', nil,
				tcell.StyleDefault.Background(tcell.ColorDarkBlue).Foreground(tcell.ColorDarkGray))
			x += 2
		}

		// Draw icon (simplified for compatibility)
		if len(comp.Icon) > 0 && comp.Icon[0] < 128 { // ASCII only
			iconStyle := tcell.StyleDefault.Background(tcell.ColorDarkBlue).Foreground(comp.Color)
			if comp.Bold {
				iconStyle = iconStyle.Bold(true)
			}
			screen.SetContent(x, statusY, rune(comp.Icon[0]), nil, iconStyle)
			x += 2
		}

		// Draw text
		textStyle := tcell.StyleDefault.Background(tcell.ColorDarkBlue).Foreground(comp.Color)
		if comp.Bold {
			textStyle = textStyle.Bold(true)
		}

		for _, r := range comp.Text {
			if x < rd.width-2 {
				screen.SetContent(x, statusY, r, nil, textStyle)
				x++
			}
		}
		x += 1
	}

	// Add help hint on the right
	helpText := "Press H for help"
	helpX := rd.width - len(helpText) - 2
	if helpX > x {
		for i, r := range helpText {
			screen.SetContent(helpX+i, statusY, r, nil,
				tcell.StyleDefault.Background(tcell.ColorDarkBlue).Foreground(tcell.ColorDarkGray))
		}
	}
}

// StatusComponent represents a component in the status bar
type StatusComponent struct {
	Icon  string
	Text  string
	Color tcell.Color
	Bold  bool
}

// Enhanced signal information panel with modern design
func (rd *Display) drawModernSignalPanel(screen tcell.Screen) {
	signal := rd.getSelectedSignal()
	if signal == nil {
		return
	}

	theme := GetModernDarkTheme()

	// Panel dimensions and positioning
	panelWidth := 40
	panelHeight := 18
	panelX := rd.width - panelWidth - 1
	panelY := 3

	if panelX < 0 || panelY < 0 {
		return
	}

	// Draw modern panel background with gradient effect
	for y := panelY; y < panelY+panelHeight && y < rd.height-2; y++ {
		for x := panelX; x < panelX+panelWidth && x < rd.width; x++ {
			// Create subtle gradient effect
			bgColor := tcell.ColorDarkBlue
			if y == panelY || y == panelY+panelHeight-1 {
				bgColor = tcell.ColorBlue // Lighter for borders
			}
			screen.SetContent(x, y, ' ', nil,
				tcell.StyleDefault.Background(bgColor))
		}
	}

	// Draw modern border with corner elements
	rd.drawModernPanelBorder(screen, panelX, panelY, panelWidth, panelHeight, theme)

	// Panel title with modern styling
	title := "▸ SIGNAL ANALYSIS"
	titleX := panelX + (panelWidth-len(title))/2
	for i, r := range title {
		screen.SetContent(titleX+i, panelY+1, r, nil,
			tcell.StyleDefault.Background(tcell.ColorDarkBlue).
				Foreground(tcell.ColorYellow).Bold(true))
	}

	// Signal details with enhanced formatting
	details := []PanelLine{
		{Label: "Name", Value: signal.Name, Color: tcell.ColorWhite},
		{Label: "Type", Value: fmt.Sprintf("%s %s", signal.Icon, signal.Type), Color: tcell.ColorBlue},
		{Label: "Strength", Value: fmt.Sprintf("%d%% (%s)", signal.Strength, rd.getStrengthDescription(signal.Strength)), Color: rd.getStrengthColor(signal.Strength)},
		{Label: "Distance", Value: fmt.Sprintf("%.0fm", signal.Distance), Color: tcell.ColorWhite},
		{Label: "Bearing", Value: fmt.Sprintf("%.1f°", signal.Angle*180/3.14159), Color: tcell.ColorWhite},
		{Label: "Phase", Value: fmt.Sprintf("%d/%d", signal.Phase, rd.config.MaxPhase), Color: tcell.ColorGray},
		{Label: "", Value: "", Color: tcell.ColorBlack}, // Spacer
		{Label: "Status", Value: rd.getSignalStatus(*signal), Color: tcell.ColorGreen},
		{Label: "First Seen", Value: signal.Lifetime.Format("15:04:05"), Color: tcell.ColorGray},
		{Label: "Last Update", Value: signal.LastSeen.Format("15:04:05"), Color: tcell.ColorGray},
		{Label: "Persistence", Value: fmt.Sprintf("%.2f", signal.Persistence), Color: tcell.ColorWhite},
		{Label: "", Value: "", Color: tcell.ColorBlack}, // Spacer
		{Label: "History", Value: fmt.Sprintf("%d waypoints", len(signal.History)), Color: tcell.ColorBlue},
	}

	// Draw signal details with modern styling
	for i, line := range details {
		y := panelY + 3 + i
		if y >= panelY+panelHeight-1 {
			break
		}

		if line.Label == "" {
			continue // Skip spacers
		}

		// Draw label
		labelX := panelX + 2
		for j, r := range line.Label + ":" {
			if labelX+j < panelX+panelWidth-2 {
				screen.SetContent(labelX+j, y, r, nil,
					tcell.StyleDefault.Background(tcell.ColorDarkBlue).
						Foreground(tcell.ColorDarkGray))
			}
		}

		// Draw value
		valueX := panelX + 14
		maxValueWidth := panelWidth - 16
		value := line.Value
		if len(value) > maxValueWidth {
			value = value[:maxValueWidth-3] + "..."
		}

		for j, r := range value {
			if valueX+j < panelX+panelWidth-2 {
				screen.SetContent(valueX+j, y, r, nil,
					tcell.StyleDefault.Background(tcell.ColorDarkBlue).
						Foreground(line.Color).Bold(true))
			}
		}
	}
}

// PanelLine represents a line in the information panel
type PanelLine struct {
	Label string
	Value string
	Color tcell.Color
}

// Draw modern panel border with enhanced styling
func (rd *Display) drawModernPanelBorder(screen tcell.Screen, x, y, width, height int, theme RadarTheme) {
	borderColor := theme.AccentPrimary

	// Top and bottom borders
	for i := 0; i < width; i++ {
		if x+i >= 0 && x+i < rd.width {
			if y >= 0 && y < rd.height {
				char := '─'
				if i == 0 {
					char = '┌'
				} else if i == width-1 {
					char = '┐'
				}
				screen.SetContent(x+i, y, char, nil,
					tcell.StyleDefault.Background(tcell.ColorDarkBlue).
						Foreground(borderColor).Bold(true))
			}

			if y+height-1 >= 0 && y+height-1 < rd.height {
				char := '─'
				if i == 0 {
					char = '└'
				} else if i == width-1 {
					char = '┘'
				}
				screen.SetContent(x+i, y+height-1, char, nil,
					tcell.StyleDefault.Background(tcell.ColorDarkBlue).
						Foreground(borderColor).Bold(true))
			}
		}
	}

	// Left and right borders
	for i := 1; i < height-1; i++ {
		if y+i >= 0 && y+i < rd.height {
			if x >= 0 && x < rd.width {
				screen.SetContent(x, y+i, '│', nil,
					tcell.StyleDefault.Background(tcell.ColorDarkBlue).
						Foreground(borderColor).Bold(true))
			}
			if x+width-1 >= 0 && x+width-1 < rd.width {
				screen.SetContent(x+width-1, y+i, '│', nil,
					tcell.StyleDefault.Background(tcell.ColorDarkBlue).
						Foreground(borderColor).Bold(true))
			}
		}
	}
}

// Get signal status description
func (rd *Display) getSignalStatus(signal Signal) string {
	if !signal.IsVisible() {
		return "Hidden"
	}

	timeSinceLastSeen := time.Since(signal.LastSeen)
	if timeSinceLastSeen < 5*time.Second {
		return "Active"
	} else if timeSinceLastSeen < 30*time.Second {
		return "Recent"
	} else {
		return "Stale"
	}
}

// Get signal strength description
func (rd *Display) getStrengthDescription(strength int) string {
	switch {
	case strength >= 90:
		return "Excellent"
	case strength >= 75:
		return "Very Good"
	case strength >= 60:
		return "Good"
	case strength >= 45:
		return "Fair"
	case strength >= 30:
		return "Poor"
	case strength >= 15:
		return "Very Poor"
	default:
		return "Critical"
	}
}

// Get color based on signal strength
func (rd *Display) getStrengthColor(strength int) tcell.Color {
	switch {
	case strength >= 75:
		return tcell.ColorGreen
	case strength >= 60:
		return tcell.ColorYellow
	case strength >= 30:
		return tcell.ColorOrange
	default:
		return tcell.ColorRed
	}
}
