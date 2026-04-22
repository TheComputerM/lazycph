package filepicker

import (
	"charm.land/bubbles/v2/filepicker"
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/thecomputerm/lazycph/internal/core"
)

type Model struct {
	Picker filepicker.Model
	Help   help.Model

	keyMap        KeyMap
	width, height int
}

var _ tea.Model = (*Model)(nil)

func New(currentDirectory string) Model {
	fp := filepicker.New()
	fp.AutoHeight = false

	fp.CurrentDirectory = currentDirectory
	fp.AllowedTypes = make([]string, 0, len(core.Engines))
	for engine, _ := range core.Engines {
		fp.AllowedTypes = append(fp.AllowedTypes, engine)
	}

	return Model{
		Picker: fp,
		Help:   help.New(),
		keyMap: DefaultKeyMap(),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(tea.RequestWindowSize, m.Picker.Init())
}

type FileSelectedMsg struct {
	Path      string
	TestCases core.TestCaseList
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateLayout()
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, m.keyMap.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keyMap.Help):
			m.Help.ShowAll = !m.Help.ShowAll
			m.updateLayout()
		}
	}

	var cmd tea.Cmd
	m.Picker, cmd = m.Picker.Update(msg)

	if didSelect, path := m.Picker.DidSelectFile(msg); didSelect {
		testCases, err := core.LoadTestCaseList(path)
		if err != nil {
			return m, func() tea.Msg {
				return err
			}
		}

		return m, func() tea.Msg {
			return FileSelectedMsg{Path: path, TestCases: testCases}
		}

	}

	return m, cmd
}

func (m Model) View() tea.View {
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		m.Picker.View(),
		m.Help.View(m.keyMap),
	)

	v := tea.NewView(content)
	v.AltScreen = true
	v.WindowTitle = "LazyCPH - Select File"

	return v
}
