package scanner

import (
	"context"
	"time"
)

// Signal represents a detected signal on the radar
type Signal struct {
	Type        string
	Icon        string
	Name        string
	Color       interface{} // tcell.Color
	Strength    int         // 0-100%
	Distance    float64     // radar units
	Angle       float64     // radians
	Phase       int
	Lifetime    time.Time
	LastSeen    time.Time
	Persistence float64
	History     []PositionHistory
	MaxHistory  int
}

// PositionHistory tracks signal movement over time
type PositionHistory struct {
	Distance  float64
	Angle     float64
	Strength  int
	Timestamp time.Time
	IsReal    bool
}

// Scanner defines the interface for all signal scanners
type Scanner interface {
	// Scan returns detected signals
	Scan(ctx context.Context) ([]Signal, error)

	// Name returns a human-readable name for this scanner
	Name() string

	// IsAvailable checks if this scanner can run on current system
	IsAvailable() bool
}

// Config holds scanner configuration
type Config struct {
	ScanInterval  time.Duration
	MaxSignals    int
	MaxScanRange  float64
	UseRealData   bool
	EnableConsent bool
}

// AddToHistory adds a position entry to signal history
func (s *Signal) AddToHistory(distance, angle float64, strength int, isReal bool, timestamp time.Time) {
	entry := PositionHistory{
		Distance:  distance,
		Angle:     angle,
		Strength:  strength,
		Timestamp: timestamp,
		IsReal:    isReal,
	}

	s.History = append(s.History, entry)

	// Keep only recent history
	if len(s.History) > s.MaxHistory {
		s.History = s.History[1:]
	}
}
