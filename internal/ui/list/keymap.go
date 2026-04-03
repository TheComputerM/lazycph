package list

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
)

type KeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Create key.Binding
	Delete key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "move down"),
		),
		Create: key.NewBinding(
			key.WithKeys("+"),
			key.WithHelp("+", "create"),
		),
		Delete: key.NewBinding(
			key.WithKeys("-"),
			key.WithHelp("-", "delete"),
		),
	}
}

var _ help.KeyMap = (*Model)(nil)

func (m *Model) ShortHelp() []key.Binding {
	return []key.Binding{m.keyMap.Create, m.keyMap.Delete}
}

func (m *Model) FullHelp() [][]key.Binding {
	return [][]key.Binding{{m.keyMap.Up, m.keyMap.Down}, {m.keyMap.Create, m.keyMap.Delete}}
}
