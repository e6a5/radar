package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/e6a5/radar/radar"
	"github.com/gdamore/tcell/v2"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// Check for existing consent or ask for permission to collect real data
	if !hasConsent() && !askForPermission() {
		fmt.Println("Permission denied. Exiting.")
		os.Exit(0)
	}

	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("Error creating screen: %v", err)
	}
	if err := screen.Init(); err != nil {
		log.Fatalf("Error initializing screen: %v", err)
	}
	defer screen.Fini()

	// Get initial terminal size
	width, height := screen.Size()
	display := radar.NewDisplay(width, height)

	// Main loop
	for {
		if !display.HandleInput(screen) {
			break
		}

		display.Render(screen)
		display.UpdatePhases()
		time.Sleep(display.RefreshRate())
	}
}

// askForPermission prompts the user for permission to collect real data
func askForPermission() bool {
	fmt.Println("🎯 Radar v2.0 - Real-time Network Monitoring")
	fmt.Println("============================================")
	fmt.Println()
	fmt.Println("This application will collect real data from your device:")
	fmt.Println("  • WiFi network scanning (network names and signal strength)")
	fmt.Println("  • Network interface monitoring (active connections)")
	fmt.Println("  • System process monitoring (network-related processes)")
	fmt.Println("  • Local network device discovery (ping scanning)")
	fmt.Println()
	fmt.Println("No data is transmitted or stored externally.")
	fmt.Println("All data remains local on your device.")
	fmt.Println()
	fmt.Print("Do you give permission to collect this data? (Y/n): ")

	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error reading input: %v\n", err)
		return false
	}

	response = strings.TrimSpace(strings.ToLower(response))

	// Accept Y, y, yes, YES, or empty (default to yes)
	if response == "y" || response == "yes" || response == "" {
		// Store consent for future runs
		saveConsent()
		return true
	}

	return false
}

// getConsentFilePath returns the path to the consent file
func getConsentFilePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory if can't get home dir
		return ".radar_consent"
	}
	return filepath.Join(homeDir, ".radar_consent")
}

// hasConsent checks if user has previously given consent
func hasConsent() bool {
	consentFile := getConsentFilePath()
	if _, err := os.Stat(consentFile); err == nil {
		return true
	}
	return false
}

// saveConsent saves the user's consent to a file
func saveConsent() {
	consentFile := getConsentFilePath()
	consent := fmt.Sprintf("# Radar Data Collection Consent\n# Generated: %s\n# User has consented to real data collection\nconsent=true\n", time.Now().Format(time.RFC3339))

	err := os.WriteFile(consentFile, []byte(consent), 0644)
	if err != nil {
		fmt.Printf("Warning: Could not save consent preferences: %v\n", err)
	}
}
