package gocovshtest

import (
	"os"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/orlangure/gocovsh/internal/model"
	"github.com/orlangure/gocovsh/internal/styles"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	lipgloss.SetColorProfile(termenv.TrueColor)
	styles.SetTheme()

	os.Exit(m.Run())
}

type modelTest struct {
	*testing.T

	profileFilename string
	codeRoot        string
	requestedFiles  []string
	filteredLines   map[string][]int

	m *model.Model
}

func (t *modelTest) init() tea.Cmd {
	t.m = model.New(
		model.WithProfileFilename(t.profileFilename),
		model.WithCodeRoot(t.codeRoot),
		model.WithRequestedFiles(t.requestedFiles),
		model.WithFilteredLines(t.filteredLines),
	)

	initCmd := t.m.Init()
	require.NotNil(t, initCmd)

	return initCmd
}

// nolint: unparam
func (t *modelTest) sendWindowSizeMsg(width, height int) (tea.Model, tea.Cmd) {
	msg := tea.WindowSizeMsg{Width: width, Height: height}

	return t.m.Update(msg)
}

func (t *modelTest) sendProfilesMsg(profilesMsg tea.Msg) (tea.Model, tea.Cmd) {
	return t.m.Update(profilesMsg)
}

func (t *modelTest) sendFileContentsMsg(fileContents tea.Msg) (tea.Model, tea.Cmd) {
	return t.m.Update(fileContents)
}

func (t *modelTest) sendErrorMsg(err tea.Msg) (tea.Model, tea.Cmd) {
	return t.m.Update(err)
}

func (t *modelTest) sendEnterKey() (tea.Model, tea.Cmd) {
	return t.m.Update(tea.KeyMsg(tea.Key{Type: tea.KeyEnter}))
}

func (t *modelTest) sendEscKey() (tea.Model, tea.Cmd) {
	return t.m.Update(tea.KeyMsg(tea.Key{Type: tea.KeyEsc}))
}

func (t *modelTest) sendLetterKey(letter rune) (tea.Model, tea.Cmd) {
	return t.m.Update(tea.KeyMsg(tea.Key{
		Type:  tea.KeyRunes,
		Runes: []rune{letter},
	}))
}
