package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// Color palette
var (
	// Primary colors
	colorRed     = lipgloss.Color("#FF6B6B")
	colorGreen   = lipgloss.Color("#69DB7C")
	colorYellow  = lipgloss.Color("#FFD43B")
	colorBlue    = lipgloss.Color("#74C0FC")
	colorGray    = lipgloss.Color("#868E96")
	colorDimGray = lipgloss.Color("#495057")

	// Icons
	iconError   = "✗"
	iconSuccess = "✓"
	iconWarning = "!"
	iconInfo    = "→"
	iconHint    = "💡"
)

// Text styles
var (
	// Error styles
	errorIcon = lipgloss.NewStyle().
			Foreground(colorRed).
			Bold(true)

	errorTitle = lipgloss.NewStyle().
			Foreground(colorRed).
			Bold(true)

	errorMessage = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#DEE2E6"))

	// Success styles
	successIcon = lipgloss.NewStyle().
			Foreground(colorGreen).
			Bold(true)

	successTitle = lipgloss.NewStyle().
			Foreground(colorGreen).
			Bold(true)

	successMessage = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#DEE2E6"))

	// Warning styles
	warningIcon = lipgloss.NewStyle().
			Foreground(colorYellow).
			Bold(true)

	warningTitle = lipgloss.NewStyle().
			Foreground(colorYellow).
			Bold(true)

	// Info styles
	infoIcon = lipgloss.NewStyle().
			Foreground(colorBlue).
			Bold(true)

	infoMessage = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#DEE2E6"))

	// Hint style
	hintStyle = lipgloss.NewStyle().
			Foreground(colorDimGray).
			Italic(true)

	// Path style - for file/directory paths
	pathStyle = lipgloss.NewStyle().
			Foreground(colorBlue).
			Bold(true)

	// Command style - for CLI commands
	cmdStyle = lipgloss.NewStyle().
			Foreground(colorYellow)

	// Dim text
	dimStyle = lipgloss.NewStyle().
			Foreground(colorGray)
)

// Box styles for framed messages
var (
	errorBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorRed).
			Padding(1, 2)

	successBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorGreen).
			Padding(1, 2)

	warningBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorYellow).
			Padding(1, 2)

	infoBox = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorBlue).
		Padding(1, 2)
)
