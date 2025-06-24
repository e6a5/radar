package radar

import (
	"github.com/gdamore/tcell/v2"
)

// ZoomMode represents different zoom interaction modes
type ZoomMode int

const (
	ZoomNormal ZoomMode = iota
	ZoomIn
	ZoomOut
	PanMode
)

// Enhanced key bindings with zoom and pan controls
func (rd *Display) HandleAdvancedInput(screen tcell.Screen) bool {
	if !screen.HasPendingEvent() {
		return true // No event pending, continue
	}

	event := screen.PollEvent()
	if event == nil {
		return true
	}

	switch ev := event.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyEscape:
			return false
		case tcell.KeyEnter:
			rd.paused = !rd.paused
		case tcell.KeyUp:
			if rd.config.EnablePan {
				rd.config.PanY -= 5
			}
		case tcell.KeyDown:
			if rd.config.EnablePan {
				rd.config.PanY += 5
			}
		case tcell.KeyLeft:
			if rd.config.EnablePan {
				rd.config.PanX -= 5
			}
		case tcell.KeyRight:
			if rd.config.EnablePan {
				rd.config.PanX += 5
			}
		case tcell.KeyHome:
			// Reset pan to center
			if rd.config.EnablePan {
				rd.config.PanX = 0
				rd.config.PanY = 0
			}
		default:
			if ev.Rune() != 0 {
				switch ev.Rune() {
				case 'q', 'Q':
					return false
				case ' ':
					rd.paused = !rd.paused
				case '+', '=':
					// Zoom in
					if rd.config.EnableZoom {
						rd.zoomIn()
					}
				case '-', '_':
					// Zoom out
					if rd.config.EnableZoom {
						rd.zoomOut()
					}
				case '0':
					// Reset zoom to 1.0
					if rd.config.EnableZoom {
						rd.config.ZoomLevel = 1.0
					}
				case 'z', 'Z':
					// Toggle zoom mode
					rd.config.EnableZoom = !rd.config.EnableZoom
				case 'm', 'M':
					// Toggle pan mode
					rd.config.EnablePan = !rd.config.EnablePan
				case 'r', 'R':
					// Reset both zoom and pan
					rd.resetViewport()
				// Keep existing filter controls
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
				case 'a', 'A':
					// Toggle all signal types
					rd.filters.AllVisible = !rd.filters.AllVisible
					rd.filters.WiFiVisible = rd.filters.AllVisible
					rd.filters.BluetoothVisible = rd.filters.AllVisible
					rd.filters.CellularVisible = rd.filters.AllVisible
					rd.filters.RadioVisible = rd.filters.AllVisible
					rd.filters.IoTVisible = rd.filters.AllVisible
					rd.filters.SatelliteVisible = rd.filters.AllVisible
				case 'f', 'F':
					// Toggle filtering system on/off
					rd.config.EnableFiltering = !rd.config.EnableFiltering
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
					// Toggle real data collection
					rd.toggleDataMode()
				case 'l', 'L':
					// Toggle signal name labels
					rd.config.ShowSignalNames = !rd.config.ShowSignalNames
				case 'v', 'V':
					// Toggle performance monitor display
					rd.showPerformanceStats = !rd.showPerformanceStats
				case 'h', 'H':
					// Show help screen (handle in main display loop)
					rd.showHelp = !rd.showHelp
				}
			}
		}
	case *tcell.EventResize:
		rd.width, rd.height = screen.Size()
		rd.updateCenterAfterResize()
	}
	return true
}

// zoomIn increases the zoom level
func (rd *Display) zoomIn() {
	newZoom := rd.config.ZoomLevel * 1.25
	if newZoom <= rd.config.MaxZoom {
		rd.config.ZoomLevel = newZoom
	}
}

// zoomOut decreases the zoom level
func (rd *Display) zoomOut() {
	newZoom := rd.config.ZoomLevel / 1.25
	if newZoom >= rd.config.MinZoom {
		rd.config.ZoomLevel = newZoom
	}
}

// resetViewport resets zoom and pan to defaults
func (rd *Display) resetViewport() {
	rd.config.ZoomLevel = 1.0
	rd.config.PanX = 0.0
	rd.config.PanY = 0.0
}

// updateCenterAfterResize updates center coordinates after screen resize
func (rd *Display) updateCenterAfterResize() {
	rd.centerX = rd.width / 2
	rd.centerY = rd.height / 2
}

// toggleDataMode switches between real and simulated data
func (rd *Display) toggleDataMode() {
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
}

// getEffectiveCenterX returns the center X coordinate accounting for pan
func (rd *Display) getEffectiveCenterX() int {
	return rd.centerX + int(rd.config.PanX)
}

// getEffectiveCenterY returns the center Y coordinate accounting for pan
func (rd *Display) getEffectiveCenterY() int {
	return rd.centerY + int(rd.config.PanY)
}

// transformRadiusForZoom applies zoom transformation to a radius
func (rd *Display) transformRadiusForZoom(radius float64) float64 {
	return radius * rd.config.ZoomLevel
}

// transformCoordinateForZoomPan applies zoom and pan transformations to coordinates
func (rd *Display) transformCoordinateForZoomPan(x, y int) (int, int) {
	// Apply zoom
	relativeX := float64(x - rd.centerX)
	relativeY := float64(y - rd.centerY)

	zoomedX := relativeX * rd.config.ZoomLevel
	zoomedY := relativeY * rd.config.ZoomLevel

	// Apply pan and return to screen coordinates
	finalX := rd.centerX + int(zoomedX+rd.config.PanX)
	finalY := rd.centerY + int(zoomedY+rd.config.PanY)

	return finalX, finalY
}
