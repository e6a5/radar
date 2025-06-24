# Wave-in-Terminal v2.0

A stunning real-time terminal-based radar system that can visualize both simulated and **real device/network data** as animated waves. Features professional radar sweep trails, range rings, real WiFi network scanning, network device discovery, and beautiful signal animations. Built with Go and the tcell library for smooth terminal UI rendering.

![Radar Animation](https://img.shields.io/badge/Animation-Radar%20Waves-green)
![Go Version](https://img.shields.io/badge/Go-1.23.2-blue)
![License](https://img.shields.io/badge/License-MIT-yellow)

## Features

### 🎯 **Enhanced Radar System with Real Data**
- **Real Device Detection**: Scans actual WiFi networks, network devices, and system connections
- **Dual Mode Operation**: Toggle between real data and simulation with `D` key
- **Live Network Scanning**: Detects WiFi networks with actual signal strength (RSSI)
- **Device Discovery**: Finds active devices on local network via ping scanning
- **System Activity**: Monitors network connections and processes
- Rotating radar beam with stunning fading trail effect
- Professional range rings with distance labels
- Real-time signal detection within beam angle
- Ultra-smooth animation with 80ms refresh rate
- Background grid pattern for depth perception

### 📡 **Signal Types with Unique Icons & Movement Patterns**
- **WiFi** (≋) - Wireless network signals with wave pattern (mostly stationary)
- **Bluetooth** (β) - Short-range device connections (moderate movement)
- **Cellular** (▲) - Mobile network towers with triangle icon (high movement)
- **Radio** (◈) - AM/FM radio broadcasts with diamond pattern (stationary)
- **IoT** (◇) - Internet of Things devices with open diamond (minimal movement)
- **Satellite** (★) - Satellite communications with star icon (orbital movement)

### 🎨 **Enhanced Visual Features**
- **Professional UI Layout:**
  - Bordered top panel with title and status
  - Interactive side panel with signal type legend and filtering (on wide screens)
  - Signal count display per type with filter status indicators
  - Bottom control panel with key bindings including filter controls
  - Clean separation of information areas
- **Advanced Signal Visualization:**
  - Type-specific colored icons for each signal
  - Pulsing animation effects for all signals
  - Ripple effects around strong signals (>70%)
  - Enhanced color mixing of type and strength
  - **Signal History Trails:** Visual movement tracking with fading points
  - **Realistic Movement:** Different signal types move with characteristic patterns
  - **Trail Persistence:** 30-second visual history with age-based fading
  - **Interactive Selection:** Yellow selection boxes highlight chosen signals
  - **Detailed Information Panel:** Comprehensive signal analysis overlay
- **Radar Display Enhancements:**
  - Crosshair center point with professional styling
  - Fading sweep trail with multiple intensity levels
  - Concentric range rings with distance labels
  - Subtle background grid for depth perception
- **Adaptive Design:**
  - Automatically scales to terminal dimensions
  - Responsive layout that adapts to screen size
  - Handles terminal resize events gracefully

### 🛤️ **Signal History & Movement Tracking**

The radar now includes sophisticated signal tracking capabilities:

- **Real-time Movement**: Signals move based on their type characteristics
  - **Cellular/Mobile**: High movement probability (25% per update)
  - **Bluetooth**: Moderate movement (15% - mobile devices)
  - **Satellite**: Predictable orbital patterns (20% with consistent direction)
  - **WiFi/Radio/IoT**: Mostly stationary with minor fluctuations (5%)

- **Visual Trail System**: Press `T` to toggle trail visualization
  - **Detected positions**: Bright colored dots using signal type color
  - **Estimated positions**: Gray dots for interpolated movement
  - **Age-based fading**: Trails fade over 30 seconds with 3 intensity levels
  - **Smart filtering**: Trails respect signal type filters

- **Historical Data**: Each signal maintains 20 position points (~40 seconds of history)

### 🌐 **Real Data Collection**

Transform the radar from simulation to actual network/device monitoring:

- **WiFi Network Scanning** (macOS/Linux):
  - Uses system commands (`airport -s` on macOS, `iwlist scan` on Linux)
  - Detects real WiFi networks with actual signal strength (RSSI)
  - Converts signal strength to distance approximation
  - Shows network names and MAC addresses

- **Network Device Discovery**:
  - Scans all active network interfaces
  - Ping sweeps local network ranges to find active devices
  - Maps discovered devices to radar signals
  - Shows gateway connectivity status

- **System Network Activity**:
  - Monitors active network connections via `netstat`
  - Tracks HTTP/HTTPS, SSH, and other connection types
  - Converts connection count to signal strength
  - Real-time network activity visualization

- **Bluetooth Device Scanning** (Linux):
  - Uses `hcitool scan` for Bluetooth discovery
  - Shows nearby Bluetooth devices
  - Limited range appropriate for Bluetooth

- **Smart Fallback System**:
  - Automatically falls back to simulation if real data unavailable
  - Graceful handling of permission issues
  - Cross-platform compatibility

- **Performance Optimized**:
  - Configurable scan intervals (default: 1 second)
  - Cached results to prevent excessive system calls
  - Non-blocking scans that don't interrupt radar animation

**Usage**: Press `D` to toggle between "REAL" and "SIM" modes. Real mode automatically scans your system for actual devices and networks!

### 🔍 **Signal Selection & Detailed Analysis**

Advanced signal inspection capabilities for professional radar operation:

- **Interactive Selection**: Press `N`/`P` to cycle through visible signals
- **Visual Highlighting**: Selected signals show bright yellow selection boxes
- **Detailed Information Panel**: Press `I` to toggle comprehensive signal data
  - **Signal Identity**: Type, icon, strength percentage and label
  - **Position Data**: Precise distance, bearing in degrees
  - **Temporal Info**: Signal age, last detection time, persistence level  
  - **Movement Analysis**: Historical position count and movement calculations
  - **Real-time Updates**: All data updates live as signals move and change

- **Professional Workflow**: 
  - Select signal of interest with `N`/`P`
  - Open detailed panel with `I` for in-depth analysis
  - Track movement patterns and signal characteristics
  - Clear selection with `C` to return to overview mode

### 🎮 **Interactive Controls**

| Key | Action |
|-----|--------|
| `ESC` / `Q` | Quit application |
| `SPACE` | Pause/Resume animation |
| `+` / `=` | Increase radar speed |
| `-` / `_` | Decrease radar speed |
| `R` | Reset simulation |
| `1` | Toggle WiFi signals |
| `2` | Toggle Bluetooth signals |
| `3` | Toggle Cellular signals |
| `4` | Toggle Radio signals |
| `5` | Toggle IoT signals |
| `6` | Toggle Satellite signals |
| `0` | Toggle all signals |
| `T` | Toggle signal trails |
| `I` | Toggle information panel |
| `N` | Select next signal |
| `P` | Select previous signal |
| `C` | Clear signal selection |
| `D` | Toggle real data mode |

### 📊 **Status Information**
- Real-time signal count display (shows only visible/filtered signals)
- **Data mode indicator**: Shows "REAL" when using actual device data, "SIM" for simulation
- Signal type filtering with live counts in side panel
- Filter status indicators (dimmed when filtered out)
- Signal trail status indicator (TRAILS shows when enabled)
- Signal selection indicator (SEL:N shows selected signal number)
- Pause status indicator
- Control hints at bottom of screen

## Installation & Usage

### Prerequisites
- Go 1.23.2 or later
- Terminal with color support

### Quick Start

```bash
# Clone the repository
git clone https://github.com/your-username/wave-in-terminal.git
cd wave-in-terminal

# Run the application
go run main.go
```

### Building

```bash
# Build binary
go build -o radar main.go

# Run binary
./radar
```

### Testing Real Data Collection

Before running the main application, you can test what real data sources are available on your system:

```bash
# Run the real data collection demo
go run demo_real_data.go
```

This will show:
- Available WiFi networks and their signal strengths
- Active network interfaces and gateway connectivity
- Current network connections and activity
- Bluetooth device scanning capabilities
- System compatibility and permissions status

## Configuration

The application includes a built-in configuration system with the following defaults:

```go
Config{
    RefreshRate:    100 * time.Millisecond,  // Animation speed
    RadarSpeed:     math.Pi / 30,            // Radar rotation speed
    MaxSignals:     8,                       // Maximum concurrent signals
    SignalLifetime: 30 * time.Second,        // How long signals persist
    BeamWidth:      math.Pi / 60,            // Radar detection angle
    MaxPhase:       4,                       // Wave animation phases
}
```

## Architecture

The project is now organized into modular files for better maintainability:

### File Structure

**Main Package:**
- **`main.go`** (38 lines): Application entry point and main loop

**Radar Package** (`./radar/`):
1. **`display.go`** (413 lines): Display struct and core radar logic with input handling
2. **`renderer.go`** (713 lines): Screen rendering and status display  
3. **`signal.go`** (200+ lines): Signal struct and signal management with history
4. **`config.go`** (35 lines): Configuration structure with real data settings
5. **`utils.go`** (28 lines): Utility functions (min/max helpers)
6. **`realdata.go`** (350+ lines): Real device and network data collection

**Demo Files:**
- **`demo_real_data.go`** (180 lines): Test script for real data collection capabilities

### Core Components

1. **radar.Display**: Main display controller
   - Manages screen dimensions and radar state
   - Handles input events and coordinates components
   - Signal lifecycle management

2. **radar.Signal**: Individual signal representation
   - Type, strength, position, and timing data
   - Animation phase tracking
   - Lifetime management

3. **radar.Config**: Centralized configuration
   - Performance tuning parameters
   - Visual and behavioral settings

4. **Renderer**: Display logic separation
   - Screen drawing operations
   - Status information display
   - Color and visual styling

### Key Improvements Over Original

- ✅ **Interactive Controls**: Full keyboard input handling
- ✅ **Dynamic Sizing**: Adapts to any terminal size
- ✅ **Signal Lifecycle**: Signals appear/disappear naturally
- ✅ **Configuration System**: Centralized, adjustable parameters
- ✅ **Package Structure**: Organized into main + radar package (370 lines total)
- ✅ **Enhanced Visuals**: Status information and improved scaling
- ✅ **Pause Functionality**: Start/stop animation control
- ✅ **Speed Control**: Adjustable radar rotation speed
- ✅ **Clean Separation**: Rendering, logic, and data clearly separated

## Technical Details

### Dependencies
- `github.com/gdamore/tcell/v2` - Terminal cell manipulation
- Standard Go libraries (math, time, fmt, log)

### Performance
- Optimized rendering with minimal screen updates
- Efficient trigonometric calculations
- Memory-conscious signal management

### Compatibility
- Cross-platform (Windows, macOS, Linux)
- Works in most terminal environments
- Supports various terminal sizes

## Future Enhancements

Potential areas for expansion:
- [x] Signal filtering by type ✅
- [x] Signal history tracking and trails ✅
- [x] Signal selection and detailed info panel ✅
- [x] **Real device/network data collection** ✅
- [ ] Configuration file support & user preferences
- [ ] Export signal data (CSV/JSON)
- [ ] Enhanced real data sources (more Bluetooth support, process monitoring)
- [ ] Multiple radar displays and split-screen modes
- [ ] Signal strength logging and historical charts
- [ ] Custom signal types and threat classifications
- [ ] Audio alerts for signal detection
- [ ] GPS integration for mobile radar systems

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built with [tcell](https://github.com/gdamore/tcell) - Excellent terminal UI library
- Inspired by classic radar displays and signal analysis tools 