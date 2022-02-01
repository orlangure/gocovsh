// Package codeview provides a bubbletea component for displaying code. It adds
// line numbers and makes sure the code fits on the screen. It does not support
// wrapping long lines; instead, it trims them.
package codeview

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/truncate"
)

const (
	ellipsis        = "…"
	newLine         = "\n"
	lineNumberColor = "#505050"
)

var (
	fileTitleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return fileTitleStyle.Copy().BorderStyle(b)
	}()

	helpStyle = lipgloss.NewStyle().Padding(0, 0, 1, 4)

	// lineNumberStylePlaceholder should not be used directly; instead, a copy of this
	// style should be adjusted with the actual line width for better looks.
	lineNumberStylePlaceholder = lipgloss.NewStyle().
					Padding(0).
					Margin(0, 1, 0, 0).
					Foreground(lipgloss.Color(lineNumberColor)).
					Faint(true).
					Align(lipgloss.Right).
					BorderForeground(lipgloss.Color(lineNumberColor)).
					Border(lipgloss.NormalBorder(), false, true, false, false)
)

// New creates a new codeview model which is rendered into the provided width
// and height.
func New(width, height int) Model {
	return Model{
		viewport: viewport.New(width, height),
		help:     help.New(),
		showHelp: true,
	}
}

// Model is the codeview model. Use New to create a new instance.
type Model struct {
	viewport viewport.Model
	help     help.Model
	width    int
	height   int
	title    string
	lines    []string
	showHelp bool
}

// Update is used to update the internal model state based on the external
// events.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	// TODO: support number-based navigation <29-01-22, yury> //
	if msg, ok := msg.(tea.KeyMsg); ok {
		if key.Matches(msg, DefaultKeyMap.Home) {
			_ = m.viewport.GotoTop()
			return m, nil
		}

		if key.Matches(msg, DefaultKeyMap.End) {
			_ = m.viewport.GotoBottom()
			return m, nil
		}
	}

	vp, cmd := m.viewport.Update(msg)
	m.viewport = vp

	return m, cmd
}

// View renders the model to be displayed.
func (m *Model) View() string {
	headerView := m.headerView()
	footerView := m.footerView()
	helpView := m.helpView()
	codeView := m.viewport.View()

	sections := make([]string, 0, 4)
	sections = append(sections, headerView, codeView, footerView)

	if helpView != "" {
		sections = append(sections, helpView)
	}

	return strings.Join(sections, "\n")
}

// SetContent sets the content of the codeview.
func (m *Model) SetContent(lines []string) {
	// save the original lines to not lose content in case of window resizing
	m.lines = lines
	m.redrawLines()
	m.viewport.SetYOffset(0)
}

func (m *Model) redrawLines() {
	content := m.formatLines(m.lines)
	m.viewport.SetContent(content)
}

// SetWidth sets the width of the codeview.
func (m *Model) SetWidth(width int) {
	m.setSize(width, m.height)
}

// SetHeight sets the height of the codeview.
func (m *Model) SetHeight(height int) {
	m.setSize(m.width, height)
}

// Width returns the width of the codeview.
func (m *Model) Width() int {
	return m.width
}

// Height returns the height of the codeview.
func (m *Model) Height() int {
	return m.height
}

// SetTitle sets the title of the codeview, usually the file name.
func (m *Model) SetTitle(title string) {
	m.title = title
}

// ShortHelp implements  help.KeyMap interface.
func (m *Model) ShortHelp() []key.Binding {
	return []key.Binding{
		DefaultKeyMap.Up,
		DefaultKeyMap.Down,
		DefaultKeyMap.Home,
		DefaultKeyMap.End,
		DefaultKeyMap.Back,
	}
}

// FullHelp implements  help.KeyMap interface.
func (m *Model) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{DefaultKeyMap.Up, DefaultKeyMap.Down, DefaultKeyMap.Home, DefaultKeyMap.End},
		{DefaultKeyMap.HalfScreenDown, DefaultKeyMap.HalfScreenUp},
		{DefaultKeyMap.Back, DefaultKeyMap.Quit},
	}
}

// SetShowHelp allows to hide or show the help section.
func (m *Model) SetShowHelp(showHelp bool) {
	m.showHelp = showHelp
	m.setSize(m.width, m.height)
}

// SetShowFullHelp allows to view extended help section, if visible.
func (m *Model) SetShowFullHelp(showFullHelp bool) {
	m.help.ShowAll = showFullHelp
	m.setSize(m.width, m.height)
}

func (m *Model) setSize(width, height int) {
	m.height = height
	m.width = width
	m.help.Width = width // this is required for full help
	m.recalculateSize()
	m.redrawLines()
}

func (m *Model) recalculateSize() {
	headerView := m.headerView()
	footerView := m.footerView()
	helpView := m.helpView()

	height := m.height

	height -= lipgloss.Height(helpView)
	height -= lipgloss.Height(headerView)
	height -= lipgloss.Height(footerView)

	// if viewport size changes, the text should be reformatted
	m.viewport.Height = max(height, 1)
}

func (m *Model) formatLines(lines []string) string {
	lineNumberStyle := lineNumberStylePlaceholder.Copy().Width(len(fmt.Sprintf("%d", len(lines))) + 1)
	lineNumberPlaceholder := lineNumberStyle.Render(fmt.Sprintf("%d", 1))
	availableWidth := m.width - lipgloss.Width(lineNumberPlaceholder) - lipgloss.Width(ellipsis)

	var buf strings.Builder

	for i, line := range lines {
		line = m.replaceTabsWithSpaces(line)
		lineNumber := lineNumberStyle.Render(fmt.Sprintf("%d", i+1))

		if lipgloss.Width(line) > availableWidth {
			line = truncate.StringWithTail(line, uint(availableWidth), ellipsis)
		}

		buf.WriteString(lipgloss.JoinHorizontal(lipgloss.Left, lineNumber, line))
		buf.WriteString(newLine)
	}

	return buf.String()
}

func (m *Model) replaceTabsWithSpaces(line string) string {
	return strings.ReplaceAll(line, "\t", "    ")
}

func (m *Model) headerView() string {
	truncatedTitle := m.title

	if maxWidth := m.width - 5; len(m.title) > maxWidth {
		truncatedTitle = fmt.Sprintf("%s%s", ellipsis, m.title[len(m.title)-maxWidth:])
	}

	title := fileTitleStyle.Render(truncatedTitle)
	line := strings.Repeat("─", max(0, m.width-lipgloss.Width(title)))

	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m *Model) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.width-lipgloss.Width(info)))

	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func (m *Model) helpView() (res string) {
	if m.showHelp {
		return helpStyle.Render(m.help.View(m))
	}

	return ""
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}
