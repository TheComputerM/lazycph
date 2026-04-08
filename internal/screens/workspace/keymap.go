package workspace

import (
	"charm.land/bubbles/v2/key"
)

type KeyMap struct {
	Quit key.Binding
	Run  key.Binding
	Next key.Binding
	Prev key.Binding
	Help key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		Quit: key.NewBinding(key.WithKeys("ctrl+c", "esc"), key.WithHelp("ctrl+c", "quit")),
		Run:  key.NewBinding(key.WithKeys("ctrl+r"), key.WithHelp("ctrl+r", "run")),
		Next: key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next")),
		Prev: key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "prev")),
		Help: key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
	}
}

func (m *Model) ShortHelp() []key.Binding {
	global := []key.Binding{m.keyMap.Run, m.keyMap.Help, m.keyMap.Quit}
	focused := m.currentlyFocused().ShortHelp()
	return append(focused, global...)
}

func (m *Model) FullHelp() [][]key.Binding {
	global := [][]key.Binding{
		{m.keyMap.Run},
		{m.keyMap.Next, m.keyMap.Prev},
		{m.keyMap.Help, m.keyMap.Quit},
	}
	focused := m.currentlyFocused().FullHelp()
	return append(focused, global...)
}

func (m *Model) HelpView() string {
	return m.Help.View(m)
}
