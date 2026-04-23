package companion

import (
	"fmt"
	"os"
	"path/filepath"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/thecomputerm/lazycph/internal/core"
	"github.com/thecomputerm/lazycph/internal/screens/filepicker"
)

type Model struct {
	Input textinput.Model
	Help  help.Model

	base          string
	data          Data
	keyMap        KeyMap
	width, height int
}

var _ tea.Model = (*Model)(nil)

func New(data Data) Model {
	ti := textinput.New()

	ti.Placeholder = "TODO"
	ti.Validate = func(s string) error {
		ext := filepath.Ext(s)
		if _, exists := core.Engines[ext]; !exists {
			return fmt.Errorf("no engine for '%s'", ext)
		}

		return nil
	}
	ti.SetValue(fmt.Sprintf("%s/%s.xxx", data.Group, data.Name))

	ti.Focus()

	base, _ := os.Getwd()

	return Model{
		Input: ti,
		Help:  help.New(),

		base:   base,
		data:   data,
		keyMap: DefaultKeyMap(),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(tea.RequestWindowSize, tea.RequestBackgroundColor, textinput.Blink, requestServer)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.Input.SetWidth(min(m.width-2, 64))
	case tea.BackgroundColorMsg:
		m.Input.SetStyles(textinput.DefaultStyles(msg.IsDark()))
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, m.keyMap.Create):
			value := m.Input.Value()
			if m.Input.Err != nil {
				return m, nil
			}
			return m, func() tea.Msg {
				// fpath is the full path to the new source code file
				fpath := filepath.Join(m.base, value)
				if err := os.MkdirAll(filepath.Dir(fpath), 0o755); err != nil {
					return err
				}

				file, err := os.Create(fpath)
				if err != nil {
					return err
				}
				if err := m.data.toTestCaseList().Save(fpath); err != nil {
					return err
				}
				if err := file.Close(); err != nil {
					return err
				}

				return filepicker.FileSelectedMsg{Path: fpath}
			}
		case key.Matches(msg, m.keyMap.Back):
			// TODO: navigate back to the previous screen
			return m, tea.Quit
		case key.Matches(msg, m.keyMap.Quit):
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.Input, cmd = m.Input.Update(msg)

	return m, cmd
}

func (m Model) View() tea.View {
	helpView := m.Help.View(m.keyMap)

	container := lipgloss.NewStyle().
		Align(lipgloss.Center, lipgloss.Center).
		Width(m.width).
		Height(m.height - lipgloss.Height(helpView))

	var errorMessage string

	if m.Input.Err != nil {
		errorMessage = m.Input.Err.Error()
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		container.Render(lipgloss.JoinVertical(
			lipgloss.Center,
			lipgloss.NewStyle().Hyperlink(m.data.URL).Render(m.data.URL),
			"",
			m.Input.View(),
			lipgloss.NewStyle().Margin(1, 0).Foreground(lipgloss.Red).Render(errorMessage),
			lipgloss.NewStyle().Faint(true).Render("CWD: "+m.base),
		)),
		helpView,
	)

	v := tea.NewView(content)
	v.AltScreen = true
	v.WindowTitle = "LazyCPH - Competitive Companion"

	return v
}
