package workspace

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
)

type KeyMap struct {
	Quit key.Binding

	Run    key.Binding
	RunAll key.Binding

	Next key.Binding
	Prev key.Binding

	SelectFile key.Binding
	Help       key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		Quit: key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "quit")),

		Run:    key.NewBinding(key.WithKeys("ctrl+r"), key.WithHelp("ctrl+r", "run")),
		RunAll: key.NewBinding(key.WithKeys("ctrl+shift+r"), key.WithHelp("ctrl+shift+r", "run all")),

		Next: key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next")),
		Prev: key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "prev")),

		SelectFile: key.NewBinding(key.WithKeys("ctrl+f", "esc"), key.WithHelp("esc", "select file")),
		Help:       key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
	}
}

func (m *Model) focusedKeyMap() help.KeyMap {
	switch m.focused {
	case 0:
		return m.TestCaseList.KeyMap
	case 1:
		return m.Input.KeyMap
	case 2:
		return m.Expected.KeyMap
	case 3:
		return m.Output.KeyMap
	}
	return nil
}

func (m *Model) ShortHelp() []key.Binding {
	global := []key.Binding{m.keyMap.Run, m.keyMap.Help, m.keyMap.Quit}
	focused := m.focusedKeyMap().ShortHelp()
	return append(focused, global...)
}

func (m *Model) FullHelp() [][]key.Binding {
	global := [][]key.Binding{
		{m.keyMap.Run, m.keyMap.RunAll},
		{m.keyMap.Next, m.keyMap.Prev},
		{m.keyMap.Help, m.keyMap.Quit},
	}
	focused := m.focusedKeyMap().FullHelp()
	return append(focused, global...)
}

func (m *Model) helpView() string {
	return m.Help.View(m)
}
