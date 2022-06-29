package styles

import "github.com/charmbracelet/lipgloss"

const (
	// TODO: support themes + dark/light mode.
	PrimaryColor   = "#00ff00"
	SecondaryColor = "#ff0000"
	InactiveColor  = "#7f7f7f"
)

var (
	NeutralLine   = lipgloss.NewStyle().Foreground(lipgloss.Color(InactiveColor))
	CoveredLine   = lipgloss.NewStyle().Foreground(lipgloss.Color(PrimaryColor))
	UncoveredLine = lipgloss.NewStyle().Foreground(lipgloss.Color(SecondaryColor))
)
