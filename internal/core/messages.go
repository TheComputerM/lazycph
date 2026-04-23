package core

import tea "charm.land/bubbletea/v2"

// NavigateMsg is sent to navigate to a new path.
// If Path is empty, the app reverts to the active model
// derived from the current state, simulating popup dismissal.
// If Path is non-empty, the app state updates and the
// active model is derived from the new path (dir → filepicker, file → workspace).
type NavigateMsg struct {
	Path string
}

// TestCaseExecutedMsg is dispatched when a TestCase has finished running.
type TestCaseExecutedMsg struct {
	TestCase *TestCase
}

func (tc *TestCase) ExecuteCmd(srcPath string) tea.Cmd {
	tc.Status = TestCaseStatusPending
	tc.Details = "Running..."

	return func() tea.Msg {
		tc.Execute(srcPath)
		return TestCaseExecutedMsg{TestCase: tc}
	}
}

func (list TestCaseList) SaveCmd(filePath string) tea.Cmd {
	return func() tea.Msg {
		return list.Save(filePath)
	}
}
