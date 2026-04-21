package textarea

import (
	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
)

type Model struct {
	*textarea.Model

	KeyMap KeyMap
	value  *string
}

func New(placeholder string) Model {
	keymap := DefaultKeyMap()

	model := textarea.New()
	model.Placeholder = placeholder
	model.KeyMap = keymap.KeyMap

	return Model{
		Model:  &model,
		KeyMap: keymap,
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	model, cmd := m.Model.Update(msg)
	m.Model = &model

	switch msg.(type) {
	case tea.KeyPressMsg:
		content := m.Value()
		if m.value != nil && content != *m.value {
			*m.value = content
		}
	}

	return m, cmd
}

func (m *Model) BindValue(value *string) {
	m.value = value
	m.SetValue(*value)
}

func Blink() tea.Msg {
	return textarea.Blink()
}
