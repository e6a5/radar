package radar

import (
	"math"
	"sync"
	"time"
)

// PerformanceMonitor tracks rendering and processing performance
type PerformanceMonitor struct {
	mutex           sync.RWMutex
	frameCount      int
	lastFrameTime   time.Time
	frameRate       float64
	renderLatency   time.Duration
	lastRenderStart time.Time
	avgRenderTime   time.Duration
	totalRenderTime time.Duration
	enabled         bool
}

// SpatialCache caches expensive spatial calculations
type SpatialCache struct {
	mutex       sync.RWMutex
	circleCache map[float64][]Point
	sinCache    map[float64]float64
	cosCache    map[float64]float64
	maxEntries  int
	enabled     bool
}

// Point represents a cached coordinate point
type Point struct {
	X, Y int
}

// PerformanceStats contains current performance metrics
type PerformanceStats struct {
	FrameRate     float64
	RenderLatency time.Duration
	AvgRenderTime time.Duration
	FrameCount    int
}

// NewPerformanceMonitor creates a new performance monitor
func NewPerformanceMonitor() *PerformanceMonitor {
	return &PerformanceMonitor{
		enabled:       true,
		lastFrameTime: time.Now(),
	}
}

// NewSpatialCache creates a new spatial calculation cache
func NewSpatialCache(maxEntries int) *SpatialCache {
	return &SpatialCache{
		circleCache: make(map[float64][]Point),
		sinCache:    make(map[float64]float64),
		cosCache:    make(map[float64]float64),
		maxEntries:  maxEntries,
		enabled:     true,
	}
}

// StartFrame begins performance measurement for a frame
func (pm *PerformanceMonitor) StartFrame() {
	if !pm.enabled {
		return
	}

	pm.lastRenderStart = time.Now()
}

// EndFrame completes performance measurement for a frame
func (pm *PerformanceMonitor) EndFrame() {
	if !pm.enabled {
		return
	}

	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	now := time.Now()

	// Calculate render latency
	pm.renderLatency = now.Sub(pm.lastRenderStart)
	pm.totalRenderTime += pm.renderLatency
	pm.frameCount++

	// Calculate average render time
	pm.avgRenderTime = pm.totalRenderTime / time.Duration(pm.frameCount)

	// Calculate frame rate
	if !pm.lastFrameTime.IsZero() {
		frameTime := now.Sub(pm.lastFrameTime).Seconds()
		if frameTime > 0 {
			pm.frameRate = 1.0 / frameTime
		}
	}
	pm.lastFrameTime = now
}

// GetStats returns current performance statistics
func (pm *PerformanceMonitor) GetStats() PerformanceStats {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	return PerformanceStats{
		FrameRate:     pm.frameRate,
		RenderLatency: pm.renderLatency,
		AvgRenderTime: pm.avgRenderTime,
		FrameCount:    pm.frameCount,
	}
}

// IsPerformanceGood returns true if performance is acceptable
func (pm *PerformanceMonitor) IsPerformanceGood() bool {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	// Consider performance good if:
	// - Frame rate > 10 FPS
	// - Render latency < 50ms
	return pm.frameRate > 10.0 && pm.renderLatency < 50*time.Millisecond
}

// GetCirclePoints returns cached circle points or calculates them
func (sc *SpatialCache) GetCirclePoints(radius float64) []Point {
	if !sc.enabled {
		return sc.calculateCirclePoints(radius)
	}

	sc.mutex.RLock()
	if points, exists := sc.circleCache[radius]; exists {
		sc.mutex.RUnlock()
		return points
	}
	sc.mutex.RUnlock()

	// Calculate and cache
	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	// Double-check after acquiring write lock
	if points, exists := sc.circleCache[radius]; exists {
		return points
	}

	// Limit cache size
	if len(sc.circleCache) >= sc.maxEntries {
		// Remove a random entry to make space
		for k := range sc.circleCache {
			delete(sc.circleCache, k)
			break
		}
	}

	points := sc.calculateCirclePoints(radius)
	sc.circleCache[radius] = points
	return points
}

// calculateCirclePoints computes circle points for given radius
func (sc *SpatialCache) calculateCirclePoints(radius float64) []Point {
	points := make([]Point, 0, 128) // Pre-allocate reasonable capacity

	// Use adaptive step size based on radius
	stepSize := 0.1
	if radius < 10 {
		stepSize = 0.2
	} else if radius > 50 {
		stepSize = 0.05
	}

	for angle := 0.0; angle < 2*math.Pi; angle += stepSize {
		x := int(math.Round(radius * sc.getCos(angle)))
		y := int(math.Round(radius * sc.getSin(angle) * 0.5)) // Terminal aspect ratio
		points = append(points, Point{X: x, Y: y})
	}

	return points
}

// getCos returns cached cosine value or calculates it
func (sc *SpatialCache) getCos(angle float64) float64 {
	// Round angle to reduce cache misses
	rounded := math.Round(angle*100) / 100

	sc.mutex.RLock()
	if val, exists := sc.cosCache[rounded]; exists {
		sc.mutex.RUnlock()
		return val
	}
	sc.mutex.RUnlock()

	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	// Double-check after acquiring write lock
	if val, exists := sc.cosCache[rounded]; exists {
		return val
	}

	// Limit cache size
	if len(sc.cosCache) >= sc.maxEntries {
		// Remove oldest entry
		for k := range sc.cosCache {
			delete(sc.cosCache, k)
			break
		}
	}

	val := math.Cos(rounded)
	sc.cosCache[rounded] = val
	return val
}

// getSin returns cached sine value or calculates it
func (sc *SpatialCache) getSin(angle float64) float64 {
	// Round angle to reduce cache misses
	rounded := math.Round(angle*100) / 100

	sc.mutex.RLock()
	if val, exists := sc.sinCache[rounded]; exists {
		sc.mutex.RUnlock()
		return val
	}
	sc.mutex.RUnlock()

	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	// Double-check after acquiring write lock
	if val, exists := sc.sinCache[rounded]; exists {
		return val
	}

	// Limit cache size
	if len(sc.sinCache) >= sc.maxEntries {
		// Remove oldest entry
		for k := range sc.sinCache {
			delete(sc.sinCache, k)
			break
		}
	}

	val := math.Sin(rounded)
	sc.sinCache[rounded] = val
	return val
}

// ClearCache clears all cached data
func (sc *SpatialCache) ClearCache() {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	sc.circleCache = make(map[float64][]Point)
	sc.sinCache = make(map[float64]float64)
	sc.cosCache = make(map[float64]float64)
}

// GetCacheStats returns cache utilization statistics
func (sc *SpatialCache) GetCacheStats() (circleEntries, sinEntries, cosEntries int) {
	sc.mutex.RLock()
	defer sc.mutex.RUnlock()

	return len(sc.circleCache), len(sc.sinCache), len(sc.cosCache)
}
