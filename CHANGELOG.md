# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- OSS project structure
- MIT License
- Contributing guidelines
- Comprehensive .gitignore

### Changed
- Repository URL updated to github.com/e6a5/radar
- Removed committed binaries from git tracking
- Project name simplified from "Wave-in-Terminal" to "Radar"

## [2.0.0] - 2024-12-23

### Added
- Real-time terminal-based radar system
- Dual mode operation (simulation vs real data collection)
- Six distinct signal types with unique icons and movement patterns
- Interactive controls for signal filtering and analysis
- Signal history tracking with visual trails
- Professional UI with borders, panels, and status information
- WiFi network scanning and device discovery
- Cross-platform compatibility (macOS, Linux, Windows)
- Smooth 80ms refresh rate animation
- Signal selection and detailed information panel
- Configurable radar parameters

### Technical Features
- Built with Go 1.23.2 and tcell terminal UI library
- Modular package structure with clear separation of concerns
- Non-blocking real data collection with timeout protection
- Memory-efficient signal lifecycle management
- Responsive design adapting to terminal size
- Vendored dependencies for reproducible builds

### Visual Effects
- Radar sweep with fading trail effect
- Range rings with distance labels
- Background grid pattern for depth perception
- Color-coded signal strength visualization
- Pulsing signal animations and ripple effects
- Signal movement trails with history tracking 