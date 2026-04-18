package ui

import "github.com/charmbracelet/lipgloss"

var (
	// Sovereign palette — pure black base, zero gray/charcoal.
	ColorNeonPink  = lipgloss.Color("#FF2A6D")
	ColorNeonCyan  = lipgloss.Color("#05D9E8")
	ColorNeonGreen = lipgloss.Color("#01FFC3")
	ColorGold      = lipgloss.Color("#FFD700")
	// Replaces all steel/gray — inactive pins glow faintly teal on pure black.
	ColorDimCyan   = lipgloss.Color("#0D5252")
	ColorPureBlack = lipgloss.Color("#000000")

	MainFrame = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(ColorNeonCyan).
			Padding(1, 2).
			Background(ColorPureBlack)

	TitleStyle  = lipgloss.NewStyle().Foreground(ColorPureBlack).Background(ColorNeonCyan).Bold(true).Padding(0, 3)
	WinStyle    = lipgloss.NewStyle().Foreground(ColorNeonGreen).Bold(true)
	LossStyle   = lipgloss.NewStyle().Foreground(ColorNeonPink).Bold(true)
	GoldStyle   = lipgloss.NewStyle().Foreground(ColorGold).Bold(true)
	PinkStyle   = lipgloss.NewStyle().Foreground(ColorNeonPink).Bold(true)
	CyanStyle   = lipgloss.NewStyle().Foreground(ColorNeonCyan).Bold(true)
	MutedStyle  = lipgloss.NewStyle().Foreground(ColorDimCyan)
	HeaderStyle = lipgloss.NewStyle().Foreground(ColorNeonCyan).Bold(true).Underline(true)
)
