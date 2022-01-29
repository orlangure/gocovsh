package main

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
)

type coverProfile struct {
	profile *cover.Profile
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

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s string) string {
			return selectedItemStyle.Render("> " + s)
		}
	}

	fmt.Fprint(w, fn(profile.profile.FileName))
}
