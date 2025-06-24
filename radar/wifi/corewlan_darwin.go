//go:build darwin
// +build darwin

package wifi

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework CoreWLAN

#import <Foundation/Foundation.h>
#import <CoreWLAN/CoreWLAN.h>

typedef struct {
    char* ssid;
    int rssi;
    char* bssid;
} WiFiNetwork;

typedef struct {
    WiFiNetwork* networks;
    int count;
} WiFiScanResult;

WiFiScanResult* scanWiFiNetworks() {
    @autoreleasepool {
        CWWiFiClient* client = [CWWiFiClient sharedWiFiClient];
        CWInterface* interface = [client interface];

        if (!interface) {
            return NULL;
        }

        NSError* error = nil;
        NSSet<CWNetwork*>* networks = [interface scanForNetworksWithName:nil error:&error];

        if (error || !networks) {
            return NULL;
        }

        WiFiScanResult* result = malloc(sizeof(WiFiScanResult));
        result->count = (int)[networks count];
        result->networks = malloc(sizeof(WiFiNetwork) * result->count);

        int i = 0;
        for (CWNetwork* network in networks) {
            NSString* ssid = [network ssid];
            if (ssid) {
                const char* ssidCStr = [ssid UTF8String];
                result->networks[i].ssid = malloc(strlen(ssidCStr) + 1);
                strcpy(result->networks[i].ssid, ssidCStr);
            } else {
                result->networks[i].ssid = malloc(1);
                result->networks[i].ssid[0] = '\0';
            }

            result->networks[i].rssi = (int)[network rssiValue];

            NSString* bssid = [network bssid];
            if (bssid) {
                const char* bssidCStr = [bssid UTF8String];
                result->networks[i].bssid = malloc(strlen(bssidCStr) + 1);
                strcpy(result->networks[i].bssid, bssidCStr);
            } else {
                result->networks[i].bssid = malloc(1);
                result->networks[i].bssid[0] = '\0';
            }

            i++;
        }

        return result;
    }
}

void freeWiFiScanResult(WiFiScanResult* result) {
    if (!result) return;

    for (int i = 0; i < result->count; i++) {
        free(result->networks[i].ssid);
        free(result->networks[i].bssid);
    }
    free(result->networks);
    free(result);
}

WiFiNetwork* getCurrentWiFiNetwork() {
    @autoreleasepool {
        CWWiFiClient* client = [CWWiFiClient sharedWiFiClient];
        CWInterface* interface = [client interface];

        if (!interface) {
            return NULL;
        }

        NSString* ssid = [interface ssid];
        if (!ssid || [ssid length] == 0) {
            return NULL;
        }

        WiFiNetwork* result = malloc(sizeof(WiFiNetwork));

        const char* ssidCStr = [ssid UTF8String];
        result->ssid = malloc(strlen(ssidCStr) + 1);
        strcpy(result->ssid, ssidCStr);

        result->rssi = (int)[interface rssiValue];

        NSString* bssid = [interface bssid];
        if (bssid) {
            const char* bssidCStr = [bssid UTF8String];
            result->bssid = malloc(strlen(bssidCStr) + 1);
            strcpy(result->bssid, bssidCStr);
        } else {
            result->bssid = malloc(1);
            result->bssid[0] = '\0';
        }

        return result;
    }
}

void freeWiFiNetwork(WiFiNetwork* network) {
    if (!network) return;

    free(network->ssid);
    free(network->bssid);
    free(network);
}
*/
import "C"
import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"
	"unsafe"

	"github.com/e6a5/radar/radar/scanner"
	"github.com/gdamore/tcell/v2"
)

// CoreWLANScanner implements WiFi scanning using Apple's CoreWLAN framework
type CoreWLANScanner struct {
	lastScan time.Time
	config   *scanner.Config
}

// NewCoreWLANScanner creates a new CoreWLAN-based WiFi scanner
func NewCoreWLANScanner(config *scanner.Config) *CoreWLANScanner {
	return &CoreWLANScanner{
		config: config,
	}
}

// Name returns the scanner name
func (c *CoreWLANScanner) Name() string {
	return "CoreWLAN WiFi Scanner"
}

// IsAvailable checks if CoreWLAN is available (always true on Darwin)
func (c *CoreWLANScanner) IsAvailable() bool {
	return true
}

// Scan scans for available WiFi networks using CoreWLAN
func (c *CoreWLANScanner) Scan(ctx context.Context) ([]scanner.Signal, error) {
	signals := make([]scanner.Signal, 0)
	now := time.Now()

	// Rate limiting
	if now.Sub(c.lastScan) < c.config.ScanInterval {
		return signals, nil
	}
	c.lastScan = now

	// Scan for networks
	result := C.scanWiFiNetworks()
	if result == nil {
		return signals, fmt.Errorf("CoreWLAN scan failed")
	}
	defer C.freeWiFiScanResult(result)

	// Convert C results to Go signals
	networkCount := int(result.count)
	if networkCount == 0 {
		return signals, nil
	}

	networks := (*[1024]C.WiFiNetwork)(unsafe.Pointer(result.networks))[:networkCount:networkCount]

	for i := 0; i < networkCount && len(signals) < c.config.MaxSignals; i++ {
		network := networks[i]

		ssid := C.GoString(network.ssid)
		if ssid == "" {
			continue
		}

		rssi := int(network.rssi)
		strength := rssiToStrength(rssi)
		distance := rssiToDistance(rssi, c.config.MaxScanRange)

		// Get friendly display name
		displayName := GetFriendlyDisplayName(ssid, strength, false)

		signal := scanner.Signal{
			Type:        "WiFi",
			Icon:        "≋",
			Name:        displayName,
			Color:       tcell.ColorBlue,
			Strength:    strength,
			Distance:    distance,
			Angle:       rand.Float64() * 2 * math.Pi,
			Phase:       0,
			Lifetime:    now,
			LastSeen:    now,
			Persistence: 1.0,
			History:     make([]scanner.PositionHistory, 0, 20),
			MaxHistory:  20,
		}

		signal.AddToHistory(signal.Distance, signal.Angle, signal.Strength, true, now)
		signals = append(signals, signal)
	}

	// Add current network if connected
	if currentNetwork := c.getCurrentNetwork(now); currentNetwork != nil {
		signals = append([]scanner.Signal{*currentNetwork}, signals...)
	}

	return signals, nil
}

// getCurrentNetwork gets the currently connected WiFi network
func (c *CoreWLANScanner) getCurrentNetwork(now time.Time) *scanner.Signal {
	network := C.getCurrentWiFiNetwork()
	if network == nil {
		return nil
	}
	defer C.freeWiFiNetwork(network)

	ssid := C.GoString(network.ssid)
	if ssid == "" {
		return nil
	}

	rssi := int(network.rssi)
	strength := rssiToStrength(rssi)
	distance := rssiToDistance(rssi, c.config.MaxScanRange)

	// Get friendly display name for connected network
	displayName := GetFriendlyDisplayName(ssid, strength, true)

	signal := scanner.Signal{
		Type:        "WiFi",
		Icon:        "≋",
		Name:        displayName,
		Color:       tcell.ColorGreen,
		Strength:    strength,
		Distance:    distance,
		Angle:       rand.Float64() * 2 * math.Pi,
		Phase:       0,
		Lifetime:    now,
		LastSeen:    now,
		Persistence: 1.0,
		History:     make([]scanner.PositionHistory, 0, 20),
		MaxHistory:  20,
	}

	signal.AddToHistory(signal.Distance, signal.Angle, signal.Strength, true, now)
	return &signal
}

// rssiToStrength converts RSSI to percentage
func rssiToStrength(rssi int) int {
	if rssi >= -30 {
		return 100
	}
	if rssi <= -90 {
		return 0
	}
	return int(100 * (float64(rssi+90) / 60.0))
}

// rssiToDistance converts RSSI to radar distance
func rssiToDistance(rssi int, maxRange float64) float64 {
	distance := math.Pow(10, float64(-rssi-30)/20.0) * 2.0
	if distance < 0.5 {
		distance = 0.5
	}
	if distance > maxRange {
		distance = maxRange
	}
	return distance
}
