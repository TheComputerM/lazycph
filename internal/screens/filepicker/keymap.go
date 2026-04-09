package filepicker

import (
	"charm.land/bubbles/v2/filepicker"
	"charm.land/bubbles/v2/key"
)

type KeyMap struct {
	filepicker.KeyMap
	Quit key.Binding
	Help key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		KeyMap: filepicker.DefaultKeyMap(),
		Quit:   key.NewBinding(key.WithKeys("ctrl+c", "q"), key.WithHelp("ctrl+c", "quit")),
		Help:   key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
	}
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Open, k.Help, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Up, k.Down}, {k.Back, k.Open}, {k.Help, k.Quit}}
}
