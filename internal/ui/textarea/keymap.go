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

var _ help.KeyMap = KeyMap{}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Paste}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Paste}}
}
