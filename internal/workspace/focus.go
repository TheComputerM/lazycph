package workspace

import (
	"charm.land/bubbles/v2/help"
	tea "charm.land/bubbletea/v2"
)

type focusable interface {
	help.KeyMap
	Focus() tea.Cmd
	Blur()
}

func (m *Model) getFocusable(index uint) focusable {
	switch index {
	case 0:
		return &m.TestCaseList
	case 1:
		return &m.Input
	case 2:
		return &m.Expected
	case 3:
		return &m.Output
	}
	return nil
}

func (m *Model) focusNext() tea.Cmd {
	return m.focusOn((m.focused + 1) % 4)
}

func (m *Model) focusPrev() tea.Cmd {
	return m.focusOn((m.focused + 3) % 4)
}

func (m *Model) focusOn(index uint) tea.Cmd {
	prev := m.getFocusable(m.focused)
	prev.Blur()

	m.focused = index
	current := m.getFocusable(index)
	return current.Focus()
}
