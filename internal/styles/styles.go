package styles

import (
	"os"

	catppuccin "github.com/catppuccin/go"
	"github.com/charmbracelet/lipgloss"
)

var CurrentTheme Theme

type Theme struct {
	PrimaryColor   string
	SecondaryColor string
	InactiveColor  string

	NeutralLine   lipgloss.Style
	CoveredLine   lipgloss.Style
	UncoveredLine lipgloss.Style
}

func (t *Theme) setStyles() {
	t.NeutralLine = lipgloss.NewStyle().Foreground(lipgloss.Color(t.InactiveColor))
	t.CoveredLine = lipgloss.NewStyle().Foreground(lipgloss.Color(t.PrimaryColor))
	t.UncoveredLine = lipgloss.NewStyle().Foreground(lipgloss.Color(t.SecondaryColor))
}

func Default() Theme {
	t := Theme{
		PrimaryColor:   "#00ff00",
		SecondaryColor: "#ff0000",
		InactiveColor:  "#7f7f7f",
	}
	t.setStyles()

	return t
}

func Catppuccin(cpn catppuccin.Theme) Theme {
	t := Theme{
		PrimaryColor:   cpn.Green().Hex,
		SecondaryColor: cpn.Red().Hex,
		InactiveColor:  cpn.Subtext1().Hex,
	}
	t.setStyles()

	return t
}

func SetTheme() {
	theme := os.Getenv("GOCOVSH_THEME")

	if variant := catppuccin.Variant(theme); variant != nil {
		CurrentTheme = Catppuccin(variant)
	} else {
		CurrentTheme = Default()
	}
}
