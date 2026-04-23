package list

import (
	tea "charm.land/bubbletea/v2"
	"github.com/thecomputerm/lazycph/internal/core"
)

type TestCaseSelectedMsg struct {
	Index    int
	TestCase *core.TestCase
}

type TestCaseExecuteMsg struct {
	TestCase *core.TestCase
}

func (m *Model) SelectTestCase(index int) tea.Cmd {
	if !(index >= 0 && index < len(m.Items)) {
		return nil
	}

	m.index = index
	return func() tea.Msg {
		return TestCaseSelectedMsg{Index: index, TestCase: m.Selected()}
	}
}
