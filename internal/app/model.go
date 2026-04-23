package app

import (
	"errors"
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/thecomputerm/lazycph/internal/core"
	"github.com/thecomputerm/lazycph/internal/screens/companion"
	"github.com/thecomputerm/lazycph/internal/screens/filepicker"
	"github.com/thecomputerm/lazycph/internal/screens/workspace"
)

type Model struct {
	// a path that determines the state of the app.
	// If it is a directory, the file picker is shown.
	// If it is a file, the workspace is shown.
	state string

	active tea.Model

	companionMode bool
}

var _ tea.Model = (*Model)(nil)

func activeModelFromState(state string) (tea.Model, error) {
	info, err := os.Stat(state)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("Path %s does not exist", state)
		}
		return nil, fmt.Errorf("Failed to stat path %s: %v", state, err)
	}
	if info.IsDir() {
		return filepicker.New(state), nil
	}
	return workspace.New(state), nil
}

func New(state string, companionMode bool) (Model, error) {
	active, err := activeModelFromState(state)
	if err != nil {
		return Model{}, err
	}

	return Model{
		state:         state,
		active:        active,
		companionMode: companionMode,
	}, nil
}

func (m Model) Init() tea.Cmd {
	cmds := []tea.Cmd{m.active.Init()}
	if m.companionMode {
		cmds = append(cmds, companion.StartServer)
	}
	return tea.Batch(cmds...)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case error:
		return m, tea.Quit
	case core.NavigateMsg:
		if msg.Path != "" {
			m.state = msg.Path
		}
		var err error
		m.active, err = activeModelFromState(m.state)
		if err != nil {
			return m, func() tea.Msg {
				return err
			}
		}
		return m, m.active.Init()
	case companion.Data:
		m.active = companion.New(msg)
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
