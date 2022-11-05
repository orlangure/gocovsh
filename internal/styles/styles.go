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

func CatppuccinMocha() Theme {
	PrimaryColor := catppuccin.Mocha().Green().Hex
	SecondaryColor := catppuccin.Mocha().Red().Hex
	InactiveColor := catppuccin.Mocha().Subtext1().Hex
	t := Theme{PrimaryColor: PrimaryColor, SecondaryColor: SecondaryColor, InactiveColor: InactiveColor}
	t.setStyles()

	return t
}

func CatppuccinMacchiato() Theme {
	PrimaryColor := catppuccin.Macchiato().Green().Hex
	SecondaryColor := catppuccin.Macchiato().Red().Hex
	InactiveColor := catppuccin.Macchiato().Subtext1().Hex
	t := Theme{PrimaryColor: PrimaryColor, SecondaryColor: SecondaryColor, InactiveColor: InactiveColor}
	t.setStyles()

	return t
}

func CatppuccinLatte() Theme {
	PrimaryColor := catppuccin.Latte().Green().Hex
	SecondaryColor := catppuccin.Latte().Red().Hex
	InactiveColor := catppuccin.Latte().Subtext1().Hex
	t := Theme{PrimaryColor: PrimaryColor, SecondaryColor: SecondaryColor, InactiveColor: InactiveColor}
	t.setStyles()

	return t
}

func Default() Theme {
	PrimaryColor := "#00ff00"
	SecondaryColor := "#ff0000"
	InactiveColor := "#7f7f7f"
	t := Theme{PrimaryColor: PrimaryColor, SecondaryColor: SecondaryColor, InactiveColor: InactiveColor}
	t.setStyles()

	return t
}

func SetTheme() {
	theme := os.Getenv("GOCOVSH_THEME")
	switch theme {
	case "mocha":
		CurrentTheme = CatppuccinMocha()
	case "macchiato":
		CurrentTheme = CatppuccinMacchiato()
	case "latte":
		CurrentTheme = CatppuccinLatte()
	default:
		CurrentTheme = Default()
	}
}
