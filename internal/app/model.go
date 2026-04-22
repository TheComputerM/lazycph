package app

import (
	"os"
	"path/filepath"

	tea "charm.land/bubbletea/v2"
	"github.com/thecomputerm/lazycph/internal/screens/filepicker"
	"github.com/thecomputerm/lazycph/internal/screens/workspace"
)

type Model struct {
	active tea.Model
}

var _ tea.Model = (*Model)(nil)

func New() Model {
	currentDirectory, _ := os.Getwd()

	return Model{
		active: filepicker.New(currentDirectory),
	}
}

func (m Model) Init() tea.Cmd {
	return m.active.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case error:
		return m, tea.Quit
	case filepicker.FileSelectedMsg:
		m.active = workspace.New(msg)
		return m, m.active.Init()
	case workspace.SelectFileMsg:
		m.active = filepicker.New(filepath.Dir(msg.PreviousFile))
		return m, m.active.Init()
	default:
		var cmd tea.Cmd
		m.active, cmd = m.active.Update(msg)
		return m, cmd
	}
}

func (m Model) View() tea.View {
	return m.active.View()
}
