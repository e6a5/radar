package wifi

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// NetworkMapping represents a mapping from SSID pattern to friendly name
type NetworkMapping struct {
	Pattern     *regexp.Regexp
	Replacement string
	Type        string // "home", "business", "public", "carrier", "device"
}

// Common network mappings for better display names
var networkMappings = []NetworkMapping{
	// Home Networks
	{regexp.MustCompile(`^NETGEAR\d+$`), "Home Network (Netgear)", "home"},
	{regexp.MustCompile(`^Linksys\d+$`), "Home Network (Linksys)", "home"},
	{regexp.MustCompile(`^TP-Link_[A-F0-9]+$`), "Home Network (TP-Link)", "home"},
	{regexp.MustCompile(`^ASUS_[A-F0-9]+$`), "Home Network (ASUS)", "home"},
	{regexp.MustCompile(`^(WiFi|Network|Internet)_?\d*$`), "Home Network", "home"},
	{regexp.MustCompile(`^HOME-\w+$`), "Home Network", "home"},
	{regexp.MustCompile(`^(Family|House|My)WiFi$`), "Family Network", "home"},

	// Business Networks
	{regexp.MustCompile(`^(Guest|Visitor|Public)$`), "Guest Network", "business"},
	{regexp.MustCompile(`^(Office|Work|Company)[\w-]*$`), "Office Network", "business"},
	{regexp.MustCompile(`^[A-Z]{2,}-Office$`), "Office Network", "business"},
	{regexp.MustCompile(`^Corp[\w-]*$`), "Corporate Network", "business"},
	{regexp.MustCompile(`^Enterprise[\w-]*$`), "Enterprise Network", "business"},

	// Public WiFi
	{regexp.MustCompile(`^Starbucks$`), "Starbucks WiFi", "public"},
	{regexp.MustCompile(`^McDonalds$`), "McDonald's WiFi", "public"},
	{regexp.MustCompile(`^(Subway|SUBWAY)[\w-]*$`), "Subway WiFi", "public"},
	{regexp.MustCompile(`^(Free|Public)[\w-]*WiFi$`), "Public WiFi", "public"},
	{regexp.MustCompile(`^Hotel[\w-]*$`), "Hotel WiFi", "public"},
	{regexp.MustCompile(`^Airport[\w-]*$`), "Airport WiFi", "public"},
	{regexp.MustCompile(`^Coffee[\w-]*$`), "Coffee Shop WiFi", "public"},
	{regexp.MustCompile(`^Library[\w-]*$`), "Library WiFi", "public"},
	{regexp.MustCompile(`^Mall[\w-]*$`), "Shopping Center WiFi", "public"},

	// Carrier Networks
	{regexp.MustCompile(`^(Verizon|VZW)[\w-]*$`), "Verizon Network", "carrier"},
	{regexp.MustCompile(`^(AT&T|ATT)[\w-]*$`), "AT&T Network", "carrier"},
	{regexp.MustCompile(`^T-Mobile[\w-]*$`), "T-Mobile Network", "carrier"},
	{regexp.MustCompile(`^Sprint[\w-]*$`), "Sprint Network", "carrier"},
	{regexp.MustCompile(`^(Xfinity|Comcast)[\w-]*$`), "Xfinity Hotspot", "carrier"},
	{regexp.MustCompile(`^Spectrum[\w-]*$`), "Spectrum WiFi", "carrier"},

	// Device Networks
	{regexp.MustCompile(`^\w+['']s iPhone$`), "iPhone Hotspot", "device"},
	{regexp.MustCompile(`^\w+['']s iPad$`), "iPad Hotspot", "device"},
	{regexp.MustCompile(`^iPhone$`), "iPhone Hotspot", "device"},
	{regexp.MustCompile(`^iPad$`), "iPad Hotspot", "device"},
	{regexp.MustCompile(`^AndroidAP[\w-]*$`), "Android Hotspot", "device"},
	{regexp.MustCompile(`^Galaxy[\w-]*$`), "Samsung Galaxy Hotspot", "device"},
	{regexp.MustCompile(`^Pixel[\w-]*$`), "Google Pixel Hotspot", "device"},

	// Generic patterns
	{regexp.MustCompile(`^[A-F0-9]{12}$`), "Router Network", "home"},
	{regexp.MustCompile(`^[A-F0-9]{6}$`), "Default Network", "home"},
	{regexp.MustCompile(`^(Router|Modem)[\w-]*$`), "Home Router", "home"},
}

// GetFriendlyDisplayName converts a raw SSID to a human-readable display name
func GetFriendlyDisplayName(ssid string, strength int, isConnected bool) string {
	if ssid == "" {
		return "Unknown Network"
	}

	// Clean the SSID
	cleanSSID := strings.TrimSpace(ssid)

	// Check for hidden networks
	if cleanSSID == "" || cleanSSID == "<hidden>" || cleanSSID == "Hidden Network" {
		return fmt.Sprintf("Hidden Network (%d%%)", strength)
	}

	// Try to match against known patterns
	displayName, networkType := mapSSIDToFriendlyName(cleanSSID)

	// Add signal strength indicator for weak signals
	if strength < 30 {
		displayName += " (Weak)"
	} else if strength > 80 {
		displayName += " (Strong)"
	}

	// Add connection status
	if isConnected {
		displayName += " ★"
	}

	// Add network type emoji/indicator
	switch networkType {
	case "home":
		displayName = "🏠 " + displayName
	case "business":
		displayName = "🏢 " + displayName
	case "public":
		displayName = "📶 " + displayName
	case "carrier":
		displayName = "📡 " + displayName
	case "device":
		displayName = "📱 " + displayName
	default:
		// For unknown networks, try to infer type from characteristics
		if isLikelyHomeNetwork(cleanSSID) {
			displayName = "🏠 " + displayName
		} else if isLikelyPublicNetwork(cleanSSID) {
			displayName = "📶 " + displayName
		} else {
			displayName = "📡 " + displayName
		}
	}

	return displayName
}

// mapSSIDToFriendlyName tries to map an SSID to a friendly name using patterns
func mapSSIDToFriendlyName(ssid string) (string, string) {
	// Try exact mappings first
	for _, mapping := range networkMappings {
		if mapping.Pattern.MatchString(ssid) {
			return mapping.Replacement, mapping.Type
		}
	}

	// If no pattern matches, try to make it more readable
	friendlyName := makeMoreReadable(ssid)
	return friendlyName, "unknown"
}

// makeMoreReadable improves the readability of SSIDs
func makeMoreReadable(ssid string) string {
	// Replace underscores and hyphens with spaces
	readable := strings.ReplaceAll(ssid, "_", " ")
	readable = strings.ReplaceAll(readable, "-", " ")

	// Handle camelCase - add spaces before uppercase letters
	var result []rune
	for i, r := range readable {
		if i > 0 && unicode.IsUpper(r) && unicode.IsLower(rune(readable[i-1])) {
			result = append(result, ' ')
		}
		result = append(result, r)
	}
	readable = string(result)

	// Capitalize first letter of each word
	words := strings.Fields(readable)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
		}
	}

	return strings.Join(words, " ")
}

// isLikelyHomeNetwork tries to determine if an SSID represents a home network
func isLikelyHomeNetwork(ssid string) bool {
	homeIndicators := []string{
		"home", "house", "family", "wifi", "network", "internet",
		"router", "modem", "default", "linksys", "netgear", "dlink",
		"tplink", "asus", "belkin", "cisco",
	}

	lowerSSID := strings.ToLower(ssid)
	for _, indicator := range homeIndicators {
		if strings.Contains(lowerSSID, indicator) {
			return true
		}
	}

	// Check for MAC address patterns (often default router names)
	if matched, _ := regexp.MatchString(`[A-Fa-f0-9]{12}`, ssid); matched {
		return true
	}

	return false
}

// isLikelyPublicNetwork tries to determine if an SSID represents a public network
func isLikelyPublicNetwork(ssid string) bool {
	publicIndicators := []string{
		"guest", "public", "free", "open", "wifi", "hotspot",
		"starbucks", "mcdonalds", "subway", "hotel", "airport",
		"library", "coffee", "cafe", "restaurant", "mall",
		"store", "shop", "center", "plaza",
	}

	lowerSSID := strings.ToLower(ssid)
	for _, indicator := range publicIndicators {
		if strings.Contains(lowerSSID, indicator) {
			return true
		}
	}

	return false
}

// GetNetworkTypeDescription returns a description of the network type
func GetNetworkTypeDescription(ssid string) string {
	_, networkType := mapSSIDToFriendlyName(ssid)

	switch networkType {
	case "home":
		return "Home/Personal Network"
	case "business":
		return "Business Network"
	case "public":
		return "Public WiFi"
	case "carrier":
		return "Carrier Network"
	case "device":
		return "Mobile Hotspot"
	default:
		if isLikelyHomeNetwork(ssid) {
			return "Likely Home Network"
		} else if isLikelyPublicNetwork(ssid) {
			return "Likely Public Network"
		}
		return "Unknown Network Type"
	}
}
