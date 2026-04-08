package textarea

import (
	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
)

type Model struct {
	*textarea.Model

	keyMap KeyMap
}

func New(placeholder string) Model {
	keymap := DefaultKeyMap()

	model := textarea.New()
	model.Placeholder = placeholder
	model.KeyMap = keymap.KeyMap

	return Model{
		Model:  &model,
		keyMap: keymap,
	}
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	model, cmd := m.Model.Update(msg)
	m.Model = &model
	return cmd
}

func Blink() tea.Msg {
	return textarea.Blink()
}
