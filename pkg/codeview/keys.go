package codeview

import "github.com/charmbracelet/bubbles/key"

// KeyMap includes codeview key mappings.
type KeyMap struct {
	Up   key.Binding
	Down key.Binding
	Home key.Binding
	End  key.Binding
	Back key.Binding
	Quit key.Binding
}

// DefaultKeyMap is the default KeyMap used by codeview package.
var DefaultKeyMap = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	Home: key.NewBinding(
		key.WithKeys("home", "g"),
		key.WithHelp("g/home", "top"),
	),
	End: key.NewBinding(
		key.WithKeys("end", "G"),
		key.WithHelp("G/end", "bottom"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}
