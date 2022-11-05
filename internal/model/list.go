package model

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/orlangure/gocovsh/internal/styles"
	"golang.org/x/tools/cover"
)

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2).MarginTop(1)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2)
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	statusBarStyle    = lipgloss.NewStyle().MarginLeft(4)
	percentageStyle   = lipgloss.NewStyle().PaddingLeft(1)
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
		line = selectedItemStyle.Foreground(lipgloss.Color(styles.CurrentTheme.PrimaryColor)).Render("> " + line)
	} else {
		line = itemStyle.Render(line)
	}

	fmt.Fprint(w, line)
}

func (d coverProfileDelegate) renderBaseLine(p *coverProfile) string {
	percentage := percentageStyle.Foreground(lipgloss.Color(styles.CurrentTheme.InactiveColor)).Render(fmt.Sprintf("%.2f%%", p.percentage))
	return fmt.Sprintf("%s %s", p.profile.FileName, percentage)
}
