package textarea

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textarea"
)

type KeyMap struct {
	textarea.KeyMap
}

func DefaultKeyMap() KeyMap {
	keymap := textarea.DefaultKeyMap()
	return KeyMap{keymap}
}

var _ help.KeyMap = (*Model)(nil)

func (m *Model) ShortHelp() []key.Binding {
	return []key.Binding{m.keyMap.Paste}
}

func (m *Model) FullHelp() [][]key.Binding {
	return [][]key.Binding{{m.keyMap.Paste}}
}
