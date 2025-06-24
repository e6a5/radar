package radar

import "github.com/gdamore/tcell/v2"

// Theme selection for enhanced visuals
type ThemeType int

const (
	ThemeModernDark ThemeType = iota
	ThemeClassicGreen
	ThemeBlueNeon
	ThemeMilitary
)

// Get the appropriate theme based on selection
func GetRadarTheme(themeType ThemeType) RadarTheme {
	switch themeType {
	case ThemeClassicGreen:
		return getClassicGreenTheme()
	case ThemeBlueNeon:
		return getBlueNeonTheme()
	case ThemeMilitary:
		return getMilitaryTheme()
	default:
		return GetModernDarkTheme()
	}
}

// Classic green radar theme (traditional radar look)
func getClassicGreenTheme() RadarTheme {
	return RadarTheme{
		Background:    tcell.ColorBlack,
		GridPrimary:   tcell.ColorDarkGreen,
		GridSecondary: tcell.Color16,

		RingPrimary:   tcell.ColorGreen,
		RingSecondary: tcell.ColorDarkGreen,
		RingLabels:    tcell.ColorGreen,

		SweepPrimary:   tcell.ColorGreen,
		SweepSecondary: tcell.ColorLime,
		SweepFade:      tcell.ColorDarkGreen,
		SweepTrail:     tcell.ColorDarkSlateGray,

		SignalExcellent: tcell.ColorYellow,
		SignalGood:      tcell.ColorLime,
		SignalFair:      tcell.ColorGreen,
		SignalPoor:      tcell.ColorDarkGreen,
		SignalConnected: tcell.ColorWhite,

		AccentPrimary:   tcell.ColorGreen,
		AccentSecondary: tcell.ColorLime,
		TextPrimary:     tcell.ColorWhite,
		TextSecondary:   tcell.ColorGreen,
	}
}

// Blue neon theme (cyberpunk aesthetic)
func getBlueNeonTheme() RadarTheme {
	return RadarTheme{
		Background:    tcell.ColorBlack,
		GridPrimary:   tcell.ColorNavy,
		GridSecondary: tcell.Color16,

		RingPrimary:   tcell.ColorBlue,
		RingSecondary: tcell.ColorNavy,
		RingLabels:    tcell.ColorBlue,

		SweepPrimary:   tcell.ColorBlue,
		SweepSecondary: tcell.ColorDarkBlue,
		SweepFade:      tcell.ColorNavy,
		SweepTrail:     tcell.ColorDarkSlateGray,

		SignalExcellent: tcell.ColorRed,
		SignalGood:      tcell.ColorPurple,
		SignalFair:      tcell.ColorBlue,
		SignalPoor:      tcell.ColorNavy,
		SignalConnected: tcell.ColorWhite,

		AccentPrimary:   tcell.ColorBlue,
		AccentSecondary: tcell.ColorDarkBlue,
		TextPrimary:     tcell.ColorWhite,
		TextSecondary:   tcell.ColorBlue,
	}
}

// Military theme (tactical display)
func getMilitaryTheme() RadarTheme {
	return RadarTheme{
		Background:    tcell.ColorBlack,
		GridPrimary:   tcell.ColorMaroon,
		GridSecondary: tcell.Color16,

		RingPrimary:   tcell.ColorOlive,
		RingSecondary: tcell.ColorMaroon,
		RingLabels:    tcell.ColorYellow,

		SweepPrimary:   tcell.ColorOlive,
		SweepSecondary: tcell.ColorYellow,
		SweepFade:      tcell.ColorMaroon,
		SweepTrail:     tcell.ColorDarkSlateGray,

		SignalExcellent: tcell.ColorRed,
		SignalGood:      tcell.ColorOrange,
		SignalFair:      tcell.ColorYellow,
		SignalPoor:      tcell.ColorGray,
		SignalConnected: tcell.ColorLime,

		AccentPrimary:   tcell.ColorOlive,
		AccentSecondary: tcell.ColorYellow,
		TextPrimary:     tcell.ColorWhite,
		TextSecondary:   tcell.ColorYellow,
	}
}

// Theme names for display
func GetThemeName(themeType ThemeType) string {
	switch themeType {
	case ThemeClassicGreen:
		return "Classic Green"
	case ThemeBlueNeon:
		return "Blue Neon"
	case ThemeMilitary:
		return "Military"
	default:
		return "Modern Dark"
	}
}
