// Package errorview implements a view that displays errors in a user-friendly
// way. Pressing any key in this view will exit the program.
package errorview

import (
	"errors"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle   = lipgloss.NewStyle().Margin(1, 0, 0, 1).Bold(true).Foreground(lipgloss.Color("#ff5555"))
	messageSTyle = lipgloss.NewStyle().Margin(1, 0, 0, 1)
	errorStyle   = lipgloss.NewStyle().Margin(1, 0, 0, 1).Foreground(lipgloss.Color("#c0c0c0"))
	footerStyle  = lipgloss.NewStyle().Margin(1)

	footer             = footerStyle.Render("Press any key to exit")
	originalErrorTitle = "The original error was:"
)

// ErrorView is a type that can be used to render an error.
type ErrorView interface {
	Title() string
	Description() string
	OriginalError() error
}

// New creates a new error view.
func New(_ ErrorView) Model {
	return Model{}
}

// Model is the model for the error view.
type Model struct {
	err error
}

// Update is called by bubbletea every time there is a new event.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if _, ok := msg.(tea.KeyMsg); ok {
		return m, tea.Quit
	}

	return m, nil
}

// View renders the error view.
func (m *Model) View() string {
	var (
		errBox        ErrorView
		title         string
		description   string
		originalError string
	)

	if errors.As(m.err, &errBox) {
		title = titleStyle.Render(errBox.Title())
		description = messageSTyle.Render(errBox.Description())

		if err := errBox.OriginalError(); err != nil {
			originalError = errorStyle.Render(fmt.Sprintf("%s\n%s", originalErrorTitle, m.err))
		}
	} else {
		title = titleStyle.Render("Error")
		description = messageSTyle.Render("Unexpected error")
		originalError = errorStyle.Render(fmt.Sprintf("%s", m.err))
	}

	return fmt.Sprintf("%s\n%s\n%s\n%s", title, description, originalError, footer)
}

// SetError sets the error to be displayed.
func (m *Model) SetError(err error) {
	m.err = err
}
