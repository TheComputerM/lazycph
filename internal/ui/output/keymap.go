package output

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/viewport"
)

type KeyMap struct {
	viewport.KeyMap
}

func DefaultKeyMap() KeyMap {
	keymap := viewport.DefaultKeyMap()
	return KeyMap{KeyMap: keymap}
}

var _ help.KeyMap = (*Model)(nil)

func (m *Model) ShortHelp() []key.Binding {
	return []key.Binding{}
}

func (m *Model) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}
