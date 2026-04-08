package textarea

import (
	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
)

type Model struct {
	*textarea.Model

	keyMap KeyMap
	value  *string
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
	previousValue := m.Value()

	model, cmd := m.Model.Update(msg)
	m.Model = &model

	// TODO: need better way to track value changes
	if previousValue != m.Value() && m.value != nil {
		*m.value = m.Value()
	}

	return cmd
}

func (m *Model) BindValue(value *string) {
	m.value = value
	m.SetValue(*value)
}

func Blink() tea.Msg {
	return textarea.Blink()
}
