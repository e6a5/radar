package radar

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

// showHelpScreen displays a comprehensive help overlay
func (rd *Display) showHelpScreen(screen tcell.Screen) {
	helpLines := []string{
		"RADAR v2.0 - Enhanced Network Monitor",
		"====================================",
		"",
		"NAVIGATION & VIEW CONTROLS:",
		"  Arrow Keys - Pan radar view (when pan mode enabled)",
		"  + / =      - Zoom in",
		"  - / _      - Zoom out",
		"  0          - Reset zoom to 1.0x",
		"  R          - Reset zoom and pan",
		"  Z          - Toggle zoom mode",
		"  M          - Toggle pan mode",
		"",
		"SIGNAL CONTROLS:",
		"  SPACE      - Pause/Resume radar",
		"  1-6        - Toggle signal types (WiFi, Bluetooth, etc.)",
		"  A          - Toggle all signal types",
		"  F          - Toggle filtering system",
		"  T          - Toggle signal trails",
		"  L          - Toggle signal labels",
		"  S          - Switch real/simulated data",
		"",
		"INFORMATION & SELECTION:",
		"  N          - Select next signal",
		"  P          - Select previous signal",
		"  C          - Clear signal selection",
		"  I          - Toggle information panel",
		"  V          - Toggle performance stats",
		"",
		"ADVANCED:",
		"  H          - Show/hide this help",
		"  Q/ESC      - Quit application",
		"",
		"PERFORMANCE FEATURES:",
		"  * Adaptive refresh rate (auto-adjusts based on performance)",
		"  * Spatial caching (optimized circle/trig calculations)",
		"  * Smart signal filtering (reduces rendering load)",
		"  * Real-time performance monitoring",
		"",
		"SIGNAL TYPES:",
		"  ≋ WiFi Networks    # Cellular Towers",
		"  b Bluetooth        @ Radio Stations",
		"  ^ IoT Devices      * Satellites",
		"",
		"Press any key to close help...",
	}

	// Calculate help window dimensions
	maxWidth := 60
	height := len(helpLines) + 4
	startY := (rd.height - height) / 2
	startX := (rd.width - maxWidth) / 2

	// Draw help window background
	for y := startY; y < startY+height; y++ {
		for x := startX; x < startX+maxWidth; x++ {
			if x >= 0 && x < rd.width && y >= 0 && y < rd.height {
				screen.SetContent(x, y, ' ', nil,
					tcell.StyleDefault.Background(tcell.ColorDarkBlue).Foreground(tcell.ColorWhite))
			}
		}
	}

	// Draw help window border
	borderStyle := tcell.StyleDefault.Background(tcell.ColorDarkBlue).Foreground(tcell.ColorGreen).Bold(true)
	for x := startX; x < startX+maxWidth; x++ {
		if x >= 0 && x < rd.width {
			if startY >= 0 && startY < rd.height {
				screen.SetContent(x, startY, '─', nil, borderStyle)
			}
			if startY+height-1 >= 0 && startY+height-1 < rd.height {
				screen.SetContent(x, startY+height-1, '─', nil, borderStyle)
			}
		}
	}
	for y := startY; y < startY+height; y++ {
		if y >= 0 && y < rd.height {
			if startX >= 0 && startX < rd.width {
				screen.SetContent(startX, y, '│', nil, borderStyle)
			}
			if startX+maxWidth-1 >= 0 && startX+maxWidth-1 < rd.width {
				screen.SetContent(startX+maxWidth-1, y, '│', nil, borderStyle)
			}
		}
	}

	// Draw corners
	if startX >= 0 && startX < rd.width && startY >= 0 && startY < rd.height {
		screen.SetContent(startX, startY, '┌', nil, borderStyle)
	}
	if startX+maxWidth-1 >= 0 && startX+maxWidth-1 < rd.width && startY >= 0 && startY < rd.height {
		screen.SetContent(startX+maxWidth-1, startY, '┐', nil, borderStyle)
	}
	if startX >= 0 && startX < rd.width && startY+height-1 >= 0 && startY+height-1 < rd.height {
		screen.SetContent(startX, startY+height-1, '└', nil, borderStyle)
	}
	if startX+maxWidth-1 >= 0 && startX+maxWidth-1 < rd.width && startY+height-1 >= 0 && startY+height-1 < rd.height {
		screen.SetContent(startX+maxWidth-1, startY+height-1, '┘', nil, borderStyle)
	}

	// Draw help text
	for i, line := range helpLines {
		y := startY + 2 + i
		if y >= 0 && y < rd.height {
			// Choose text style based on content
			textStyle := tcell.StyleDefault.Background(tcell.ColorDarkBlue).Foreground(tcell.ColorWhite)

			if len(line) > 0 {
				if line[0] >= 'A' && line[0] <= 'Z' {
					textStyle = textStyle.Foreground(tcell.ColorYellow).Bold(true)
				} else if line[0] == ' ' {
					if len(line) > 2 && (line[2] >= 'A' && line[2] <= 'Z') {
						textStyle = textStyle.Foreground(tcell.ColorGreen)
					}
				}
			}

			x := startX + 2
			for j, r := range line {
				if x+j >= 0 && x+j < rd.width {
					screen.SetContent(x+j, y, r, nil, textStyle)
				}
			}
		}
	}
}

// drawEnhancedStatus draws an enhanced status bar with performance info
func (rd *Display) drawEnhancedStatus(screen tcell.Screen) {
	if rd.height < 4 {
		return // Not enough space
	}

	statusY := rd.height - 1

	// Build status components
	var statusParts []string

	// Basic status
	if rd.paused {
		statusParts = append(statusParts, "⏸ PAUSED")
	} else {
		statusParts = append(statusParts, "▶ RUNNING")
	}

	// Data mode
	if rd.config.EnableRealData {
		statusParts = append(statusParts, "📡 REAL DATA")
	} else {
		statusParts = append(statusParts, "🎲 SIMULATION")
	}

	// Zoom and pan info
	if rd.config.EnableZoom && rd.config.ZoomLevel != 1.0 {
		statusParts = append(statusParts, fmt.Sprintf("🔍 %.1fx", rd.config.ZoomLevel))
	}
	if rd.config.EnablePan && (rd.config.PanX != 0 || rd.config.PanY != 0) {
		statusParts = append(statusParts, fmt.Sprintf("📐 (%.0f,%.0f)", rd.config.PanX, rd.config.PanY))
	}

	// Performance info
	if rd.showPerformanceStats {
		stats := rd.performanceMonitor.GetStats()
		statusParts = append(statusParts,
			fmt.Sprintf("📊 %.1f FPS | %dms", stats.FrameRate, stats.RenderLatency.Milliseconds()))

		// Cache stats
		circle, sin, cos := rd.spatialCache.GetCacheStats()
		statusParts = append(statusParts,
			fmt.Sprintf("💾 C:%d S:%d T:%d", circle, sin, cos))
	}

	// Signal count
	visibleCount := rd.getVisibleSignalCount()
	totalCount := len(rd.signals)
	statusParts = append(statusParts, fmt.Sprintf("📶 %d/%d signals", visibleCount, totalCount))

	// Build full status line
	statusLine := ""
	for i, part := range statusParts {
		if i > 0 {
			statusLine += " | "
		}
		statusLine += part
	}

	// Add help hint
	statusLine += " | Press H for help"

	// Clear status line
	for x := 0; x < rd.width; x++ {
		screen.SetContent(x, statusY, ' ', nil,
			tcell.StyleDefault.Background(tcell.ColorDarkGray).Foreground(tcell.ColorWhite))
	}

	// Draw status text
	x := 1
	for _, r := range statusLine {
		if x < rd.width-1 {
			screen.SetContent(x, statusY, r, nil,
				tcell.StyleDefault.Background(tcell.ColorDarkGray).Foreground(tcell.ColorWhite).Bold(true))
			x++
		}
	}
}

// drawEnhancedSignalInfo draws enhanced signal information panel
func (rd *Display) drawEnhancedSignalInfo(screen tcell.Screen) {
	signal := rd.getSelectedSignal()
	if signal == nil {
		return
	}

	// Panel dimensions
	panelWidth := 35
	panelHeight := 15
	panelX := rd.width - panelWidth - 2
	panelY := 5

	if panelX < 0 || panelY < 0 {
		return // Not enough space
	}

	// Draw panel background
	bgStyle := tcell.StyleDefault.Background(tcell.ColorDarkBlue).Foreground(tcell.ColorWhite)
	for y := panelY; y < panelY+panelHeight && y < rd.height-2; y++ {
		for x := panelX; x < panelX+panelWidth && x < rd.width; x++ {
			screen.SetContent(x, y, ' ', nil, bgStyle)
		}
	}

	// Panel styling (no border for now)

	// Information lines
	infoLines := []string{
		"📊 SIGNAL DETAILS",
		"═════════════════",
		fmt.Sprintf("Name: %s", signal.Name),
		fmt.Sprintf("Type: %s %s", signal.Icon, signal.Type),
		fmt.Sprintf("Strength: %d%% (%s)", signal.Strength, rd.getStrengthText(signal.Strength)),
		fmt.Sprintf("Distance: %.0fm", signal.Distance),
		fmt.Sprintf("Angle: %.1f°", signal.Angle*180/3.14159),
		fmt.Sprintf("Phase: %d/%d", signal.Phase, rd.config.MaxPhase),
		"",
		fmt.Sprintf("First seen: %s", signal.Lifetime.Format("15:04:05")),
		fmt.Sprintf("Last seen: %s", signal.LastSeen.Format("15:04:05")),
		fmt.Sprintf("Persistence: %.1f", signal.Persistence),
		"",
		"History: " + fmt.Sprintf("%d points", len(signal.History)),
	}

	// Draw information
	for i, line := range infoLines {
		y := panelY + 1 + i
		if y >= rd.height-2 {
			break
		}

		textStyle := bgStyle
		if i == 0 {
			textStyle = textStyle.Foreground(tcell.ColorYellow).Bold(true)
		} else if i == 1 {
			textStyle = textStyle.Foreground(tcell.ColorGreen)
		}

		x := panelX + 1
		for _, r := range line {
			if x < panelX+panelWidth-1 {
				screen.SetContent(x, y, r, nil, textStyle)
				x++
			}
		}
	}
}

// getStrengthText returns descriptive text for signal strength
func (rd *Display) getStrengthText(strength int) string {
	switch {
	case strength >= 80:
		return "Excellent"
	case strength >= 60:
		return "Good"
	case strength >= 40:
		return "Fair"
	case strength >= 20:
		return "Poor"
	default:
		return "Very Poor"
	}
}
