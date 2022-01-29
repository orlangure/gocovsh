package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/orlangure/gocovsh/pkg/codeview"
	"golang.org/x/tools/cover"
)

const (
	// TODO: support themes + dark/light mode.
	primaryColor   = "#00ff00"
	secondaryColor = "#ff0000"
	inactiveColor  = "#505050"
)

var (
	neutralLine   = lipgloss.NewStyle().Foreground(lipgloss.Color(inactiveColor))
	coveredLine   = lipgloss.NewStyle().Foreground(lipgloss.Color(primaryColor))
	uncoveredLine = lipgloss.NewStyle().Foreground(lipgloss.Color(secondaryColor))
)

type viewName string

const (
	activeViewList viewName = "list"
	activeViewCode viewName = "code"
)

type model struct {
	activeView viewName

	list  list.Model
	items []list.Item

	code codeview.Model

	profileFilename     string
	detectedPackageName string

	err   error
	ready bool
}

func (m model) Init() tea.Cmd {
	return loadProfiles(m.profileFilename)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.updateWindowSize(msg.Width, msg.Height)

	case error:
		return m.onError(msg)

	case []*cover.Profile:
		return m.onProfilesLoaded(msg)

	case fileContents:
		return m.onFileContentLoaded(msg)

	case tea.KeyMsg:
		if m, cmd := m.onKeyPressed(msg.String()); m != nil {
			return m, cmd
		}
	}

	var cmd tea.Cmd

	switch m.activeView {
	case activeViewList:
		m.list, cmd = m.list.Update(msg)
	case activeViewCode:
		m.code, cmd = m.code.Update(msg)
	}

	return m, cmd
}

func (m model) View() string {
	if m.err != nil {
		// TODO: add error style
		return fmt.Sprintf("Error: %s\nPress any key to exit\n", m.err)
	}

	if !m.ready {
		return "Initializing..."
	}

	if m.isCodeView() {
		return m.code.View()
	}

	if m.isListView() {
		return m.list.View()
	}

	return "Unknown view"
}

func (m model) isCodeView() bool {
	return m.activeView == activeViewCode
}

func (m model) isListView() bool {
	return m.activeView == activeViewList
}

func (m *model) updateWindowSize(width, height int) (tea.Model, tea.Cmd) {
	if !m.ready {
		m.code = codeview.New(width, height)

		m.list = list.New([]list.Item{}, coverProfileDelegate{}, width, height-1)
		m.list.Title = "Code coverage is available for the following files:"
		m.list.SetShowStatusBar(true)
		m.list.SetFilteringEnabled(true)
		m.list.Styles.Title = titleStyle
		m.list.FilterInput.PromptStyle = m.list.FilterInput.PromptStyle.Copy().Margin(1, 0, 0, 0)
		m.list.Styles.PaginationStyle = paginationStyle
		m.list.Styles.HelpStyle = helpStyle
		m.list.Styles.StatusBar = statusBarStyle

		m.ready = true
	}

	m.code.SetWidth(width)
	m.code.SetHeight(height)

	m.list.SetWidth(width)
	m.list.SetHeight(height - 1)

	return *m, nil
}

func (m *model) onError(err error) (tea.Model, tea.Cmd) {
	m.err = err
	return m, nil
}

func (m *model) onProfilesLoaded(profiles []*cover.Profile) (tea.Model, tea.Cmd) {
	if len(profiles) == 0 {
		m.err = fmt.Errorf("no profiles found; you may need to run `go test -coverprofile=coverage.out`")
		return m, nil
	}

	m.items = make([]list.Item, len(profiles))

	for i, p := range profiles {
		// package name should already be set
		p.FileName = strings.TrimPrefix(p.FileName, m.detectedPackageName+"/")
		m.items[i] = &coverProfile{profile: p}
	}

	m.list.SetItems(m.items)

	return m, nil
}

func (m *model) onFileContentLoaded(content []string) (tea.Model, tea.Cmd) {
	m.code.SetContent(content)
	m.activeView = activeViewCode

	return m, nil
}

func (m *model) onKeyPressed(key string) (tea.Model, tea.Cmd) {
	// exit on any key in case of error
	if m.err != nil {
		return m, tea.Quit
	}

	// don't match any of the keys below if we're actively filtering.
	if m.list.FilterState() == list.Filtering {
		return nil, nil
	}

	switch key {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "esc":
		if m.isCodeView() {
			m.activeView = activeViewList
			return m, nil
		}

		if m.isListView() {
			// do not exit on "esc"
			if m.list.FilterState() == list.Unfiltered {
				return m, nil
			}
		}

	case "enter":
		item, ok := m.list.SelectedItem().(*coverProfile)
		if ok {
			m.code.SetTitle(item.profile.FileName)

			return m, loadFile(item.profile.FileName, item.profile)
		}

		return m, nil
	}

	return nil, nil
}

func loadProfiles(profileFilename string) tea.Cmd {
	return func() tea.Msg {
		pkg, err := determinePackageName()
		if err != nil {
			return fmt.Errorf("failed to determine package name: %w", err)
		}

		profiles, err := cover.ParseProfiles(profileFilename)
		if err != nil {
			return fmt.Errorf("failed to parse cover profiles: %w", err)
		}

		for i, p := range profiles {
			p.FileName = strings.TrimPrefix(p.FileName, pkg+"/")
			profiles[i] = p
		}

		return profiles
	}
}

func determinePackageName() (string, error) {
	bs, err := os.ReadFile("go.mod")
	if err != nil {
		return "", fmt.Errorf("cannot open go.mod file: %w", err)
	}

	matches := regexp.MustCompile(`module\s+(.+)`).FindStringSubmatch(string(bs))

	if len(matches) == 0 {
		return "", fmt.Errorf("could not determine package name; make sure go.mod file is valid")
	}

	return matches[1], nil
}

type fileContents []string

// nolint: gosec
func loadFile(filename string, profile *cover.Profile) tea.Cmd {
	return func() tea.Msg {
		f, err := os.Open(filename)
		if err != nil {
			return fmt.Errorf("could not open file %s: %w", filename, err)
		}

		defer func() { _ = f.Close() }()

		scanner := bufio.NewScanner(f)

		var lines []string

		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}

		highlightedText, err := colorize(lines, profile)
		if err != nil {
			return fmt.Errorf("could not colorize file %s: %w", filename, err)
		}

		return highlightedText
	}
}

func colorize(lines []string, profile *cover.Profile) (contents fileContents, err error) {
	defer func() {
		if rr := recover(); rr != nil {
			err = fmt.Errorf("unexpected error in coverage profile: potentially mismatching profile and source code: %v", rr)
		}
	}()

	buf := make(fileContents, 0, len(lines))

	for lineIdx, blockIdx := 0, 0; lineIdx < len(lines); lineIdx++ {
		line, block := lines[lineIdx], profile.Blocks[blockIdx]

		coverageStyle := uncoveredLine
		if block.Count > 0 {
			coverageStyle = coveredLine
		}

		adjustedStartLine, adjustedEndLine := block.StartLine-1, block.EndLine-1

		// before the first block - not covered
		if lineIdx < adjustedStartLine {
			buf = append(buf, neutralLine.Render(line))
			continue
		}

		// first line - highlight from the start col
		if lineIdx == adjustedStartLine {
			uncoveredPart := neutralLine.Render(line[:block.StartCol-1])
			coveredPart := coverageStyle.Render(line[block.StartCol-1:])
			buf = append(buf, fmt.Sprintf("%s%s", uncoveredPart, coveredPart))

			continue
		}

		// inside any block - can be anything
		if lineIdx >= adjustedStartLine && lineIdx <= adjustedEndLine {
			// TODO: support end column as well
			if block.NumStmt > 0 {
				buf = append(buf, coverageStyle.Render(line))
			} else {
				buf = append(buf, neutralLine.Render(line))
			}

			continue
		}

		// after a block - might be the last block or just bump the block
		if lineIdx > adjustedEndLine {
			// when there are more blocks, bump the block
			if blockIdx < len(profile.Blocks)-1 {
				blockIdx++
				lineIdx--
			} else {
				buf = append(buf, neutralLine.Render(line))
			}
		}
	}

	return buf, nil
}
