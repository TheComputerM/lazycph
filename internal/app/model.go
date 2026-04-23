package app

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	tea "charm.land/bubbletea/v2"
	"github.com/thecomputerm/lazycph/internal/core"
	"github.com/thecomputerm/lazycph/internal/screens/filepicker"
	"github.com/thecomputerm/lazycph/internal/screens/workspace"
)

type Model struct {
	// a path that determines the state of the app.
	// If it is a directory, the file picker is shown.
	// If it is a file, the workspace is shown.
	state string

	active tea.Model
}

var _ tea.Model = (*Model)(nil)

func New(state string) Model {
	var active tea.Model

	if info, err := os.Stat(state); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Fatalf("Path %s does not exist", state)
		}
		log.Fatalf("Failed to stat path %s: %v", state, err)
	} else if info.IsDir() {
		active = filepicker.New(state)
	} else {
		testcases, err := core.LoadTestCaseList(state)
		if err != nil {
			log.Fatalln(err)
		}
		active = workspace.New(filepicker.FileSelectedMsg{Path: state, TestCases: testcases})
	}

	return Model{
		state:  state,
		active: active,
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
		m.state = msg.Path
		m.active = workspace.New(msg)
		return m, m.active.Init()
	case workspace.SelectFileMsg:
		m.state = filepath.Dir(msg.PreviousFile)
		m.active = filepicker.New(m.state)
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
