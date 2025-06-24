package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/e6a5/radar/radar"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	
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
