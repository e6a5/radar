package main

import (
	"fmt"

	"github.com/e6a5/radar/radar/wifi"
)

func main() {
	fmt.Println("WiFi Display Name Mapping Examples:")
	fmt.Println("===================================")

	// Test various SSID patterns
	testCases := []struct {
		ssid      string
		strength  int
		connected bool
	}{
		{"NETGEAR42", 85, false},
		{"Linksys00123", 45, false},
		{"iPhone", 90, false},
		{"John's iPhone", 75, true},
		{"Starbucks", 60, false},
		{"Xfinity", 55, false},
		{"MyWiFi", 95, true},
		{"ASUS_2C5F43", 70, false},
		{"OfficeNetwork", 80, false},
		{"Guest", 40, false},
		{"AndroidAP1234", 65, false},
		{"CoffeeShop", 50, false},
		{"0123456789AB", 30, false},
		{"HOME-NET", 88, false},
		{"", 50, false},
		{"WirelessNetwork_5G", 72, false},
		{"CompanyGuest", 35, false},
	}

	for _, tc := range testCases {
		original := tc.ssid
		if original == "" {
			original = "<empty>"
		}

		friendly := wifi.GetFriendlyDisplayName(tc.ssid, tc.strength, tc.connected)
		networkType := wifi.GetNetworkTypeDescription(tc.ssid)

		fmt.Printf("Original: %-20s → Display: %-30s [%s]\n",
			original, friendly, networkType)
	}

	fmt.Println("\nKey Features:")
	fmt.Println("- 🏠 Home networks get house emoji")
	fmt.Println("- 🏢 Business networks get office emoji")
	fmt.Println("- 📶 Public WiFi gets signal emoji")
	fmt.Println("- 📡 Carrier networks get tower emoji")
	fmt.Println("- 📱 Mobile hotspots get phone emoji")
	fmt.Println("- ★ Connected networks get star indicator")
	fmt.Println("- (Weak)/(Strong) indicators for signal strength")
	fmt.Println("- Readable names from technical SSIDs")
}
