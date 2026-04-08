package app

import (
	tea "charm.land/bubbletea/v2"
	"github.com/thecomputerm/lazycph/internal/screens/workspace"
)

type Model struct {
	Workspace workspace.Model
}

var _ tea.Model = (*Model)(nil)

func New() Model {
	return Model{
		Workspace: workspace.New(),
	}
}

func (m Model) Init() tea.Cmd {
	return m.Workspace.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.Workspace, cmd = m.Workspace.Update(msg)
	return m, cmd
}

func (m Model) View() tea.View {
	v := tea.NewView(m.Workspace.View())
	v.AltScreen = true
	v.WindowTitle = "LazyCPH"
	v.MouseMode = tea.MouseModeCellMotion
	return v
}
