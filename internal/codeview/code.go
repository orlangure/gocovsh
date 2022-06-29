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
	"github.com/orlangure/gocovsh/internal/styles"
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

	blankBlockSeparatorStyle = lipgloss.NewStyle().
					MarginTop(1).MarginBottom(1).
					Foreground(lipgloss.Color(lineNumberColor))
)

type filteredLines struct {
	actualLines  []int
	contextLines map[int]bool
}

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
	viewport      viewport.Model
	help          help.Model
	width         int
	height        int
	title         string
	lines         []string
	filteredLines filteredLines
	showHelp      bool
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
	codeView := m.viewport.View()

	sections := make([]string, 0, 4)
	sections = append(sections, headerView, codeView, footerView)

	if helpView := m.helpView(); helpView != "" {
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

// SetFilteredLines sets the lines that should be displayed, while all other
// lines are hidden. If not set, everything is displayed.
func (m *Model) SetFilteredLines(filteredLines []int) {
	m.filteredLines = contextifyFilteredLines(filteredLines)
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
	if len(lines) == 0 {
		return ""
	}

	var (
		buf             strings.Builder
		filterApplied   = len(m.filteredLines.actualLines) > 0
		numberWidth     = len(fmt.Sprintf("%d", len(lines))) + 1
		lineNumberStyle = lineNumberStylePlaceholder.Copy().Width(numberWidth)
		printSingleLine = m.linePrinter(&buf, lineNumberStyle)
	)

	if filterApplied {
		lastPrintedLine := 0
		separator := blankBlockSeparatorStyle.Render(strings.Repeat("─", max(0, m.width)))

		for _, thisLineNumber := range m.filteredLines.actualLines {
			if thisLineNumber > len(lines) {
				break
			}

			if thisLineNumber-lastPrintedLine > 1 {
				buf.WriteString(separator)
				buf.WriteString(newLine)
			}

			drawPlus := false
			line := lines[thisLineNumber-1]

			if !m.filteredLines.contextLines[thisLineNumber] {
				drawPlus = true
			}

			printSingleLine(line, thisLineNumber, drawPlus)
			lastPrintedLine = thisLineNumber
		}
	} else {
		for i, line := range lines {
			printSingleLine(line, i+1, false)
		}
	}

	return buf.String()
}

type linePrinterFunc func(line string, number int, drawPlus bool)

func (m *Model) linePrinter(buf *strings.Builder, lineNumberStyle lipgloss.Style) linePrinterFunc {
	filterApplied := len(m.filteredLines.actualLines) > 0
	lineNumberPlaceholder := lineNumberStyle.Render("1")
	availableWidth := m.width - lipgloss.Width(lineNumberPlaceholder) - lipgloss.Width(ellipsis)
	renderedPlus := styles.CoveredLine.Render("+ ")
	renderedSpace := styles.NeutralLine.Render("  ")

	return func(line string, number int, drawPlus bool) {
		line = m.replaceTabsWithSpaces(line)
		lineNumber := lineNumberStyle.Render(fmt.Sprintf("%d", number))
		prefix := ""

		if filterApplied {
			if drawPlus {
				prefix = renderedPlus
			} else {
				prefix = renderedSpace
			}
		}

		if lipgloss.Width(line) > availableWidth {
			line = truncate.StringWithTail(line, uint(availableWidth), ellipsis)
		}

		buf.WriteString(lipgloss.JoinHorizontal(lipgloss.Left, prefix, lineNumber, line))
		buf.WriteString(newLine)
	}
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

func contextifyFilteredLines(input []int) filteredLines {
	if len(input) == 0 {
		return filteredLines{
			actualLines:  input,
			contextLines: map[int]bool{},
		}
	}

	extendedLines := make([]int, 0, len(input)+2)
	contextLines := make(map[int]bool, 2)
	lastAddedNumber := input[0]

	beforeFirst := lastAddedNumber - 1
	if beforeFirst > 0 {
		extendedLines = append(extendedLines, beforeFirst)
		contextLines[beforeFirst] = true
	}

	for _, lineNumber := range input {
		if lineNumber-lastAddedNumber > 1 {
			afterLast := lastAddedNumber + 1
			beforeThis := lineNumber - 1

			extendedLines = append(extendedLines, afterLast)
			contextLines[afterLast] = true

			if afterLast != beforeThis {
				extendedLines = append(extendedLines, beforeThis)
				contextLines[beforeThis] = true
			}
		}

		extendedLines = append(extendedLines, lineNumber)
		lastAddedNumber = lineNumber
	}

	lineAfterLast := input[len(input)-1] + 1
	extendedLines = append(extendedLines, lineAfterLast)
	contextLines[lineAfterLast] = true

	return filteredLines{
		actualLines:  extendedLines,
		contextLines: contextLines,
	}
}
