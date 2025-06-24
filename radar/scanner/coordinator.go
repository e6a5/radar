package scanner

import (
	"context"
	"sync"
	"time"
)

// Coordinator manages multiple scanners and aggregates their results
type Coordinator struct {
	scanners      []Scanner
	config        *Config
	lastScan      time.Time
	cachedSignals []Signal
	mutex         sync.RWMutex
	isScanning    bool
}

// NewCoordinator creates a new scanner coordinator
func NewCoordinator(config *Config) *Coordinator {
	return &Coordinator{
		scanners:      make([]Scanner, 0),
		config:        config,
		cachedSignals: make([]Signal, 0),
	}
}

// AddScanner adds a scanner to the coordinator
func (c *Coordinator) AddScanner(scanner Scanner) {
	if scanner.IsAvailable() {
		c.mutex.Lock()
		c.scanners = append(c.scanners, scanner)
		c.mutex.Unlock()
	}
}

// GetScanners returns the list of available scanners
func (c *Coordinator) GetScanners() []string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	names := make([]string, len(c.scanners))
	for i, scanner := range c.scanners {
		names[i] = scanner.Name()
	}
	return names
}

// Scan runs all scanners and aggregates results
func (c *Coordinator) Scan(ctx context.Context) ([]Signal, error) {
	c.mutex.RLock()
	now := time.Now()

	// Rate limiting
	if now.Sub(c.lastScan) < c.config.ScanInterval {
		signals := make([]Signal, len(c.cachedSignals))
		copy(signals, c.cachedSignals)
		c.mutex.RUnlock()
		return signals, nil
	}
	c.mutex.RUnlock()

	// Non-blocking scan - start background scan if not already running
	c.mutex.Lock()
	if !c.isScanning {
		c.isScanning = true
		go c.performBackgroundScan(ctx)
	}

	// Return cached signals immediately
	signals := make([]Signal, len(c.cachedSignals))
	copy(signals, c.cachedSignals)
	c.mutex.Unlock()

	return signals, nil
}

// performBackgroundScan runs all scanners in parallel
func (c *Coordinator) performBackgroundScan(ctx context.Context) {
	defer func() {
		c.mutex.Lock()
		c.isScanning = false
		c.lastScan = time.Now()
		c.mutex.Unlock()
	}()

	c.mutex.RLock()
	scanners := make([]Scanner, len(c.scanners))
	copy(scanners, c.scanners)
	c.mutex.RUnlock()

	// Create timeout context
	scanCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Run all scanners in parallel
	type scanResult struct {
		signals []Signal
		err     error
		scanner string
	}

	resultChan := make(chan scanResult, len(scanners))

	for _, scanner := range scanners {
		go func(s Scanner) {
			signals, err := s.Scan(scanCtx)
			resultChan <- scanResult{
				signals: signals,
				err:     err,
				scanner: s.Name(),
			}
		}(scanner)
	}

	// Collect results
	allSignals := make([]Signal, 0)
	for i := 0; i < len(scanners); i++ {
		select {
		case result := <-resultChan:
			if result.err == nil {
				allSignals = append(allSignals, result.signals...)
			}
		case <-scanCtx.Done():
			// Timeout - continue with what we have
			break
		}
	}

	// Limit total signals
	if len(allSignals) > c.config.MaxSignals {
		allSignals = allSignals[:c.config.MaxSignals]
	}

	// Update cached signals
	c.mutex.Lock()
	c.cachedSignals = allSignals
	c.mutex.Unlock()
}

// GetCachedSignals returns the last cached signals
func (c *Coordinator) GetCachedSignals() []Signal {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	signals := make([]Signal, len(c.cachedSignals))
	copy(signals, c.cachedSignals)
	return signals
}

// Config returns the coordinator configuration
func (c *Coordinator) GetConfig() *Config {
	return c.config
}
