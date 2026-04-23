package companion

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Model struct {
	Input textinput.Model
	Help  help.Model

	keyMap        KeyMap
	width, height int
}

var _ tea.Model = (*Model)(nil)

func New() Model {
	ti := textinput.New()
	ti.Placeholder = ""
	ti.Focus()

	return Model{
		Input: ti,
		Help:  help.New(),

		keyMap: DefaultKeyMap(),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(tea.RequestWindowSize, tea.RequestBackgroundColor, textinput.Blink)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.BackgroundColorMsg:
		m.Input.SetStyles(textinput.DefaultStyles(msg.IsDark()))
	case tea.KeyPressMsg:
		switch {
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
		Width(m.width).
		Height(m.height-lipgloss.Height(helpView)).
		Align(lipgloss.Center, lipgloss.Center)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		container.Render(lipgloss.JoinVertical(
			lipgloss.Center,
			lipgloss.NewStyle().Render("TODO: url"),
			"",
			m.Input.View(),
			"",
			lipgloss.NewStyle().Faint(true).Render("Current Working Directory"),
		)),
		helpView,
	)

	v := tea.NewView(content)
	v.AltScreen = true
	v.WindowTitle = "LazyCPH - Companion"

	return v
}
