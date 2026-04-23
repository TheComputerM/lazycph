package workspace

import (
	"charm.land/lipgloss/v2"
)

func (m *Model) updateLayout() {
	availWidth, availHeight := m.width, m.height

	// help
	m.Help.SetWidth(availWidth)
	availHeight -= lipgloss.Height(m.helpView())

	// testcase list
	availWidth -= m.TestCaseList.GetWidth()
	m.TestCaseList.SetHeight(availHeight)

	// editor
	editorHeight := availHeight/2 - 1
	halfWidth := availWidth / 2
	m.Input.SetWidth(halfWidth)
	m.Input.SetHeight(editorHeight)
	m.Expected.SetWidth(halfWidth)
	m.Expected.SetHeight(editorHeight)
	availHeight -= editorHeight

	// output
	m.Output.SetWidth(availWidth)
	m.Output.SetHeight(availHeight)
}
