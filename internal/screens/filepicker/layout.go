package filepicker

import (
	"charm.land/lipgloss/v2"
)

func (m *Model) updateLayout() {
	availWidth, availHeight := m.width, m.height

	// help
	m.Help.SetWidth(availWidth)

	availHeight -= lipgloss.Height(m.Help.View(m.keyMap))

	// picker adds extra newline https://github.com/charmbracelet/bubbles/pull/914
	availHeight -= 1

	// picker
	m.Picker.SetHeight(availHeight)
}
