package list

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
)

type KeyMap struct {
	Up      key.Binding
	Down    key.Binding
	Create  key.Binding
	Delete  key.Binding
	Execute key.Binding
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
			key.WithKeys("+", "c"),
			key.WithHelp("+/c", "create"),
		),
		Delete: key.NewBinding(
			key.WithKeys("-", "d"),
			key.WithHelp("-/d", "delete"),
		),
		Execute: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "run"),
		),
	}
}

var _ help.KeyMap = KeyMap{}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Create, k.Delete}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Create, k.Delete}, {k.Up, k.Down}}
}
