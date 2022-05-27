package model

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/tools/cover"
)

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2).MarginTop(1)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color(primaryColor))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	statusBarStyle    = lipgloss.NewStyle().MarginLeft(4).Foreground(lipgloss.Color(inactiveColor))
	percentageStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color(inactiveColor)).PaddingLeft(1)
)

type coverProfile struct {
	profile    *cover.Profile
	percentage float64
}

func (f *coverProfile) FilterValue() string { return f.profile.FileName }

type coverProfileDelegate struct{}

func (d coverProfileDelegate) Height() int                               { return 1 }
func (d coverProfileDelegate) Spacing() int                              { return 0 }
func (d coverProfileDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d coverProfileDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	profile, ok := listItem.(*coverProfile)
	if !ok {
		return
	}

	line := d.renderBaseLine(profile)

	if index == m.Index() {
		line = selectedItemStyle.Render("> " + line)
	} else {
		line = itemStyle.Render(line)
	}

	fmt.Fprint(w, line)
}

func (d coverProfileDelegate) renderBaseLine(p *coverProfile) string {
	percentageString := fmt.Sprintf("%3d%%", int(p.percentage))
	percentage := percentageStyle.Render(percentageString)
	return fmt.Sprintf("%s %s", percentage, p.profile.FileName)
}
