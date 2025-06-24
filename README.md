# Radar v2.0

A real-time terminal-based radar system that visualizes network signals and device activity as animated radar sweeps. Monitor WiFi networks, Bluetooth devices, and system connections with professional radar display effects.

![Go Version](https://img.shields.io/badge/Go-1.23.2-blue)
![License](https://img.shields.io/badge/License-MIT-yellow)

## Features

- **Real & Simulated Data**: Toggle between actual WiFi/network scanning and simulation
- **Signal Types**: WiFi (≋), Bluetooth (β), Cellular (▲), Radio (◈), IoT (◇), Satellite (★)
- **Interactive Controls**: Filter signals, adjust speed, select for detailed analysis
- **Visual Effects**: Rotating radar beam, range rings, signal trails, and smooth animations
- **Cross-Platform**: Works on macOS, Linux, and Windows terminals

## Quick Start

```bash
# Clone and run
git clone https://github.com/e6a5/radar.git
cd radar
go run main.go
```

## Controls

| Key | Action |
|-----|--------|
| `ESC`/`Q` | Quit |
| `SPACE` | Pause/Resume |
| `+`/`-` | Adjust radar speed |
| `S` | Toggle simulation mode |
| `1-6` | Toggle signal types |
| `T` | Toggle signal trails |
| `N`/`P` | Select signals |
| `I` | Show signal details |

## Requirements

- Go 1.23.2 or later
- Terminal with color support

## Features

**Real Data Collection** (Default Mode):
- **WiFi Networks**: Scans actual networks with signal strength
- **Network Activity**: Monitors active connections and system processes
- **Device Discovery**: Finds devices on local network

*Press `S` to toggle simulation mode if real data collection is unavailable.*

**Privacy Note**: On first run, you'll be asked for permission to collect device data. Your consent is saved and you won't be prompted again. To revoke consent, delete the file `~/.radar_consent`.

## Building

```bash
# Build binary
go build -o radar main.go

# Run binary
./radar
```

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Contributing

We welcome contributions! See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines. 