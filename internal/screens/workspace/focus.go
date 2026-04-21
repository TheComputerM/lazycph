package workspace

import (
	tea "charm.land/bubbletea/v2"
)

type focusable interface {
	Focus() tea.Cmd
	Blur()
	Update(msg tea.Msg) tea.Cmd
}

func (m *Model) focusableAt(index uint) focusable {
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

func (m *Model) currentlyFocused() focusable {
	return m.focusableAt(m.focused)
}

func (m *Model) focusNext() tea.Cmd {
	return m.focusOn((m.focused + 1) % 4)
}

func (m *Model) focusPrev() tea.Cmd {
	return m.focusOn((m.focused + 3) % 4)
}

func (m *Model) focusOn(index uint) tea.Cmd {
	prev := m.focusableAt(m.focused)
	prev.Blur()

	m.focused = index
	current := m.focusableAt(m.focused)
	return current.Focus()
}
